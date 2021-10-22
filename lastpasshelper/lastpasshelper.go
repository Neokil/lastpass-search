package lastpasshelper

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ansd/lastpass-go"
)

const CACHE_FILE_LOCATION = "/tmp/lastpass-accounts.json"

type Account struct {
	ID       string
	Name     string
	Username string
	Password string
	URL      string
	Notes    string
}

func UpdateAccounts(username string, password string, otpCode string) error {
	opts := []lastpass.ClientOption{}
	if otpCode != "" {
		opts = append(opts, lastpass.WithOneTimePassword(otpCode))
	}
	cl, err := lastpass.NewClient(context.Background(), username, password, opts...)
	if err != nil {
		return fmt.Errorf("Failed to log into lastpass-account: %w", err)
	}

	accounts, err := cl.Accounts(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to log retrieve accounts: %w", err)
	}

	_ = cl.Logout(context.Background())

	b, err := json.Marshal(accounts)
	if err != nil {
		return fmt.Errorf("Failed to marshal accounts: %w", err)
	}

	err = os.WriteFile(CACHE_FILE_LOCATION, b, 0600)
	if err != nil {
		return fmt.Errorf("Failed to write accounts to file: %w", err)
	}
	return nil
}

func GetAccounts() ([]Account, error) {
	b, err := os.ReadFile(CACHE_FILE_LOCATION)
	if err != nil {
		return nil, fmt.Errorf("Failed to read accounts from file: %w", err)
	}
	accounts := []Account{}
	err = json.Unmarshal(b, &accounts)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal accounts from file: %w", err)
	}
	return accounts, nil
}
