package main

import (
	"log"
	"sync"
)

// executes the transactions in a single thread
func executeTransactionsSingle(bankService BankAccountRepository, transactions []Transaction) {
	for _, transaction := range transactions {
		err := bankService.TransferBalance(transaction, 0)
		if err != nil {
			log.Fatalf("Error transferring balance: %v", err)
		}
	}
}

// executes the transactions in multiple threads
func executeTransactionsGoroutine(bankService BankAccountRepository, transactions []Transaction, maxConnections int, delay_transaction float64) {
	var wg sync.WaitGroup

	// Check if bankService is of type SQLBankAccountRepository
	if sqlBankService, ok := bankService.(*SQLBankAccountRepository); ok {
		semaphore := make(chan struct{}, maxConnections*2)
		for _, transaction := range transactions {
			wg.Add(1)
			semaphore <- struct{}{}
			go func(tx Transaction, delay float64) {
				defer wg.Done()
				defer func() { <-semaphore }()

				err := sqlBankService.TransferBalance(tx, delay)
				if err != nil {
					log.Printf("Error transferring balance: %v", err)
				}
			}(transaction, delay_transaction)
		}
		// Check if bankService is of type PostgRESTBankAccountRepository
	} else if restBankService, ok := bankService.(*PostgRESTBankAccountRepository); ok {
		semaphore := make(chan struct{}, maxConnections)

		for _, transaction := range transactions {
			wg.Add(1)
			semaphore <- struct{}{}
			go func(tx Transaction) {
				defer wg.Done()
				defer func() { <-semaphore }()

				err := restBankService.TransferBalance(tx, 0.0)
				if err != nil {
					log.Printf("Error transferring balance: %v", err)
				}
			}(transaction)
		}
	}

	wg.Wait() // wait for all goroutines to finish
}
