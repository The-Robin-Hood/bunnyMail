package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/The-Robin-Hood/bunnymail/internal/store"
)

type App struct {
	DB           *store.DB
	AccountStore *store.AccountStore
	MessageStore *store.MessageStore
}

func InitializeApp() (*App, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home dir: %w", err)
	}

	appDir := filepath.Join(home, ".bunnymail")
	if err := os.MkdirAll(appDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create app dir: %w", err)
	}

	dbPath := filepath.Join(appDir, "mail.db")

	db, err := store.InitializeDB(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.RunMigrations(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	accountStore := store.NewAccountStore(db)
	messageStore := store.NewMessageStore(db)

	return &App{
		DB:           db,
		AccountStore: accountStore,
		MessageStore: messageStore,
	}, nil
}

func (a *App) TerminateApp() error {
	return a.DB.CloseDB()
}
