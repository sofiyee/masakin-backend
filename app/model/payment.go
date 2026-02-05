package model

import "time"

type Payment struct {
	ID         int
	OrderID    int
	ProofImage string
	VerifiedBy int
	CreatedAt  time.Time
}


type PaymentAdminView struct {
	ID           int       `json:"id"`
	OrderID      int       `json:"order_id"`
	CustomerName string    `json:"customer_name"`
	OrderType    string    `json:"order_type"`
	TotalPrice   int       `json:"total_price"`
	OrderStatus  string    `json:"order_status"`
	ProofImage   string    `json:"proof_image"`
	PaidAt       time.Time `json:"paid_at"`
}
