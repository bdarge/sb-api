package util

import (
	"database/sql"
	"errors"
	. "github.com/bdarge/api/pkg/config"
	"github.com/bdarge/api/pkg/db"
	"github.com/bdarge/api/pkg/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
)

func Migrate(conf Config, handler db.Handler) error {
	dbInstance, _ := sql.Open("mysql", conf.DSN+"&multiStatements=true")

	err := handler.DB.AutoMigrate(
		&models.Business{},
		&models.Role{},
		&models.User{},
		&models.Account{},
		&models.Address{},
		&models.Customer{},
		&models.Transaction{},
		&models.TransactionItem{},
	)
	if err != nil {
		return err
	}

	driver, _ := mysql.WithInstance(dbInstance, &mysql.Config{})

	log.Printf("read migraction script from: %s", conf.MigrationDir)
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+conf.MigrationDir,
		conf.Database,
		driver,
	)

	if err != nil {
		log.Println("ERROR => ", err)
		return err
	}

	err = m.Steps(1)

	if errors.Is(err, os.ErrNotExist) || errors.Is(err, migrate.ErrShortLimit{}) || errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	if err != nil {
		log.Println("ERROR => ", err)
		return err
	}
	log.Printf("Applied migrations")

	return nil
}
