package main

import (
	"encoding/json"
	"fmt"
	"github.com/cuissto59/smtp-server/pkg"
	_ "github.com/mattn/go-sqlite3"
	"github.com/phires/go-guerrilla"
	"github.com/phires/go-guerrilla/backends"
	guerrillaLog "github.com/phires/go-guerrilla/log"
	"log"
	"net/http"
)

// Global database handler instance
var dbHandler *pkg.DatabaseHandler

// Email represents a simplified email structure
type Email struct {
	MailID  int    `json:"mail_id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func initDatabase() {
	var err error
	dbHandler = &pkg.DatabaseHandler{}

	if err = dbHandler.Open(); err != nil {
		log.Fatalf("Failed to open database %v", err)
	}

	if err = dbHandler.CreateEmailsTable(); err != nil {
		log.Fatalf("Failed to create emails table : %v", err)
	}

	// Defer closing the database connection until the main function exits

}

// Handle fetching emails from Redis
func handleListEmails(w http.ResponseWriter, r *http.Request) {

	if dbHandler == nil || dbHandler.DB == nil {
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	rows, err := dbHandler.DB.Query("SELECT `mail_id`, `from`, `to`, `subject`, `body` FROM emails")
	if err != nil {
		log.Fatalf("failed to fetch emails : %v", err)
		http.Error(w, "Failed to fetch emails", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var emails []Email
	for rows.Next() {
		var email Email
		if err := rows.Scan(&email.MailID, &email.From, &email.To, &email.Subject, &email.Body); err != nil {
			log.Fatalf("failed to scan email with id %d for : %v", email.MailID, err)
			http.Error(w, "Error scanning emails", http.StatusInternalServerError)
			return
		}
		emails = append(emails, email)
	}

	// Encode the emails as JSON and write to the response
	if err := json.NewEncoder(w).Encode(emails); err != nil {
		http.Error(w, "Failed to encode emails to JSON", http.StatusInternalServerError)
		return
	}
}

func setupRouters() {
	http.HandleFunc("/emails", handleListEmails)
}

func main() {

	// initialize database
	initDatabase()
	log.Println(dbHandler)

	cfg := &guerrilla.AppConfig{LogFile: guerrillaLog.OutputStdout.String(), AllowedHosts: []string{"*"}}
	sc := guerrilla.ServerConfig{
		ListenInterface: "127.0.0.1:2525",
		IsEnabled:       true,
	}
	cfg.Servers = append(cfg.Servers, sc)
	sqlstr := "INSERT INTO emails "
	sqlstr += "(`date`, `to`, `from`, `subject`, `body`,  `mail`, `spam_score`, "
	sqlstr += "`hash`, `content_type`, `recipient`, `has_attach`, `ip_addr`, "
	sqlstr += "`return_path`, `is_tls`, `message_id`, `reply_to`, `sender`)"
	sqlstr += " VALUES "

	values := "(CURRENT_TIMESTAMP, ?, ?, ?, ? , ?, 0, ?, ?, ?, 0, ?, ?, ?, ?, ?, ?)"
	bcfg := backends.BackendConfig{
		"save_process":       "sql|Debugger|headersParser",
		"mail_table":         "emails",
		"primary_mail_host":  "*",
		"log_received_mails": true,
		"sql_driver":         "sqlite3",
		"sql_dsn":            "./sqlite3.db",
		"sql_insert":         sqlstr,
		"sql_values":         values,
	}
	cfg.BackendConfig = bcfg

	guerrillaDeamon := guerrilla.Daemon{Config: cfg}
	if err := guerrillaDeamon.Start(); err != nil {
		fmt.Printf("Failed to start server: %s", err)
	}

	go func() {
		setupRouters()

		log.Println("Start HTTP server on:8080 ...")
		if err := http.ListenAndServe(":8083", nil); err != nil {
			log.Fatalf("Http server failed %s", err)
		}
	}()

	defer func() {
		if err := dbHandler.Close(); err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}()

	// Block indefinitely
	select {}
}
