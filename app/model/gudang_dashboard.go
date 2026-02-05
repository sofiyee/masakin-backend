package model

type WarehouseDashboard struct {
	Date         string `json:"date"`
	TotalPortion int    `json:"total_portion"`
	TotalMenu    int    `json:"total_menu"`
}
