package imap

import (
	"fmt"

	"github.com/The-Robin-Hood/bunnymail/internal/logger"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type Client struct {
	client *client.Client
}

type IMAPSecurity string

const (
	IMAPSecurityTLS      IMAPSecurity = "ssl_tls"
	IMAPSecuritySTARTTLS IMAPSecurity = "start_tls"
	IMAPSecurityNone     IMAPSecurity = "none"
)

func InitiateIMAPSession(host string, port int, username, password string, security IMAPSecurity) (*Client, error) {

	addr := fmt.Sprintf("%s:%d", host, port)

	var c *client.Client
	var err error

	switch security {

	case IMAPSecurityTLS:
		c, err = client.DialTLS(addr, nil)

	case IMAPSecuritySTARTTLS:
		c, err = client.Dial(addr)
		if err == nil {
			if err = c.StartTLS(nil); err != nil {
				c.Logout()
				return nil, fmt.Errorf("STARTTLS failed: %w", err)
			}
		}

	case IMAPSecurityNone:
		c, err = client.Dial(addr)

	default:
		return nil, fmt.Errorf("unsupported IMAP security: %s", security)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	if err := c.Login(username, password); err != nil {
		c.Logout()
		return nil, fmt.Errorf("login failed: %w", err)
	}

	return &Client{client: c}, nil
}

func (c *Client) CloseIMAPSession() error {
	if c.client != nil {
		return c.client.Logout()
	}
	return nil
}

func (c *Client) ListMailboxes() ([]string, error) {
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)

	go func() {
		done <- c.client.List("", "*", mailboxes)
	}()

	var names []string
	for m := range mailboxes {
		names = append(names, m.Name)
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("failed to list mailboxes: %w", err)
	}

	return names, nil
}

func TestConnection(host string, port int, username, password string) error {
	client, err := InitiateIMAPSession(host, port, username, password, IMAPSecurityTLS)
	if err != nil {
		return err
	}
	defer client.CloseIMAPSession()

	test, _ := client.ListMailboxes()

	for i := range test {
		logger.Info("Found mailbox: %s", test[i])
	}
	return nil
}
