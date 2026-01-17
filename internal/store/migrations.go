package store

import "fmt"

func (db *DB) RunMigrations() error {
	migrations := []string{
		// Accounts table
		`CREATE TABLE IF NOT EXISTS accounts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password TEXT,
			remember_password BOOLEAN NOT NULL DEFAULT FALSE,

			imap_host TEXT,
			imap_port INTEGER,
			imap_username TEXT,
			imap_password TEXT,
			imap_security TEXT,
			imap_auth_type TEXT,

			smtp_host TEXT,
			smtp_port INTEGER,
			smtp_username TEXT,
			smtp_password TEXT,
			smtp_security TEXT,
			smtp_auth_type TEXT,

			last_sync_at DATETIME,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,

		// Messages table
		`CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			account_id INTEGER NOT NULL,

			message_id TEXT NOT NULL,       -- RFC 5322 Message-ID
			uid INTEGER NOT NULL,           -- IMAP UID
			folder TEXT NOT NULL DEFAULT 'INBOX',

			subject TEXT,

			from_name TEXT,
			from_address TEXT NOT NULL,

			to_addresses TEXT,              -- JSON array
			cc_addresses TEXT,              -- JSON array

			body_text TEXT,
			body_html TEXT,

			received_at DATETIME NOT NULL,

			is_read BOOLEAN NOT NULL DEFAULT FALSE,
			is_starred BOOLEAN NOT NULL DEFAULT FALSE,
			has_attachments BOOLEAN NOT NULL DEFAULT FALSE,

			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

			FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
			UNIQUE (account_id, message_id),
			UNIQUE (account_id, folder, uid)
		);`,

		// Indexes
		`CREATE INDEX IF NOT EXISTS idx_messages_account_id
			ON messages(account_id);`,

		`CREATE INDEX IF NOT EXISTS idx_messages_received_at
			ON messages(received_at);`,

		`CREATE INDEX IF NOT EXISTS idx_messages_account_folder
			ON messages(account_id, folder);`,
	}

	for i, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration %d failed: %w", i, err)
		}
	}

	return nil
}
