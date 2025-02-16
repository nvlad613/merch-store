package store

import "time"

type Purchase struct {
	Username    string
	ProductName string
	Timestamp   time.Time
	Quantity    int
}

type InventoryItem struct {
	ProductName string
	Quantity    int
}
