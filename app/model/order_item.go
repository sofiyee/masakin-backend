package model

import "time"

type OrderItem struct {
	ID        int
	OrderID  int
	MenuID   int
	OrderDate time.Time
	Quantity int
}

type RecentOrder struct {
	ID           int       `json:"id"`
	CustomerName string    `json:"customer_name"`
	OrderType    string    `json:"order_type"`
	Status       string    `json:"status"`
	TotalPrice   int       `json:"total_price"`
	CreatedAt    time.Time `json:"created_at"`
}

// Response struct untuk API
type RecentOrdersResponse struct {
	Success bool           `json:"success"`
	Data    []RecentOrder  `json:"data"`
	Message string         `json:"message,omitempty"`
}