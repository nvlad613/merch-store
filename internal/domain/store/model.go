package store

import "time"

type Purchase struct {
	ProductName string
	Timestamp   time.Time
	Quantity    int
}

type InventoryItem struct {
	ProductName string
	Quantity    int
}
