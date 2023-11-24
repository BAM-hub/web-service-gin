package dbConnection

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func goDotEnvVariable(key string) string {
    // load .env file
    err := godotenv.Load(".env")
  
    if err != nil {
      log.Fatalf("Error loading .env file")
    }
  
    return os.Getenv(key)
}

func Connect() (*sql.DB, error) {
    cfg := mysql.Config {
        User: goDotEnvVariable("DB_USER"),
        Passwd: goDotEnvVariable("DB_PASS"),
        Net: "tcp",
        Addr: "127.0.0.1:3306",
        DBName: "test",
    }

    var err error
    db, err = sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        return nil, err
    }
	println("Connected to database")
    return db, nil
}