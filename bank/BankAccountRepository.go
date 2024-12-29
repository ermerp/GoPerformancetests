package main

type BankAccountRepository interface {
	CreateAccount(account Account) error
	DeleteAllAccounts() error
	TransferBalance(transaction Transaction, delay_transaction float64) error
}
