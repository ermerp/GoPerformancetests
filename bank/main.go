package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Account struct {
	ID      string  `json:"account_id"`
	Balance float64 `json:"balance"`
}

type Transaction struct {
	From    string  `json:"from"`
	To      string  `json:"to"`
	Balance float64 `json:"balance"`
}

func main() {
	// Retrieve the environment variables
	interfaceType := os.Getenv("INTERFACE_TYPE")
	if interfaceType == "" {
		interfaceType = "PGX"
	}
	algorithm := os.Getenv("ALGORITHM")
	if algorithm == "" {
		algorithm = "GOROUTINE"
	}
	numberOfAccounts := os.Getenv("NUMBER_Of_ACCOUNTS")
	if numberOfAccounts == "" {
		numberOfAccounts = "10"
	}
	numberOfTransactions := os.Getenv("NUMBER_OF_TRANSACTIONS")
	if numberOfTransactions == "" {
		numberOfTransactions = "100"
	}

	maxConn := os.Getenv("MAX_CONNECTIONS")
	if maxConn == "" {
		maxConn = "10"
	}
	maxConnections, err := strconv.Atoi(maxConn)
	if err != nil {
		log.Fatalf("Invalid max connections: %d", maxConnections)
	}

	fmt.Printf("Go:Bank - Interface: %s, Algorithm: %s, Max Connections: %s, Number of Accounts: %s, Number of Transactions: %s\n",
		interfaceType, algorithm, maxConn, numberOfAccounts, numberOfTransactions)

	var bankService BankAccountRepository
	switch interfaceType {
	case "PGX":
		bankService = NewSQLBankAccountRepository(int32(maxConnections))
		defer bankService.(*SQLBankAccountRepository).pool.Close()
	case "REST":
		bankService = NewPostgRESTBankAccountRepository()
	default:
		log.Fatalf("Unknown interface type: %s", interfaceType)
	}

	// Clean up the database
	err = bankService.DeleteAllAccounts()
	if err != nil {
		log.Fatalf("Error deleting all accounts: %v", err)
	}

	// Importiere Accounts aus einer Datei
	importAccounts(bankService, fmt.Sprintf("bankData/BankAccounts%s.txt", numberOfAccounts))

	// Importiere Transaktionen aus einer Datei
	transactions := importTransactions(fmt.Sprintf("bankData/BankTransactions%s-%s.txt", numberOfTransactions, numberOfAccounts))

	fmt.Println("File imported.")

	start := time.Now()
	// Führe die Transaktionen durch
	switch algorithm {
	case "SINGLE":
		executeTransactionsSingle(bankService, transactions)
	case "GOROUTINE":

		executeTransactionsGoroutine(bankService, transactions, maxConnections)
	default:
		log.Fatalf("Unknown algorithm: %s", algorithm)
	}
	duration := time.Since(start)

	log.Printf("Go:Bank - Time: %v", duration)
}

// importAccounts liest die Datei ein und schreibt die Accounts in die Datenbank
func importAccounts(bankService BankAccountRepository, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var accounts []Account
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Teile die Zeile in ID und Balance auf
		parts := strings.SplitN(line, ", ", 2)
		if len(parts) != 2 {
			log.Fatalf("Invalid line format: %s", line)
		}

		// Ersetze das Komma durch einen Punkt im Balance-Teil
		balanceStr := strings.Replace(parts[1], ",", ".", 1)
		balance, err := strconv.ParseFloat(balanceStr, 64)
		if err != nil {
			log.Fatalf("Invalid balance format: %s", balanceStr)
		}

		account := Account{
			ID:      parts[0],
			Balance: balance,
		}
		accounts = append(accounts, account)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for _, acc := range accounts {
		err := bankService.CreateAccount(acc)
		if err != nil {
			log.Fatalf("Error creating account: %v", err)
		}
	}
}

// importTransactions liest die Datei ein und speichert jede Zeile als Transaction in einer Liste
func importTransactions(fileName string) []Transaction {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var transactions []Transaction
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Teile die Zeile in From, To und Balance auf
		parts := strings.SplitN(line, ", ", 3)
		if len(parts) != 3 {
			log.Fatalf("Invalid line format: %s", line)
		}

		// Ersetze das Komma durch einen Punkt im Balance-Teil
		balanceStr := strings.Replace(parts[2], ",", ".", 1)
		balance, err := strconv.ParseFloat(balanceStr, 64)
		if err != nil {
			log.Fatalf("Invalid balance format: %s", balanceStr)
		}

		transaction := Transaction{
			From:    parts[0],
			To:      parts[1],
			Balance: balance,
		}
		transactions = append(transactions, transaction)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return transactions
}
