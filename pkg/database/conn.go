package database

import (
        "database/sql"

        "github.com/golang-migrate/migrate/v4"
        "github.com/golang-migrate/migrate/v4/database/sqlite3"
        _ "github.com/golang-migrate/migrate/v4/source/file"
        _ "github.com/mattn/go-sqlite3"
  )

var DB *sql.DB

func Conn() error {
    var err error
    DB, err = sql.Open("sqlite3", "./database.db")
    if err != nil {
            return err
    }

    driver, err := sqlite3.WithInstance(DB, &sqlite3.Config{})
    if err != nil {
        return err
    }

    m, err := migrate.NewWithDatabaseInstance(
            "file://db/migrations",
            "sqlite3", driver,
    )
    if err != nil {
        return err
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }

    return nil
}
