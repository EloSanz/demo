package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/elosanz/demo/internal/notification"
	notificationhandler "github.com/elosanz/demo/internal/notification/handler"
	"github.com/elosanz/demo/internal/rickandmorty"
	rickandmortyhandler "github.com/elosanz/demo/internal/rickandmorty/handler"
	"github.com/elosanz/demo/internal/storage"
	storagehandler "github.com/elosanz/demo/internal/storage/handler"
	"github.com/elosanz/demo/internal/user"
	userhandler "github.com/elosanz/demo/internal/user/handler"
	"github.com/elosanz/demo/pkg/web"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	_ "github.com/elosanz/demo/docs"
	"gorm.io/gorm"
)

// NewHandler constructs the application router with all routes and middlewares.
func NewHandler(
	db *gorm.DB,
	userSvc user.UserService,
	rmSvc rickandmorty.RickAndMortyService,
	s3Svc storage.StorageService,
	notifSvc notification.NotificationService,
) http.Handler {
	mux := http.NewServeMux()
	
	sqlDB, _ := db.DB()
	
	// Handlers
	userH := userhandler.NewUserHandler(userSvc)
	rmH := rickandmortyhandler.NewRickAndMortyHandler(rmSvc)
	notifH := notificationhandler.NewNotificationHandler(notifSvc)

	// ─── System Routes ───────────────────────────────────────────────────────
	mux.HandleFunc("GET /health", web.HealthHandler(sqlDB))
	mux.Handle("GET /metrics", web.MetricsHandler())
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	// ─── Notification Routes ─────────────────────────────────────────────────
	mux.HandleFunc("POST /api/notifications", web.Adapt(notifH.Publish))

	// ─── User Routes ─────────────────────────────────────────────────────────
	mux.HandleFunc("GET /api/users", web.Adapt(userH.GetAll))
	mux.HandleFunc("GET /api/users/external", web.Adapt(userH.FetchExternal))
	mux.HandleFunc("POST /api/users/sync/{externalID}", web.Adapt(userH.SyncExternal))
	mux.HandleFunc("POST /api/users/transfer", web.Adapt(userH.TransferPoints))
	mux.HandleFunc("GET /api/users/{id}", web.Adapt(userH.GetByID))
	mux.HandleFunc("POST /api/users", web.Adapt(userH.Create))
	mux.HandleFunc("PUT /api/users/{id}", web.Adapt(userH.Update))
	mux.HandleFunc("DELETE /api/users/{id}", web.Adapt(userH.Delete))

	// ─── Rick and Morty Routes ───────────────────────────────────────────────
	mux.HandleFunc("GET /api/rickandmorty/characters", web.Adapt(rmH.Search))
	mux.HandleFunc("GET /api/rickandmorty/characters/{id}", web.Adapt(rmH.GetByID))

	// ─── Storage Routes ──────────────────────────────────────────────────────
	if s3Svc != nil {
		storageH := storagehandler.NewStorageHandler(s3Svc)
		mux.HandleFunc("POST /api/storage/upload", web.Adapt(storageH.Upload))
		mux.HandleFunc("GET /api/storage/files", web.Adapt(storageH.List))
		mux.HandleFunc("GET /api/storage/files/{name}", web.Adapt(storageH.Download))
		mux.HandleFunc("DELETE /api/storage/files/{name}", web.Adapt(storageH.Delete))
	}

	// ─── Debug/Test Routes ───────────────────────────────────────────────────
	mux.HandleFunc("GET /api/test/slow", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("slow request started")
		time.Sleep(15 * time.Second)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"slow success"}`))
	})

	// Wrap everything in Middlewares
	return web.Recovery(web.MetricsMiddleware(web.RequestLogger(mux)))
}
