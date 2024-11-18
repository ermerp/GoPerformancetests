package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	MAX_RETRIES = 100
)

// PostgRESTBankAccountRepository implementiert das BankAccountRepository Interface
type PostgRESTBankAccountRepository struct{}

// CreateAccount sendet eine HTTP-POST-Anfrage an die PostgREST-API, um einen neuen Account zu erstellen
func (b PostgRESTBankAccountRepository) CreateAccount(account Account) error {
	url := "http://localhost:3000/rpc/create_account"

	// JSON-Daten erstellen
	jsonData, err := json.Marshal(map[string]interface{}{
		"account_id": account.ID,
		"balance":    account.Balance,
	})
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	// HTTP-POST-Anfrage senden
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending POST request: %v", err)
	}
	defer resp.Body.Close()

	// Überprüfen, ob die Anfrage erfolgreich war
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("received non-200 response: %s, body: %s", resp.Status, string(body))
	}

	return nil
}

// DeleteAllAccounts sendet eine HTTP-POST-Anfrage an die PostgREST-API, um alle Accounts zu löschen
func (b PostgRESTBankAccountRepository) DeleteAllAccounts() error {
	url := "http://localhost:3000/rpc/delete_all_accounts"

	// Leere JSON-Daten erstellen
	jsonData, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	// HTTP-POST-Anfrage senden
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending POST request: %v", err)
	}
	defer resp.Body.Close()

	// Überprüfen, ob die Anfrage erfolgreich war
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("received non-200 response: %s, body: %s", resp.Status, string(body))
	}

	return nil
}

// TransferBalance sendet eine HTTP-POST-Anfrage an die PostgREST-API, um eine Transaktion durchzuführen
func (b PostgRESTBankAccountRepository) TransferBalance(transaction Transaction) error {
	url := "http://localhost:3000/rpc/transfer_balance"

	// JSON-Daten erstellen
	jsonData, err := json.Marshal(map[string]interface{}{
		"from_id": transaction.From,
		"to_id":   transaction.To,
		"amount":  transaction.Balance,
	})
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	attempt := 1
	for attempt < MAX_RETRIES {
		// HTTP-POST-Anfrage senden
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("error sending POST request: %v", err)
		}

		// Überprüfen, ob die Anfrage erfolgreich war
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close() // Verbindung sofort schließen

		if resp.StatusCode == http.StatusOK {
			return nil
		}

		switch resp.StatusCode {
		case http.StatusInternalServerError:
			attempt++
			sleepTime := calculateRetryDelay(attempt)
			fmt.Printf("Deadlock detected. Transfer attempt %d failed: %s -> %s, amount: %.2f. Retrying after %v... Response body: %s\n", attempt, transaction.From, transaction.To, transaction.Balance, sleepTime, string(body))
			time.Sleep(sleepTime)
		case http.StatusGatewayTimeout:
			attempt++
			sleepTime := calculateRetryDelay(attempt)
			fmt.Printf("Gateway Timeout. Transfer attempt %d failed: %s -> %s, amount: %.2f. Retrying after %v... Response body: %s\n", attempt, transaction.From, transaction.To, transaction.Balance, sleepTime, string(body))
			time.Sleep(sleepTime)
		default:
			sleepTime := calculateRetryDelay(attempt)
			fmt.Printf("Error!!! Transfer failed - Code: %d Retrying after %v ... Response body: %s\n", resp.StatusCode, sleepTime, string(body))
			time.Sleep(sleepTime)
		}
	}

	fmt.Printf("Error!!! Transfer failed after %d attempts due to persistent deadlocks: %s -> %s, amount: %.2f\n", MAX_RETRIES, transaction.From, transaction.To, transaction.Balance)
	return fmt.Errorf("transfer failed after %d attempts", MAX_RETRIES)
}

// calculateRetryDelay berechnet die Verzögerung basierend auf der Anzahl der Versuche
func calculateRetryDelay(attempt int) time.Duration {
	return time.Duration(attempt*1000) * time.Millisecond
}
