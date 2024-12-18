package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// SQLBankAccountRepository implementiert das BankAccountRepository Interface
type SQLBankAccountRepository struct {
	pool *pgxpool.Pool
}

// NewSQLBankAccountRepository erstellt ein neues SQLBankAccountRepository
func NewSQLBankAccountRepository(maxConns int32) *SQLBankAccountRepository {
	return &SQLBankAccountRepository{pool: initializePool(maxConns)}
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
func (r *SQLBankAccountRepository) TransferBalance(transaction Transaction, delay_transaction float64) error {
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

		// Führe pg_sleep in der Transaktion aus, um eine Verzögerung in der Datenbank zu erzeugen
		if delay_transaction > 0 {
			_, err = tx.Exec(ctx, fmt.Sprintf("SELECT pg_sleep(%f);", delay_transaction))
			if err != nil {
				tx.Rollback(ctx)
				return fmt.Errorf("error during pg_sleep: %v", err)
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

func initializePool(maxConns int32) *pgxpool.Pool {
	// PostgreSQL Verbindungsinformationen
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	psqlInfo := fmt.Sprintf("postgres://myuser:mypassword@%s:5432/mydatabase?sslmode=disable", dbHost)

	// Erstelle einen pgxpool
	config, err := pgxpool.ParseConfig(psqlInfo)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}

	// Setze die maximale Anzahl von Verbindungen
	config.MaxConns = maxConns

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	return pool
}
