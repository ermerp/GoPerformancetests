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
	fmt.Println("Go:Bank")

	numberOfAccounts := 1000
	numberOfTransactions := 10000

	// Erstelle eine Instanz von Bank, die das BankAccountRepository Interface implementiert
	var bankService BankAccountRepository = PostgRESTBankAccountRepository{}

	// Lösche alle bestehenden Accounts
	err := bankService.DeleteAllAccounts()
	if err != nil {
		log.Fatalf("Error deleting all accounts: %v", err)
	}

	// Importiere Accounts aus einer Datei
	importAccounts(bankService, fmt.Sprintf("BankAccounts%d.txt", numberOfAccounts))

	// Importiere Transaktionen aus einer Datei
	transactions := importTransactions(fmt.Sprintf("BankTransactions%d-%d.txt", numberOfTransactions, numberOfAccounts))

	fmt.Println("File imported.")

	// Führe die Transaktionen durch
	start := time.Now() // Startzeitpunkt

	//executeTransactionsSingle(bankService, transactions)
	executeTransactionsGoroutine(bankService, transactions)

	duration := time.Since(start) // Dauer berechnen
	log.Printf("Go:Bank -  Time: %v", duration)
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
