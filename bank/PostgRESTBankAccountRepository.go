package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type PostgRESTBankAccountRepository struct {
	url string
}

func NewPostgRESTBankAccountRepository() *PostgRESTBankAccountRepository {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	newUrl := fmt.Sprintf("http://%s:3000/rpc/", dbHost)
	return &PostgRESTBankAccountRepository{url: newUrl}
}

func (r *PostgRESTBankAccountRepository) CreateAccount(account Account) error {

	jsonData, err := json.Marshal(map[string]interface{}{
		"account_id": account.ID,
		"balance":    account.Balance,
	})
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("%screate_account", r.url), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("received non-200 response: %s, body: %s", resp.Status, string(body))
	}

	return nil
}

func (r *PostgRESTBankAccountRepository) DeleteAllAccounts() error {

	jsonData, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	attempt := 1
	for attempt < MAX_RETRIES {

		resp, err := http.Post(fmt.Sprintf("%sdelete_all_accounts", r.url), "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("error sending POST request: %v", err)
		}

		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return nil
		}

		switch resp.StatusCode {
		case http.StatusInternalServerError:
			attempt++
			sleepTime := calculateRetryDelay(attempt)
			fmt.Printf("Deadlock detected. Delete attempt %d failed. Retrying after %v... Response body: %s\n", attempt, sleepTime, string(body))
			time.Sleep(sleepTime)
		case http.StatusGatewayTimeout:
			attempt++
			sleepTime := calculateRetryDelay(attempt)
			fmt.Printf("Gateway Timeout. Delete attempt %d failed. Retrying after %v... Response body: %s\n", attempt, sleepTime, string(body))
			time.Sleep(sleepTime)
		default:
			sleepTime := calculateRetryDelay(attempt)
			fmt.Printf("Error!!! Delete failed - Code: %d Retrying after %v ... Response body: %s\n", resp.StatusCode, sleepTime, string(body))
			time.Sleep(sleepTime)
		}
	}

	fmt.Printf("Error!!! Delete failed after %d attempts due to persistent issues.\n", MAX_RETRIES)
	return fmt.Errorf("delete failed after %d attempts", MAX_RETRIES)
}

func (r *PostgRESTBankAccountRepository) TransferBalance(transaction Transaction, delay_transaction float64) error {

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

		resp, err := http.Post(fmt.Sprintf("%stransfer_balance", r.url), "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("error sending POST request: %v", err)
		}

		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

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
