package database

import (
	"database/sql"
	"fmt"
	"tunnerse-server/logger"
	"tunnerse-server/variables"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB *sql.DB
}

func InitDB() *Database {
	dbPath := variables.GetExecPath() + "/db.sqlite"
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.Log("ERROR", "failed to open database", []logger.LogDetail{
			{Key: "error", Value: err.Error()},
		})
		return nil
	}

	if err := db.Ping(); err != nil {
		logger.Log("ERROR", "failed to ping database", []logger.LogDetail{
			{Key: "error", Value: err.Error()},
		})
		return nil
	}

	if err := createTables(db); err != nil {
		logger.Log("ERROR", "failed to create tables", []logger.LogDetail{
			{Key: "error", Value: err.Error()},
		})
		return nil
	}

	return &Database{DB: db}
}

func createTables(db *sql.DB) error {
	createInfoTable := `
	CREATE TABLE IF NOT EXISTS Info (
		ID TEXT PRIMARY KEY,
		Pid INTEGER,
		Requests INTEGER,
		Healthchecks INTEGER,
		Warns INTEGER,
		Errors INTEGER
	);`
	if _, err := db.Exec(createInfoTable); err != nil {
		return fmt.Errorf("failed to create Info table: %w", err)
	}

	createTunnelTable := `
	CREATE TABLE IF NOT EXISTS Tunnel (
		ID TEXT PRIMARY KEY,
		Port TEXT NOT NULL,
		Url TEXT,
		Domain TEXT,
		Active INTEGER NOT NULL CHECK (Active IN (0,1)),
		CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(createTunnelTable); err != nil {
		return fmt.Errorf("failed to create Tunnel table: %w", err)
	}

	return nil
}
