package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func NewSQLite(dsn string) *sql.DB {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	if _, err := db.Exec("PRAGMA journal_mode = WAL;"); err != nil {
		log.Fatalf("failed to set WAL mode: %v", err)
	}

	migrate(db)

	fmt.Println("[database] SQLite initialized")
	return db
}

func migrate(db *sql.DB) {
	createTableInvoice := `
    CREATE TABLE IF NOT EXISTS invoices (
        id             TEXT PRIMARY KEY,
        invoice_number TEXT NOT NULL,
        currency       TEXT NOT NULL,
        amount         REAL NOT NULL,
        status         TEXT NOT NULL,
        issued_at      TEXT NOT NULL,
        expired_at     TEXT NOT NULL,
        customer       TEXT NOT NULL,
        items          TEXT NOT NULL,
        summary        TEXT NOT NULL,
        metadata       TEXT NOT NULL,
		config         TEXT NOT NULL,
        created_at     TEXT NOT NULL
    );`

	if _, err := db.Exec(createTableInvoice); err != nil {
		log.Fatalf("failed to create table: %v", err)
	}

	if _, err := db.Exec("CREATE INDEX IF NOT EXISTS idx_invoices_id ON invoices(id);"); err != nil {
		log.Fatalf("failed to create index: %v", err)
	}

}
