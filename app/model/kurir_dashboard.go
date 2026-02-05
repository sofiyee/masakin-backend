package model

type CourierDashboard struct {
	Date           string `json:"date"`
	TotalOrders    int    `json:"total_orders"`
	PendingOrders  int    `json:"pending_orders"`
	DeliveredOrders int   `json:"delivered_orders"`
}
