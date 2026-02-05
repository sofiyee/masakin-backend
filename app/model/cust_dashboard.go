package model

type CustomerDashboard struct {
	HasActiveOrder   bool   `json:"has_active_order"`
	NextDeliveryDate string `json:"next_delivery_date,omitempty"`
	MonthlyOrders    int    `json:"monthly_orders"`
	MonthlySpending  int    `json:"monthly_spending"`
}

