package imap

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/The-Robin-Hood/bunnymail/internal/logger"
	"github.com/The-Robin-Hood/bunnymail/internal/model"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-message/mail"
)

func (c *Client) FetchMessages(mailbox string, limit int) ([]*model.Message, error) {
	mbox, err := c.client.Select(mailbox, false)
	if err != nil {
		return nil, fmt.Errorf("failed to select mailbox: %w", err)
	}

	if mbox.Messages == 0 {
		return []*model.Message{}, nil
	}

	// Calculate range
	from := uint32(1)
	to := mbox.Messages

	if limit > 0 && int(mbox.Messages) > limit {
		from = mbox.Messages - uint32(limit) + 1
	}

	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- c.client.Fetch(seqset, []imap.FetchItem{
			imap.FetchEnvelope,
			imap.FetchUid,
			"BODY[]",
		}, messages)
	}()

	result := []*model.Message{}

	for msg := range messages {
		message, err := c.parseMessage(msg, mailbox)
		if err != nil {
			logger.Warn("Failed to parse message", "error", err)
			continue
		}
		result = append(result, message)
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("fetch failed: %w", err)
	}

	return result, nil
}

// parseMessage converts IMAP message to our model
func (c *Client) parseMessage(msg *imap.Message, folder string) (*model.Message, error) {
	if msg.Envelope == nil {
		return nil, fmt.Errorf("envelope is nil")
	}

	env := msg.Envelope

	fromName := ""
	fromAddr := "unknown@example.com"

	if len(env.From) > 0 {
		from := env.From[0]
		if from.PersonalName != "" {
			fromName = from.PersonalName
		}
		if from.MailboxName != "" && from.HostName != "" {
			fromAddr = fmt.Sprintf("%s@%s", from.MailboxName, from.HostName)
		}
		if fromName == "" {
			fromName = fromAddr
		}
	}

	toAddrs := ""
	if len(env.To) > 0 {
		addr := env.To[0]
		if addr.MailboxName != "" && addr.HostName != "" {
			toAddrs = fmt.Sprintf("%s@%s", addr.MailboxName, addr.HostName)
		}
	}

	var ccList []string
	for _, addr := range env.Cc {
		if addr.MailboxName != "" && addr.HostName != "" {
			ccList = append(ccList, fmt.Sprintf("%s@%s", addr.MailboxName, addr.HostName))
		}
	}
	ccAddrs := strings.Join(ccList, ", ")

	messageID := env.MessageId
	if messageID == "" {
		messageID = fmt.Sprintf("local-%d", time.Now().UnixNano())
	}

	// Handle date
	receivedAt := time.Now()
	if !env.Date.IsZero() {
		receivedAt = env.Date
	}

	// Extract body (text and HTML)
	bodyText, bodyHTML, hasAttachments := c.extractBody(msg)

	uid := msg.Uid

	message := &model.Message{
		UID:            uid,
		MessageID:      messageID,
		Subject:        env.Subject,
		FromName:       fromName,
		FromAddress:    fromAddr,
		ToAddresses:    toAddrs,
		CcAddresses:    ccAddrs,
		BodyText:       bodyText,
		BodyHTML:       bodyHTML,
		ReceivedAt:     receivedAt.Format(time.RFC3339),
		Folder:         folder,
		IsRead:         false,
		IsStarred:      false,
		HasAttachments: hasAttachments,
	}

	// debug print expect html instead size of html
	logger.Debug(fmt.Sprintf(" Parsed message UID=%d\n Subject=%q\n From=%q\n To=%q\n HasAttachments=%v\n CC=%q\n ReceivedAt=%v\n Folder=%q\n MessageID=%q\n IssRead=%v\n IsStarred=%v\n",
		message.UID,
		message.Subject,
		message.FromAddress,
		message.ToAddresses,
		message.HasAttachments,
		message.CcAddresses,
		message.ReceivedAt,
		message.Folder,
		message.MessageID,
		message.IsRead,
		message.IsStarred,
	))
	logger.Debug(fmt.Sprintf(" BodyText: %.30q\n", message.BodyText))
	logger.Debug(fmt.Sprintf(" BodyHTML: %.30q\n", message.BodyHTML))

	return message, nil
}

// extractBody extracts and decodes text and HTML parts
func (c *Client) extractBody(msg *imap.Message) (string, string, bool) {
	var textPart, htmlPart string
	hasAttachments := false

	// Get the entire message body
	section := msg.GetBody(&imap.BodySectionName{})
	if section == nil {
		return "", "", hasAttachments
	}

	// Create MIME message reader
	mr, err := mail.CreateReader(section)
	if err != nil {
		bodyBytes, _ := io.ReadAll(section)
		return string(bodyBytes), "", hasAttachments
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Warn("Failed to read MIME part", "error", err)
			continue
		}

		switch h := part.Header.(type) {
		case *mail.AttachmentHeader:
			hasAttachments = true
			// filename, _ := h.Filename()

		case *mail.InlineHeader:
			contentType, params, err := h.ContentType()
			if err != nil {
				continue
			}

			partCharset := params["charset"]
			if partCharset == "" {
				partCharset = "utf-8"
			}

			encoding := h.Get("Content-Transfer-Encoding")

			bodyBytes, err := io.ReadAll(part.Body)
			if err != nil {
				continue
			}

			// Decode transfer encoding
			decodedBody := decodeTransferEncoding(bodyBytes, encoding)

			// Convert charset to UTF-8 if needed
			decodedBody = convertCharset(decodedBody, partCharset)

			switch contentType {
			case "text/plain":
				if textPart == "" {
					textPart = decodedBody
				}

			case "text/html":
				if htmlPart == "" {
					htmlPart = decodedBody
				}
			}
		}
	}

	return textPart, htmlPart, hasAttachments
}
