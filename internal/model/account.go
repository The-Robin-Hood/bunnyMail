package model

type Account struct {
	ID               int64  `db:"id" json:"id"`
	Name             string `db:"name" json:"name"`
	Email            string `db:"email" json:"email"`
	Password         string `db:"password" json:"password"`
	RememberPassword bool   `db:"remember_password" json:"remember_password"`

	// IMAP settings
	IMAPHost     string `db:"imap_host" json:"imap_host"`
	IMAPPort     int    `db:"imap_port" json:"imap_port"`
	IMAPUsername string `db:"imap_username" json:"imap_username"`
	IMAPPassword string `db:"imap_password" json:"imap_password"`
	IMAPSecurity string `db:"imap_security" json:"imap_security"`
	IMAPAuthType string `db:"imap_auth_type" json:"imap_auth_type"`

	// SMTP settings
	SMTPHost     string `db:"smtp_host" json:"smtp_host"`
	SMTPPort     int    `db:"smtp_port" json:"smtp_port"`
	SMTPUsername string `db:"smtp_username" json:"smtp_username"`
	SMTPPassword string `db:"smtp_password" json:"smtp_password"`
	SMTPSecurity string `db:"smtp_security" json:"smtp_security"`
	SMTPAuthType string `db:"smtp_auth_type" json:"smtp_auth_type"`

	LastSyncAt *string `db:"last_sync_at" json:"last_sync_at"`
	CreatedAt  string  `db:"created_at" json:"created_at"`
}
