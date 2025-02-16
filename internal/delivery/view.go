package delivery

import (
	"merch-store/internal/domain/balance"
	"merch-store/internal/domain/balance/transaction_type"
)

type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type IncomingTransaction struct {
	From   string `json:"fromUser"`
	Amount int    `json:"amount"`
}

type OutgoingTransaction struct {
	To     string `json:"toUser"`
	Amount int    `json:"amount"`
}

type TransactionsReport struct {
	Received []IncomingTransaction `json:"received"`
	Sent     []OutgoingTransaction `json:"sent"`
}

func (v *TransactionsReport) FromModel(model balance.TransactionsReport) TransactionsReport {
	incoming := make([]IncomingTransaction, 0, len(model.Transactions))
	outgoing := make([]OutgoingTransaction, 0, len(model.Transactions))

	for _, transaction := range model.Transactions {
		switch transaction.Type {
		case transaction_type.Income:
			incoming = append(incoming, IncomingTransaction{
				From:   transaction.Participant,
				Amount: transaction.Amount,
			})
		case transaction_type.Outgo:
			outgoing = append(outgoing, OutgoingTransaction{
				To:     transaction.Participant,
				Amount: transaction.Amount,
			})
		}
	}

	*v = TransactionsReport{
		Received: incoming,
		Sent:     outgoing,
	}

	return *v
}
