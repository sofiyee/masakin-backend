package repository

import "database/sql"
import "time"

type CustomerDashboardRepository struct {
	DB *sql.DB
}

func NewCustomerDashboardRepository(db *sql.DB) *CustomerDashboardRepository {
	return &CustomerDashboardRepository{DB: db}
}

func (r *CustomerDashboardRepository) GetActiveOrder(customerID int) (bool, string, error) {
	var nextDate *time.Time

	err := r.DB.QueryRow(`
		SELECT MIN(start_date)
		FROM orders
		WHERE customer_id = $1
		AND status = 'paid'
	`, customerID).Scan(&nextDate)

	if err != nil {
		return false, "", err
	}

	if nextDate == nil {
		return false, "", nil
	}

	return true, nextDate.Format("2006-01-02"), nil
}





func (r *CustomerDashboardRepository) GetMonthlyOrders(customerID int, month string) (int, error) {
	var total int
	err := r.DB.QueryRow(`
		SELECT COUNT(*)
		FROM orders
		WHERE customer_id = $1
		AND to_char(created_at, 'YYYY-MM') = $2
	`, customerID, month).Scan(&total)

	return total, err
}

func (r *CustomerDashboardRepository) GetMonthlySpending(customerID int, month string) (int, error) {
	var total int
	err := r.DB.QueryRow(`
		SELECT COALESCE(SUM(total_price),0)
		FROM orders
		WHERE customer_id = $1
		AND to_char(created_at, 'YYYY-MM') = $2
		AND status IN ('paid','delivered')
	`, customerID, month).Scan(&total)

	return total, err
}

