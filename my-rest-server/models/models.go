package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

const CREATE_MAILS_TABLE_SQL_FORMAT = `
CREATE TABLE %v(
	id VARCHAR(255) PRIMARY KEY,
	account_email VARCHAR(255),
	sender_email VARCHAR(255),
	snippet VARCHAR(255)
)`

const INSERT_MAILS_TABLE_SQL_FORMAT = `
INSERT INTO %v(id, account_email, sender_email, snippet)
VALUES(?, ?, ?, ?)
`

// Setup initializes the database instance
func Setup() {
	// Capture connection properties.
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DB_USER")
	cfg.Passwd = os.Getenv("DB_PASSWORD")
	cfg.Net = "tcp"
	cfg.Addr = fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
	cfg.DBName = os.Getenv("DB_NAME")
	cfg.ParseTime = true

	// mysql is configured to use utc as checked by SELECT @@system_time_zone;
	cfg.Loc = time.UTC

	// Get a database handle.
	var err error
	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := DB.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	// note: schema of database is to do prior as defined sql/schema.sql
}
