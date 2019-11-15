package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/AndrewDonelson/golog"
	// Blank Import
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	// Blank Import
	_ "github.com/golang-migrate/migrate/source/file"
)

// CreateDatabase Creates a new MySQL Database
func CreateDatabase() (*sql.DB, error) {
	serverName := "localhost:3306"
	user := "myuser"
	password := "pw"
	dbName := "demo"

	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&multiStatements=true", user, password, serverName, dbName)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	if err := migrateDatabase(db); err != nil {
		return db, err
	}

	return db, nil
}

func migrateDatabase(db *sql.DB) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	dir, err := os.Getwd()
	if err != nil {
		golog.Log.Fatal(err.Error())
	}

	migration, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s/db/migrations", dir),
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	golog.Log.Printf("Applying database migrations")
	err = migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	version, _, err := migration.Version()
	if err != nil {
		return err
	}

	golog.Log.Printf("Active database version: %d", version)

	return nil
}
