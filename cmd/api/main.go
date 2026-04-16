package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/elosanz/demo/internal/config"
	"github.com/elosanz/demo/internal/rickandmorty"
	rickandmortyhandler "github.com/elosanz/demo/internal/rickandmorty/handler"
	"github.com/elosanz/demo/internal/storage"
	storagehandler "github.com/elosanz/demo/internal/storage/handler"
	"github.com/elosanz/demo/internal/user"
	userhandler "github.com/elosanz/demo/internal/user/handler"
	"github.com/elosanz/demo/pkg/web"

	awsinfra "github.com/elosanz/demo/infrastructure/aws"
	"github.com/elosanz/demo/infrastructure/httpclient"
	infrapostgres "github.com/elosanz/demo/infrastructure/postgres"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	// ─── Database ────────────────────────────────────────────────────────────
	var dialector gorm.Dialector
	if cfg.DBEngine == "postgres" {
		slog.Info("using postgres database")
		dialector = postgres.Open(cfg.DBDSN)
	} else {
		slog.Info("using sqlite database", "path", cfg.DBDSN)
		dialector = sqlite.Open(cfg.DBDSN)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}

	if err := db.AutoMigrate(&user.User{}); err != nil {
		return fmt.Errorf("running auto-migrations: %w", err)
	}

	sqlDB, _ := db.DB()

	// ─── External HTTP clients ───────────────────────────────────────────────
	httpClient := &http.Client{Timeout: 10 * time.Second}
	
	externalUserClient := httpclient.NewJSONPlaceholderClient(cfg.JSONPlaceholderURL, httpClient)
	externalRMClient := httpclient.NewRickAndMortyHTTPClient(cfg.RickAndMortyURL, httpClient)

	// ─── AWS S3 ──────────────────────────────────────────────────────────────
	s3Client, err := awsinfra.NewS3Client(ctx)
	if err != nil {
		slog.Warn("AWS S3 unavailable — storage endpoints disabled", "error", err)
	}

	// ─── Dependencies ────────────────────────────────────────────────────────
	userRepo := infrapostgres.NewUserGORMRepository(db)
	userSvc := user.NewUserService(userRepo, externalUserClient)
	rmSvc := rickandmorty.NewRickAndMortyService(externalRMClient)

	userH := userhandler.NewUserHandler(userSvc)
	rmH := rickandmortyhandler.NewRickAndMortyHandler(rmSvc)

	// ─── Routing ─────────────────────────────────────────────────────────────
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", web.HealthHandler(sqlDB))
	mux.Handle("GET /metrics", web.MetricsHandler())

	mux.HandleFunc("GET /api/users", web.Adapt(userH.GetAll))
	mux.HandleFunc("GET /api/users/external", web.Adapt(userH.FetchExternal))
	mux.HandleFunc("POST /api/users/sync/{externalID}", web.Adapt(userH.SyncExternal))
	mux.HandleFunc("GET /api/users/{id}", web.Adapt(userH.GetByID))
	mux.HandleFunc("POST /api/users", web.Adapt(userH.Create))
	mux.HandleFunc("PUT /api/users/{id}", web.Adapt(userH.Update))
	mux.HandleFunc("DELETE /api/users/{id}", web.Adapt(userH.Delete))

	mux.HandleFunc("GET /api/rickandmorty/characters", web.Adapt(rmH.Search))
	mux.HandleFunc("GET /api/rickandmorty/characters/{id}", web.Adapt(rmH.GetByID))

	if s3Client != nil {
		storageH := storagehandler.NewStorageHandler(storage.NewS3StorageService(s3Client))
		mux.HandleFunc("POST /api/storage/upload", web.Adapt(storageH.Upload))
		mux.HandleFunc("GET /api/storage/files", web.Adapt(storageH.List))
		mux.HandleFunc("GET /api/storage/files/{name}", web.Adapt(storageH.Download))
		mux.HandleFunc("DELETE /api/storage/files/{name}", web.Adapt(storageH.Delete))
	}

	// Slow endpoint for testing Graceful Shutdown
	mux.HandleFunc("GET /api/test/slow", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("slow request started")
		time.Sleep(15 * time.Second)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"slow success"}`))
		slog.Info("slow request finished")
	})

	// ─── Server Start ────────────────────────────────────────────────────────
	// Chain middlewares: Recovery -> Metrics -> Logger -> Mux
	handler := web.Recovery(web.MetricsMiddleware(web.RequestLogger(mux)))

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,
	}

	serverError := make(chan error, 1)
	go func() {
		slog.Info("starting server", "addr", server.Addr)
		serverError <- server.ListenAndServe()
	}()

	select {
	case err := <-serverError:
		return err
	case <-ctx.Done():
		slog.Info("shutting down server")
		
		// Increased timeout to 20s to allow slow requests to finish
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			server.Close()
			return err
		}
		slog.Info("graceful shutdown complete")
	}

	return nil
}
