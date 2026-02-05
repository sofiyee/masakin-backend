package repository

import "database/sql"

type CourierDashboardRepository struct {
	DB *sql.DB
}

func NewCourierDashboardRepository(db *sql.DB) *CourierDashboardRepository {
	return &CourierDashboardRepository{DB: db}
}

func (r *CourierDashboardRepository) GetCourierStats(courierID int, date string) (total, pending, delivered int, err error) {

	query := `
	SELECT
		COUNT(*) FILTER (WHERE o.status IN ('paid','delivered')) AS total,
		COUNT(*) FILTER (WHERE o.status = 'paid') AS pending,
		COUNT(*) FILTER (WHERE o.status = 'delivered') AS delivered
	FROM orders o
	JOIN customers c ON c.id = o.customer_id
	JOIN couriers cr ON cr.region = c.region
	WHERE cr.user_id = $1
	AND DATE(o.created_at) = $2
	`

	err = r.DB.QueryRow(query, courierID, date).
		Scan(&total, &pending, &delivered)

	return
}
