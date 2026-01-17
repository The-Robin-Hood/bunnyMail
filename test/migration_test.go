package test

import (
	"database/sql"
	"testing"

	"github.com/The-Robin-Hood/bunnymail/internal/store"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *store.DB {
	t.Helper()
	db, err := store.InitializeDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	return db
}

func TestRunMigrations(t *testing.T) {
	db := setupTestDB(t)

	if err := db.RunMigrations(); err != nil {
		t.Fatalf("RunMigrations failed: %v", err)
	}

	// Verify tables exist
	assertTableExists(t, db, "accounts")
	assertTableExists(t, db, "messages")

	// Verify critical columns exist
	assertColumnExists(t, db, "messages", "uid")
	assertColumnExists(t, db, "messages", "cc_addresses")
	assertColumnExists(t, db, "messages", "has_attachments")
	assertColumnExists(t, db, "messages", "folder")
}

func assertTableExists(t *testing.T, db *store.DB, table string) {
	t.Helper()

	var name string
	err := db.QueryRow(
		`SELECT name FROM sqlite_master WHERE type='table' AND name=?`,
		table,
	).Scan(&name)

	if err != nil {
		t.Fatalf("table %s does not exist: %v", table, err)
	}
}

func assertColumnExists(t *testing.T, db *store.DB, table, column string) {
	t.Helper()

	rows, err := db.Query(`PRAGMA table_info(` + table + `);`)
	if err != nil {
		t.Fatalf("failed to inspect table %s: %v", table, err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			cid        int
			name       string
			colType    string
			notNull    int
			defaultVal sql.NullString
			pk         int
		)

		if err := rows.Scan(&cid, &name, &colType, &notNull, &defaultVal, &pk); err != nil {
			t.Fatalf("failed scanning table info: %v", err)
		}

		if name == column {
			return
		}
	}

	t.Fatalf("column %s.%s does not exist", table, column)
}
