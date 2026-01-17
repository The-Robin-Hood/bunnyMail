package store

import (
	"database/sql"
	"fmt"

	"github.com/The-Robin-Hood/bunnymail/internal/model"
)

type AccountStore struct {
	db *DB
}

func NewAccountStore(db *DB) *AccountStore {
	return &AccountStore{db: db}
}

func (s *AccountStore) Create(acc *model.Account) error {
	query := `
		INSERT INTO accounts (
			name, 
			email, 
			password, 
			remember_password,
			imap_host, 
			imap_port, 
			imap_username, 
			imap_password, 
			imap_security, 
			imap_auth_type,
			smtp_host, 
			smtp_port, 
			smtp_username, 
			smtp_password, 
			smtp_security, 
			smtp_auth_type,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))
	`

	imapPassword := ""
	if acc.RememberPassword {
		imapPassword = acc.IMAPPassword
	}

	result, err := s.db.Exec(
		query,
		acc.Name,
		acc.Email,
		acc.Password,
		acc.RememberPassword,
		acc.IMAPHost,
		acc.IMAPPort,
		acc.IMAPUsername,
		imapPassword,
		acc.IMAPSecurity,
		acc.IMAPAuthType,
		acc.SMTPHost,
		acc.SMTPPort,
		acc.SMTPUsername,
		acc.SMTPPassword,
		acc.SMTPSecurity,
		acc.SMTPAuthType,
	)

	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	acc.ID = id
	return nil
}

func (s *AccountStore) List() ([]*model.Account, error) {
	var accounts []*model.Account
	query := `SELECT * FROM accounts ORDER BY created_at DESC`

	err := s.db.Select(&accounts, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}

	return accounts, nil
}

func (s *AccountStore) GetByID(id int64) (*model.Account, error) {
	var acc model.Account
	query := `SELECT * FROM accounts WHERE id = ?`

	err := s.db.Get(&acc, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("account not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &acc, nil
}

func (s *AccountStore) Delete(id int64) error {
	query := `DELETE FROM accounts WHERE id = ?`
	_, err := s.db.Exec(query, id)
	return err
}
