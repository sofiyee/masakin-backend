package repository

import (
	"database/sql"
	"masakin-backend/app/model"
)

type ReportRepository struct {
	DB *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{DB: db}
}

// ==========================
// TOTAL PORSI
// ==========================
func (r *ReportRepository) GetTotalPortion(month string) (int, error) {
	var total int
	err := r.DB.QueryRow(`
		SELECT COALESCE(SUM(oi.quantity), 0)
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		WHERE to_char(o.created_at, 'YYYY-MM') = $1
		AND o.status IN ('paid','delivered')
	`, month).Scan(&total)

	return total, err
}

// ==========================
// TOTAL OMZET
// ==========================

func (r *ReportRepository) GetTotalRevenue(month string) (int, error) {
	var total int
	err := r.DB.QueryRow(`
		SELECT COALESCE(SUM(o.total_price), 0)
		FROM orders o
		WHERE to_char(o.created_at, 'YYYY-MM') = $1
		AND o.status IN ('paid','delivered')
	`, month).Scan(&total)

	return total, err
}

// ==========================
// MENU TERLARIS
// ==========================
func (r *ReportRepository) GetTopMenus(month string) ([]model.TopMenuReport, error) {
	rows, err := r.DB.Query(`
		SELECT
			m.id,
			m.name,
			SUM(oi.quantity) AS total_qty
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		JOIN menus m ON m.id = oi.menu_id
		WHERE to_char(o.created_at, 'YYYY-MM') = $1
		AND o.status IN ('paid','delivered')
		GROUP BY m.id, m.name
		ORDER BY total_qty DESC
		LIMIT 5
	`, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.TopMenuReport
	for rows.Next() {
		var t model.TopMenuReport
		if err := rows.Scan(&t.MenuID, &t.MenuName, &t.TotalQty); err != nil {
			return nil, err
		}
		res = append(res, t)
	}

	return res, nil
}




