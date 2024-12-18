package main

import (
	"log"
	"sync"
)

// executeTransactionsSingle führt die Transaktionen einzeln aus
func executeTransactionsSingle(bankService BankAccountRepository, transactions []Transaction) {
	for _, transaction := range transactions {
		err := bankService.TransferBalance(transaction, 0)
		if err != nil {
			log.Fatalf("Error transferring balance: %v", err)
		}
	}
}

// executeTransactionsGoroutine führt die Transaktionen parallel in Goroutinen aus
func executeTransactionsGoroutine(bankService BankAccountRepository, transactions []Transaction, maxConnections int, delay_transaction float64) {
	var wg sync.WaitGroup

	// Check if bankService is of type SQLBankAccountRepository
	if sqlBankService, ok := bankService.(*SQLBankAccountRepository); ok {
		semaphore := make(chan struct{}, maxConnections*2) // Begrenze die Anzahl der gleichzeitigen Goroutinen auf maxConnections*2
		for _, transaction := range transactions {
			wg.Add(1)
			semaphore <- struct{}{} // Blockiert, wenn das Limit erreicht ist
			go func(tx Transaction, delay float64) {
				defer wg.Done()
				defer func() { <-semaphore }() // Gibt den Platz im Semaphore frei

				err := sqlBankService.TransferBalance(tx, delay)
				if err != nil {
					log.Printf("Error transferring balance: %v", err)
				}
			}(transaction, delay_transaction)
		}
	} else if restBankService, ok := bankService.(*PostgRESTBankAccountRepository); ok {
		semaphore := make(chan struct{}, maxConnections) // Begrenze die Anzahl der gleichzeitigen Goroutinen

		for _, transaction := range transactions {
			wg.Add(1)
			semaphore <- struct{}{} // Blockiert, wenn das Limit erreicht ist
			go func(tx Transaction) {
				defer wg.Done()
				defer func() { <-semaphore }() // Gibt den Platz im Semaphore frei

				err := restBankService.TransferBalance(tx, 0.0)
				if err != nil {
					log.Printf("Error transferring balance: %v", err)
				}
			}(transaction)
		}
	}

	wg.Wait() // Warte, bis alle Goroutinen abgeschlossen sind
}
