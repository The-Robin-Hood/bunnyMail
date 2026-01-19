package main

import (
	"context"
	"fmt"

	"github.com/The-Robin-Hood/bunnymail/internal/app"
	"github.com/The-Robin-Hood/bunnymail/internal/logger"
	"github.com/The-Robin-Hood/bunnymail/internal/mail/imap"
	"github.com/The-Robin-Hood/bunnymail/internal/model"
)

type App struct {
	ctx     context.Context
	appCore *app.App
}

func BunnyMailApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	appCore, err := app.InitializeApp()
	if err != nil {
		logger.Error("Failed to initialize app", "error", err)
		return
	}

	a.appCore = appCore
}

func (a *App) G_TestConnection(accountID int64) error {
	if a.appCore == nil {
		return fmt.Errorf("app not initialized")
	}

	account, err := a.appCore.AccountStore.GetByID(accountID)
	if err != nil {
		return err
	}

	return imap.TestConnection(
		account.IMAPHost,
		account.IMAPPort,
		account.Email,
		account.Password,
	)
}

func (a *App) G_GetAccounts() ([]model.Account, error) {
	if a.appCore == nil {
		return nil, fmt.Errorf("app not initialized")
	}

	accounts, err := a.appCore.AccountStore.List()
	if err != nil {
		logger.Error("Error retrieving stored accounts", "error", err)
		return nil, err
	}

	result := make([]model.Account, len(accounts))
	for i, acc := range accounts {
		result[i] = *acc
	}

	return result, nil
}

// AddAccount adds a new email account
func (a *App) G_AddAccount(account model.Account) error {
	if a.appCore == nil {
		return fmt.Errorf("app not initialized")
	}

	return a.appCore.AccountStore.Create(&account)
}

// SyncAccount syncs emails from an account
func (a *App) G_SyncAccount(accountID int64, limit int) (int, error) {
	if a.appCore == nil {
		return 0, fmt.Errorf("app not initialized")
	}

	account, err := a.appCore.AccountStore.GetByID(accountID)
	if err != nil {
		return 0, err
	}

	// Connect to IMAP
	client, err := imap.InitiateIMAPSession(
		account.IMAPHost,
		account.IMAPPort,
		account.IMAPUsername,
		account.IMAPPassword,
		imap.IMAPSecurity(account.IMAPSecurity),
	)
	if err != nil {
		return 0, err
	}
	defer client.CloseIMAPSession()

	messages, err := client.FetchMessages("INBOX", limit)
	if err != nil {
		return 0, err
	}

	saved := 0
	for _, msg := range messages {
		msg.AccountID = accountID
		if err := a.appCore.MessageStore.Create(msg); err == nil {
			saved++
		}
	}

	return saved, nil
}

// GetMessagesByAccount returns messages from specific account
func (a *App) G_GetMessagesByAccount(accountID int64, limit int) ([]model.Message, error) {
	if a.appCore == nil {
		return nil, fmt.Errorf("app not initialized")
	}

	messages, err := a.appCore.MessageStore.List(accountID, limit)
	if err != nil {
		return nil, err
	}

	result := make([]model.Message, len(messages))
	for i, msg := range messages {
		result[i] = *msg
	}

	return result, nil
}

// GetMessage returns a specific message
// func (a *App) GetMessage(messageID int64) (*model.Message, error) {
// 	if a.appCore == nil {
// 		return nil, fmt.Errorf("app not initialized")
// 	}

// 	return a.appCore.MessageStore.GetByID(messageID)
// }

// // MarkAsRead marks a message as read
// func (a *App) MarkAsRead(messageID int64) error {
// 	if a.appCore == nil {
// 		return fmt.Errorf("app not initialized")
// 	}

// 	return a.appCore.MessageStore.MarkAsRead(messageID)
// }

// // DeleteAccount deletes an account
// func (a *App) DeleteAccount(accountID int64) error {
// 	if a.appCore == nil {
// 		return fmt.Errorf("app not initialized")
// 	}

// 	return a.appCore.AccountStore.Delete(accountID)
// }
