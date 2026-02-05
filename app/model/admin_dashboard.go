package model

type AdminDashboard struct {
	TodayOrders    int `json:"today_orders"`
	MonthlyOrders  int `json:"monthly_orders"`
	MonthlyRevenue int `json:"monthly_revenue"`
	ActiveCustomer int `json:"active_customer"`
	ActiveCourier  int `json:"active_courier"`
}

