package repository

import "database/sql"

type DashboardRepository struct {
	DB *sql.DB
}

func NewDashboardRepository(db *sql.DB) *DashboardRepository {
	return &DashboardRepository{DB: db}
}

func (r *DashboardRepository) GetTodayOrders() (int, error) {
	var total int
	err := r.DB.QueryRow(`
		SELECT COUNT(*)
		FROM orders
		WHERE DATE(created_at) = CURRENT_DATE
	`).Scan(&total)
	return total, err
}

func (r *DashboardRepository) GetMonthlyOrders(month string) (int, error) {
	var total int
	err := r.DB.QueryRow(`
		SELECT COUNT(*)
		FROM orders
		WHERE to_char(created_at, 'YYYY-MM') = $1
	`, month).Scan(&total)
	return total, err
}

func (r *DashboardRepository) GetMonthlyRevenue(month string) (int, error) {
	var total int
	err := r.DB.QueryRow(`
		SELECT COALESCE(SUM(total_price),0)
		FROM orders
		WHERE to_char(created_at, 'YYYY-MM') = $1
		AND status IN ('paid','delivered')
	`, month).Scan(&total)
	return total, err
}

func (r *DashboardRepository) GetActiveCustomer() (int, error) {
	var total int
	err := r.DB.QueryRow(`
		SELECT COUNT(*)
		FROM users
		WHERE role = 'customer'
	`).Scan(&total)
	return total, err
}

func (r *DashboardRepository) GetActiveCourier() (int, error) {
	var total int
	err := r.DB.QueryRow(`
		SELECT COUNT(*)
		FROM couriers
		WHERE active = true
	`).Scan(&total)
	return total, err
}
