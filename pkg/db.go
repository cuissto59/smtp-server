package pkg

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DatabaseHandler struct {
	DB *sql.DB
}

// Open initializes and opens a database connection.
func (handler *DatabaseHandler) Open() error {
	log.Println("Start opening database ...")
	db, err := sql.Open("sqlite3", "./sqlite3.db")
	if err != nil {
		return err
	}
	handler.DB = db
	log.Println(handler)
	return nil
}

// Close terminates the database connection.
func (handler *DatabaseHandler) Close() error {
	log.Println("Start closing database ....")
	if handler.DB != nil {
		return handler.DB.Close()
	}
	return nil
}

// CreateEmailsTable creates the 'emails' table if it does not exist.
func (handler *DatabaseHandler) CreateEmailsTable() error {
	log.Println("Start creating emails table...")
	createTableQuery := `
 CREATE TABLE IF NOT EXISTS emails (
    mail_id INTEGER PRIMARY KEY AUTOINCREMENT,
    "date" DATETIME NOT NULL,
    "to" VARCHAR(255) NOT NULL,
    "from" VARCHAR(255) NOT NULL,
    subject VARCHAR(255),
    body TEXT,
    mail BLOB,
    spam_score FLOAT,
    hash VARCHAR(64),
    content_type VARCHAR(100),
    recipient VARCHAR(255),
    has_attach BOOLEAN,
    ip_addr VARCHAR(45),
    return_path VARCHAR(255),
    is_tls BOOLEAN,
    message_id VARCHAR(255),
    reply_to VARCHAR(255),
    sender VARCHAR(255)
);
`
	_, err := handler.DB.Exec(createTableQuery)
	if err != nil {
		return err
	}
	log.Println(handler)

	return nil
}
