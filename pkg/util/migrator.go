package util

import (
	"database/sql"
	"errors"
	"log/slog"
	"os"

	"github.com/bdarge/api/pkg/config"
	"github.com/bdarge/api/pkg/db"
	"github.com/bdarge/api/pkg/models"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
)

func logger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

// Migrate migrates
func Migrate(conf config.Config, handler db.Handler) error {
	logger()
	dbInstance, _ := sql.Open("mysql", conf.DSN+"&multiStatements=true")

	err := handler.DB.AutoMigrate(
		&models.Business{},
		&models.Role{},
		&models.User{},
		&models.Account{},
		&models.Address{},
		&models.Customer{},
		&models.CustomerAddress{},
		&models.Transaction{},
		&models.TransactionItem{},
		&models.Lang{},
	)
	if err != nil {
		return err
	}

	driver, _ := mysql.WithInstance(dbInstance, &mysql.Config{})

	slog.Info("read migration scripts", "directory", conf.MigrationDir)
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+conf.MigrationDir,
		conf.Database,
		driver,
	)

	if err != nil {
		slog.Error("Migration scripts failed", "error", err)
		return err
	}

	err = m.Steps(1)

	if errors.Is(err, os.ErrNotExist) || errors.Is(err, migrate.ErrShortLimit{}) || errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	if err != nil {
		slog.Error("Migration scripts failed", "error", err)
		return err
	}
	slog.Info("Migration done successfully")

	return nil
}
