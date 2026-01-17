package model

type Message struct {
	ID        int64 `db:"id" json:"id"`
	AccountID int64 `db:"account_id" json:"account_id"`

	MessageID string `db:"message_id" json:"message_id"` // RFC 5322 Message-ID
	UID       uint32 `db:"uid" json:"uid"`               // IMAP UID
	Folder    string `db:"folder" json:"folder"`         // INBOX, Sent, Drafts

	Subject string `db:"subject" json:"subject"`

	FromName    string `db:"from_name" json:"from_name"`
	FromAddress string `db:"from_address" json:"from_address"`

	ToAddresses string `db:"to_addresses" json:"to_addresses"` // JSON array
	CcAddresses string `db:"cc_addresses" json:"cc_addresses"` // JSON array (optional but small)

	BodyText string `db:"body_text" json:"body_text"`
	BodyHTML string `db:"body_html" json:"body_html"`

	ReceivedAt     string `db:"received_at" json:"received_at"`
	IsRead         bool   `db:"is_read" json:"is_read"`
	IsStarred      bool   `db:"is_starred" json:"is_starred"`
	HasAttachments bool   `db:"has_attachments" json:"has_attachments"`

	CreatedAt string `db:"created_at" json:"created_at"`
}
