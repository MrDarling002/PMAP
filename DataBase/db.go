package DataBase

import (
    "database/sql"
    "log"
    "os"

    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    connectionString := os.Getenv("DATABASE_URL")
    var dbErr error
    DB, dbErr = sql.Open("postgres", connectionString)
    if dbErr != nil {
        log.Fatal(dbErr)
    }
    if dbErr = DB.Ping(); dbErr != nil {
        log.Fatal(dbErr)
    }
    log.Println("Database connected!")
}