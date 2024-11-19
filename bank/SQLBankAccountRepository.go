package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// SQLBankAccountRepository implementiert das BankAccountRepository Interface
type SQLBankAccountRepository struct {
	pool *pgxpool.Pool
}

// NewSQLBankAccountRepository erstellt ein neues SQLBankAccountRepository
func NewSQLBankAccountRepository(pool *pgxpool.Pool) *SQLBankAccountRepository {
	return &SQLBankAccountRepository{pool: pool}
}

// CreateAccount fügt einen neuen Account in die Datenbank ein
func (r *SQLBankAccountRepository) CreateAccount(account Account) error {
	ctx := context.Background()
	query := "INSERT INTO account (id, balance) VALUES ($1, $2)"
	_, err := r.pool.Exec(ctx, query, account.ID, account.Balance)
	if err != nil {
		return fmt.Errorf("error creating account: %v", err)
	}
	return nil
}

// DeleteAllAccounts löscht alle Accounts aus der Datenbank
func (r *SQLBankAccountRepository) DeleteAllAccounts() error {
	ctx := context.Background()
	query := "DELETE FROM account"
	attempt := 1
	for attempt < MAX_RETRIES {
		_, err := r.pool.Exec(ctx, query)
		if err == nil {
			return nil
		}
		sleepTime := calculateRetryDelay(attempt)
		fmt.Printf("Delete attempt %d failed. Retrying after %v...\n", attempt, sleepTime)
		time.Sleep(sleepTime)
		attempt++
	}
	return fmt.Errorf("delete failed after %d attempts", MAX_RETRIES)
}

// TransferBalance führt eine Transaktion durch, um Guthaben von einem Account auf einen anderen zu übertragen
func (r *SQLBankAccountRepository) TransferBalance(transaction Transaction) error {
	ctx := context.Background()
	attempt := 1
	for attempt < MAX_RETRIES {
		tx, err := r.pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("error beginning transaction: %v", err)
		}

		// Sperre die Zeilen in der richtigen Reihenfolge
		if transaction.From < transaction.To {
			_, err = tx.Exec(ctx, "SELECT 1 FROM account WHERE id = $1 FOR UPDATE", transaction.From)
			retry, _ := handleTransactionError(ctx, tx, err, attempt)
			if retry {
				attempt++
				continue
			}

			_, err = tx.Exec(ctx, "SELECT 1 FROM account WHERE id = $1 FOR UPDATE", transaction.To)
			retry, _ = handleTransactionError(ctx, tx, err, attempt)
			if retry {
				attempt++
				continue
			}
		} else {
			_, err = tx.Exec(ctx, "SELECT 1 FROM account WHERE id = $1 FOR UPDATE", transaction.To)
			retry, _ := handleTransactionError(ctx, tx, err, attempt)
			if retry {
				attempt++
				continue
			}

			_, err = tx.Exec(ctx, "SELECT 1 FROM account WHERE id = $1 FOR UPDATE", transaction.From)
			retry, _ = handleTransactionError(ctx, tx, err, attempt)
			if retry {
				attempt++
				continue
			}
		}

		// Führe die Updates durch
		_, err = tx.Exec(ctx, "UPDATE account SET balance = balance - $1 WHERE id = $2", transaction.Balance, transaction.From)
		retry, _ := handleTransactionError(ctx, tx, err, attempt)
		if retry {
			attempt++
			continue
		}

		_, err = tx.Exec(ctx, "UPDATE account SET balance = balance + $1 WHERE id = $2", transaction.Balance, transaction.To)
		retry, _ = handleTransactionError(ctx, tx, err, attempt)
		if retry {
			attempt++
			continue
		}

		err = tx.Commit(ctx)
		if err != nil {
			retry, _ := handleTransactionError(ctx, tx, err, attempt)
			if retry {
				attempt++
				continue
			}
		}

		return nil
	}
	return fmt.Errorf("transfer failed after %d attempts", MAX_RETRIES)
}

func handleTransactionError(ctx context.Context, tx pgx.Tx, err error, attempt int) (bool, error) {
	if err != nil {
		tx.Rollback(ctx)
		sleepTime := calculateRetryDelay(attempt)
		fmt.Printf("Transfer attempt %d failed. Retrying after %v...\n", attempt, sleepTime)
		time.Sleep(sleepTime)
		return true, nil
	}
	return false, nil
}
