package model

type MonthlyReport struct {
	Month        string           `json:"month"`
	TotalPortion int              `json:"total_portion"`
	TotalRevenue int              `json:"total_revenue"`
	TopMenus     []TopMenuReport  `json:"top_menus"`
}

type TopMenuReport struct {
	MenuID   int    `json:"menu_id"`
	MenuName string `json:"menu_name"`
	TotalQty int    `json:"total_qty"`
}
