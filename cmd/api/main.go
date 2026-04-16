package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/elosanz/demo/internal/api"
	"github.com/elosanz/demo/internal/config"
	"github.com/elosanz/demo/internal/database"
	"github.com/elosanz/demo/internal/notification"
	"github.com/elosanz/demo/internal/rickandmorty"
	"github.com/elosanz/demo/internal/storage"
	"github.com/elosanz/demo/internal/user"

	awsinfra "github.com/elosanz/demo/infrastructure/aws"
	"github.com/elosanz/demo/infrastructure/httpclient"
	infrapostgres "github.com/elosanz/demo/infrastructure/postgres"
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

	// 1. Database
	db, err := database.Init(cfg)
	if err != nil {
		return err
	}

	// 2. HTTP Clients
	httpClient := &http.Client{Timeout: 10 * time.Second}
	externalUserClient := httpclient.NewJSONPlaceholderClient(cfg.JSONPlaceholderURL, httpClient)
	externalRMClient := httpclient.NewRickAndMortyHTTPClient(cfg.RickAndMortyURL, httpClient)

	// 3. Infrastructure
	s3Client, err := awsinfra.NewS3Client(ctx)
	var s3Svc storage.StorageService
	if err == nil {
		s3Svc = storage.NewS3StorageService(s3Client)
	}

	// --------------------------------------------------------------------------
	// NOTIFICATION SERVICE: Dynamic Selection
	// --------------------------------------------------------------------------
	var notifSvc notification.NotificationService
	if cfg.NotificationEngine == "sqs" {
		sqsClient, err := awsinfra.NewSQSClient(ctx)
		if err != nil {
			return fmt.Errorf("initializing AWS SQS client: %w", err)
		}
		notifSvc = notification.NewSQSNotificationService(sqsClient, cfg.SQSQueueURL)
	} else {
		// Default to Memory for local development
		notifSvc = notification.NewMemoryNotificationService(100)
	}
	notifSvc.StartWorker(ctx)

	// 4. Services
	userRepo := infrapostgres.NewUserGORMRepository(db)
	userSvc := user.NewUserService(userRepo, externalUserClient)
	rmSvc := rickandmorty.NewRickAndMortyService(externalRMClient)

	// 5. Router Index
	handler := api.NewHandler(db, userSvc, rmSvc, s3Svc, notifSvc)

	// 6. Server Start
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
