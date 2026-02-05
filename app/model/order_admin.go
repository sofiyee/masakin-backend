package model

import "time"

type OrderAdminList struct {
	ID           int       `json:"id"`
	CustomerName string    `json:"customer_name"`
	OrderType    string    `json:"order_type"`
	Status       string    `json:"status"`
	TotalPrice   int       `json:"total_price"`
	CreatedAt    time.Time `json:"created_at"`
}

