package database

import (
	"fmt"
	"log/slog"

	"github.com/elosanz/demo/internal/config"
	"github.com/elosanz/demo/internal/user"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Init initializes the database connection based on config.
func Init(cfg config.Config) (*gorm.DB, error) {
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
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	slog.Info("running auto-migrations")
	if err := db.AutoMigrate(&user.User{}); err != nil {
		return nil, fmt.Errorf("running auto-migrations: %w", err)
	}

	return db, nil
}
