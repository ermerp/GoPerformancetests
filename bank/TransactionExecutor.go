package main

import (
	"log"
	"sync"
)

// executeTransactionsSingle führt die Transaktionen einzeln aus
func executeTransactionsSingle(bankService BankAccountRepository, transactions []Transaction) {
	for _, transaction := range transactions {
		err := bankService.TransferBalance(transaction)
		if err != nil {
			log.Fatalf("Error transferring balance: %v", err)
		}
	}
}

// executeTransactionsGoroutine führt die Transaktionen parallel in Goroutinen aus
func executeTransactionsGoroutine(bankService BankAccountRepository, transactions []Transaction) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10) // Begrenze die Anzahl der gleichzeitigen Goroutinen auf 10

	for _, transaction := range transactions {
		wg.Add(1)
		semaphore <- struct{}{} // Blockiert, wenn das Limit erreicht ist
		go func(tx Transaction) {
			defer wg.Done()
			defer func() { <-semaphore }() // Gibt den Platz im Semaphore frei

			err := bankService.TransferBalance(tx)
			if err != nil {
				log.Printf("Error transferring balance: %v", err)
			}
		}(transaction)
	}
	wg.Wait() // Warte, bis alle Goroutinen abgeschlossen sind
}
