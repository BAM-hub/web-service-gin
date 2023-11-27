package dbConnection

import (
	"database/sql"
	utils "example/web-service-gin/utils"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var db *sql.DB



func Connect() (*sql.DB, error) {
    cfg := fmt.Sprintf("%s:%s@%s(%s)/%s",
         utils.GoDotEnvVariable("DB_USER"),
         utils.GoDotEnvVariable("DB_PASS"),
         "tcp",
         "127.0.0.1:3306",
         "test",
    )

    var err error
    db, err := sql.Open("mysql", cfg)

    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        return nil, err
    }

    driver, driverErr := mysql.WithInstance(db, &mysql.Config{})
    if driverErr != nil {
        log.Fatalf("Could not initialize migrate instance: %v\n", driverErr)
        return nil, driverErr
    }
    m, migrationErr := migrate.NewWithDatabaseInstance(
        "file://C:/Users/bshar/OneDrive/Desktop/fullstack/3tiers/web-service-gin/db/migrations",
        "mysql", driver)

    if migrationErr != nil {
        log.Fatalf("Could not initialize migrate instance: %v\n", migrationErr)
        return nil, migrationErr
    }

    err = m.Up()
    if err != nil && err != migrate.ErrNoChange {
        log.Fatalf("Could not migrate: %v\n", err)
        return nil, err
    }

    fmt.Println("Migration completed")
	println("Connected to database")
    return db, nil
}