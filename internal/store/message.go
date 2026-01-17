package store

import (
	"database/sql"
	"fmt"

	"github.com/The-Robin-Hood/bunnymail/internal/model"
)

type MessageStore struct {
	db *DB
}

func NewMessageStore(db *DB) *MessageStore {
	return &MessageStore{db: db}
}

// Create inserts a new message (or ignores if exists)
func (s *MessageStore) Create(msg *model.Message) error {
	query := `
		INSERT INTO messages (
			account_id,
			message_id,
			uid,
			folder,
			subject,
			from_name,
			from_address,
			to_addresses,
			cc_addresses,
			body_text,
			body_html,
			received_at,
			is_read,
			is_starred,
			has_attachments,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))
	`

	result, err := s.db.Exec(query,
		msg.AccountID,
		msg.MessageID,
		msg.UID,
		msg.Folder,
		msg.Subject,
		msg.FromName,
		msg.FromAddress,
		msg.ToAddresses,
		msg.CcAddresses,
		msg.BodyText,
		msg.BodyHTML,
		msg.ReceivedAt,
		msg.IsRead,
		msg.IsStarred,
		msg.HasAttachments,
	)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to fetch inserted id: %w", err)
	}

	msg.ID = id
	return nil
}

// List returns messages for an account
func (s *MessageStore) List(accountID int64, limit int) ([]*model.Message, error) {
	var messages []*model.Message
	query := `
		SELECT * FROM messages 
		WHERE account_id = ?
		ORDER BY received_at DESC
		LIMIT ?
	`

	err := s.db.Select(&messages, query, accountID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}

	return messages, nil
}

// GetByID retrieves a message by ID
func (s *MessageStore) GetByID(id int64) (*model.Message, error) {
	var msg model.Message
	query := `SELECT * FROM messages WHERE id = ?`

	err := s.db.Get(&msg, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("message not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return &msg, nil
}

// MarkAsRead marks a message as read
func (s *MessageStore) MarkAsRead(id int64) error {
	query := `UPDATE messages SET is_read = TRUE WHERE id = ?`
	_, err := s.db.Exec(query, id)
	return err
}

// Count returns total messages for an account
func (s *MessageStore) Count(accountID int64) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM messages WHERE account_id = ?`
	err := s.db.Get(&count, query, accountID)
	return count, err
}
