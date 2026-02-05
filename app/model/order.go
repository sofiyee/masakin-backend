package model

import "time"

type Order struct {
	ID         int
	CustomerID int
	OrderType  string
	StartDate  time.Time
	EndDate    time.Time
	Status     string
	TotalPrice int
	CreatedAt  time.Time
}

type OrderAdminView struct {
	ID         int    `json:"id"`
	CustomerID int    `json:"customer_id"`
	OrderType  string `json:"order_type"`
	Status     string `json:"status"`
	TotalPrice int    `json:"total_price"`
	CreatedAt  string `json:"created_at"`
}

type OrderCourierView struct {
	OrderID      int       `json:"order_id"`
	CustomerName string    `json:"customer_name"`
	Region       string    `json:"region"`
	FullAddress  string    `json:"full_address"`
	OrderType    string    `json:"order_type"`
	CreatedAt    time.Time `json:"created_at"`
}

type UnpaidOrderAdminView struct {
	ID           int    `json:"id"`
	CustomerName string `json:"customer_name"`
	OrderType    string `json:"order_type"`
	TotalPrice   int    `json:"total_price"`
}

type RecentOrderCustomerView struct {
	OrderID    int    `json:"order_id"`
	OrderType  string `json:"order_type"`
	Status     string `json:"status"`
	TotalPrice int    `json:"total_price"`
	CreatedAt  string `json:"created_at"`
}