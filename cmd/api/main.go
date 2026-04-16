package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"time"

	awsinfra "github.com/elosanz/demo/infrastructure/aws"
	"github.com/elosanz/demo/infrastructure/httpclient"
	"github.com/elosanz/demo/infrastructure/postgres"
	"github.com/elosanz/demo/internal/rickandmorty"
	rickandmortyhandler "github.com/elosanz/demo/internal/rickandmorty/handler"
	"github.com/elosanz/demo/internal/storage"
	storagehandler "github.com/elosanz/demo/internal/storage/handler"
	"github.com/elosanz/demo/internal/user"
	userhandler "github.com/elosanz/demo/internal/user/handler"
	"github.com/elosanz/demo/pkg/web"

	_ "modernc.org/sqlite"
)

func main() {
	ctx := context.Background()

	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "file::memory:?cache=shared"
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		slog.Error("opening database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := runMigrations(db); err != nil {
		slog.Error("running migrations", "error", err)
		os.Exit(1)
	}

	httpClient := &http.Client{Timeout: 10 * time.Second}

	jpBaseURL := envOr("JSONPLACEHOLDER_BASE_URL", "https://jsonplaceholder.typicode.com")
	rmBaseURL := envOr("RICKANDMORTY_BASE_URL", "https://rickandmortyapi.com/api")

	externalUserClient := httpclient.NewJSONPlaceholderClient(jpBaseURL, httpClient)
	externalRMClient := httpclient.NewRickAndMortyHTTPClient(rmBaseURL, httpClient)

	s3Client, err := awsinfra.NewS3Client(ctx)
	if err != nil {
		slog.Warn("AWS S3 unavailable — storage endpoints disabled", "error", err)
	}

	userRepo := postgres.NewUserSQLiteRepository(db)
	userSvc := user.NewUserService(userRepo, externalUserClient)
	rmSvc := rickandmorty.NewRickAndMortyService(externalRMClient)

	userH := userhandler.NewUserHandler(userSvc)
	rmH := rickandmortyhandler.NewRickAndMortyHandler(rmSvc)

	mux := http.NewServeMux()

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

	addr := ":" + envOr("PORT", "8080")
	slog.Info("starting server", "addr", addr)
	if err := http.ListenAndServe(addr, web.RequestLogger(mux)); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func runMigrations(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id         INTEGER  PRIMARY KEY AUTOINCREMENT,
		name       TEXT     NOT NULL,
		email      TEXT     NOT NULL UNIQUE,
		phone      TEXT,
		website    TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`)
	return err
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
