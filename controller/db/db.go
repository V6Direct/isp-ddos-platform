package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func Open(path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("create db dir: %w", err)
	}
	conn, err := sql.Open("sqlite3", path+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(1) // SQLite: single writer
	return conn, nil
}

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS routers (
			id			INTEGER PRIMARY KEY AUTOINCREMENT,
			name		TEXT NOT NULL UNIQUE,
			ip			TEXT NOT NULL,
			token		TEXT NOT NULL UNIQUE,
			mode		TEXT NOT NULL DEFAULT 'unknown',
			iface		TEXT NOT NULL DEFAULT '',
			last_seen	DATETIME,
			created_at	DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS rules (
			id				 INTEGER PRIMARY KEY AUTOINCREMENT,
			router_id		 INTEGER REFERENCES routers(id) ON DELETE CASCADE,
			allowlist_json	 TEXT NOT NULL DEFAULT '[]',
			denylist_json	 TEXT NOT NULL DEFAULT '[]',
			rate_limits_json TEXT NOT NULL DEFAULT '{}',
			updated_at		 DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS global_rules (
			id				 INTEGER PRIMARY KEY CHECK (id = 1),
			allowlist_json	 TEXT NOT NULL DEFAULT '[]',
			denylist_json	 TEXT NOT NULL DEFAULT '[]',
			rate_limits_json TEXT NOT NULL DEFAULT '{}',
			updated_at		 DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		INSERT OR IGNORE INTO global_rules (id) VALUES (1);

		CREATE TABLE IF NOT EXISTS stats (
			id			 INTEGER PRIMARY KEY AUTOINCREMENT,
			router_id	 INTEGER REFERENCES routers(id) ON DELETE CASCADE,
			timestamp	 DATETIME DEFAULT CURRENT_TIMESTAMP,
			xdp_drops	 INTEGER DEFAULT 0,
			tc_drops	 INTEGER DEFAULT 0,
			nft_drops	 INTEGER DEFAULT 0,
			pps			 INTEGER DEFAULT 0,
			bps			 INTEGER DEFAULT 0
		);

		CREATE INDEX IF NOT EXISTS idx_stats_router_time
			ON stats(router_id, timestamp);
	`)
	return err
}
