package balance

import (
	"merch-store/internal/domain/balance/transaction_type"
)

type TransactionsReport struct {
	User         string
	Coins        int
	Transactions []Transaction
}

type Transaction struct {
	Type        transaction_type.TransactionType
	Participant string
	Amount      int
}
