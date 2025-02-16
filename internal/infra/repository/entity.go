package repository

import (
	"github.com/uptrace/bun"
	"merch-store/internal/domain/auth"
	"merch-store/internal/domain/balance"
	"merch-store/internal/domain/balance/transaction_type"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	Id       int    `bun:"id,pk,autoincrement"`
	Name     string `bun:"name,type:varchar(32),notnull"`
	PassHash string `bun:"passhash,type:varchar(64),notnull"`
	Coins    int    `bun:"coins,type:int,notnull,default:1"`
}

func (e User) ToModel() auth.User {
	return auth.User{
		Id:           e.Id,
		Username:     e.Name,
		PasswordHash: []byte(e.PassHash),
	}
}

type Transaction struct {
	bun.BaseModel `bun:"table:transactions"`

	Id          int64     `bun:"id,pk,autoincrement"`
	SenderId    int       `bun:"sender_id,type:int,notnull"`
	RecipientId int       `bun:"recipient_id,type:int,notnull"`
	Amount      int       `bun:"amount,type:int,notnull"`
	Occurred    time.Time `bun:"occurred,type:timestamp,notnull"`
}

type TransactionPreview struct {
	Amount        int    `bun:"amount"`
	SenderName    string `bun:"sender_name"`
	RecipientName string `bun:"recipient_name"`

	UserId int `bun:"user_id"`
}

func (e TransactionPreview) ToModel(username string) balance.Transaction {
	var (
		transactionType = transaction_type.Income
		participant     = e.SenderName
	)

	if e.SenderName == username {
		transactionType = transaction_type.Outgo
		participant = e.RecipientName
	}

	return balance.Transaction{
		Type:        transactionType,
		Participant: participant,
		Amount:      e.Amount,
	}
}

type Merch struct {
	bun.BaseModel `bun:"table:merch"`

	Id    int    `bun:"id,pk,autoincrement"`
	Name  string `bun:"name,type:varchar(32),unique"`
	Price int    `bun:"price,type:int,notnull"`
}

type Purchase struct {
	bun.BaseModel `bun:"table:purchases"`

	UserId   int       `bun:"user_id,type:int,notnull"`
	MerchId  int       `bun:"merch_id,type:int,notnull"`
	Quantity int       `bun:"quantity,type:int,notnull,default:1"`
	Occurred time.Time `bun:"occurred,type:timestamp,notnull"`
}

type PurchasePreview struct {
	Quantity  int       `bun:"quantity"`
	Occurred  time.Time `bun:"occurred"`
	MerchName string    `bun:"merch_name"`
}
