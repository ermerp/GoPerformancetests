package main

// BankAccountRepository definiert die Methoden, die von der Bank-Implementierung bereitgestellt werden müssen
type BankAccountRepository interface {
	CreateAccount(account Account) error
	DeleteAllAccounts() error
	TransferBalance(transaction Transaction, delay_transaction float64) error
}
