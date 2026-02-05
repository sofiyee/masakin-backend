package model

type KitchenSummary struct {
	MenuID   int    `json:"menu_id"`
	MenuName string `json:"menu_name"`
	TotalQty int    `json:"total_qty"`
}
