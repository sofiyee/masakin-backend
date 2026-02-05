package repository

import "database/sql"

type WarehouseDashboardRepository struct {
	DB *sql.DB
}

func NewWarehouseDashboardRepository(db *sql.DB) *WarehouseDashboardRepository {
	return &WarehouseDashboardRepository{DB: db}
}

func (r *WarehouseDashboardRepository) GetTodayTotalPortion(date string) (int, error) {
	var total int
	err := r.DB.QueryRow(`
		SELECT COALESCE(SUM(oi.quantity),0)
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		WHERE oi.order_date = $1
		AND o.status IN ('paid','delivered')
	`, date).Scan(&total)
	return total, err
}

func (r *WarehouseDashboardRepository) GetTodayTotalMenu(date string) (int, error) {
	var total int
	err := r.DB.QueryRow(`
		SELECT COUNT(DISTINCT oi.menu_id)
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		WHERE oi.order_date = $1
		AND o.status IN ('paid','delivered')
	`, date).Scan(&total)
	return total, err
}
