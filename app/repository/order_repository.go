package repository

import (
	
	"masakin-backend/app/model"
)

type OrderRepository struct {
	DB DBExecutor
}

func NewOrderRepository(db DBExecutor) *OrderRepository {
	return &OrderRepository{DB: db}
}


// ====================
// CREATE
// ====================
func (r *OrderRepository) Create(o *model.Order) (int, error) {
	var id int
	err := r.DB.QueryRow(`
		INSERT INTO orders
		(customer_id, order_type, start_date, end_date, status, total_price)
		VALUES ($1, $2, $3, $4, 'unpaid', 0)
		RETURNING id
	`,
		o.CustomerID,
		o.OrderType,
		o.StartDate,
		o.EndDate,
	).Scan(&id)

	return id, err
}

// ====================
// UPDATE TOTAL
// ====================
func (r *OrderRepository) UpdateTotal(orderID int, total int) error {
	_, err := r.DB.Exec(`
		UPDATE orders
		SET total_price = $1
		WHERE id = $2
	`, total, orderID)
	return err
}

func (r *OrderRepository) GetAll() ([]model.OrderAdminList, error) {
	rows, err := r.DB.Query(`
		SELECT
			o.id,
			u.name AS customer_name,
			o.order_type,
			o.status,
			o.total_price,
			o.created_at
		FROM orders o
		JOIN customers c ON c.id = o.customer_id
		JOIN users u ON u.id = c.user_id
		ORDER BY o.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.OrderAdminList

	for rows.Next() {
		var o model.OrderAdminList
		if err := rows.Scan(
			&o.ID,
			&o.CustomerName,
			&o.OrderType,
			&o.Status,
			&o.TotalPrice,
			&o.CreatedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, o)
	}

	return res, nil
}


// ====================
// GET BY ID (DETAIL)
// ====================
func (r *OrderRepository) GetByID(orderID int) (*model.Order, error) {
	var o model.Order

	err := r.DB.QueryRow(`
		SELECT id, customer_id, order_type, start_date, end_date, status, total_price, created_at
		FROM orders
		WHERE id = $1
	`, orderID).Scan(
		&o.ID,
		&o.CustomerID,
		&o.OrderType,
		&o.StartDate,
		&o.EndDate,
		&o.Status,
		&o.TotalPrice,
		&o.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &o, nil
}

// ====================
// UPDATE STATUS (ADMIN)
// ====================
func (r *OrderRepository) UpdateStatus(orderID int, status string) error {
	_, err := r.DB.Exec(`
		UPDATE orders
		SET status = $1
		WHERE id = $2
	`, status, orderID)

	return err
}

// ================================
// GET ORDERS FOR COURIER (BY REGION)
// ================================
func (r *OrderRepository) GetForCourier(kurirUserID int) ([]model.Order, error) {
	rows, err := r.DB.Query(`
		SELECT
			o.id,
			o.customer_id,
			o.order_type,
			o.start_date,
			o.end_date,
			o.status,
			o.total_price,
			o.created_at
		FROM orders o
		JOIN customers c ON c.id = o.customer_id
		JOIN couriers cr ON cr.region = c.region
		WHERE cr.user_id = $1
		AND cr.active = true
		AND o.status = 'paid'
		ORDER BY o.created_at ASC
	`, kurirUserID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		err := rows.Scan(
			&o.ID,
			&o.CustomerID,
			&o.OrderType,
			&o.StartDate,
			&o.EndDate,
			&o.Status,
			&o.TotalPrice,
			&o.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}

// ================================
// DELIVER ORDER
// ================================
func (r *OrderRepository) Deliver(orderID int) error {
	_, err := r.DB.Exec(`
		UPDATE orders
		SET status = 'delivered'
		WHERE id = $1
		AND status = 'paid'
	`, orderID)

	return err
}

// ================================
// GET RECENT ORDERS FOR DASHBOARD
// ================================
func (r *OrderRepository) GetRecentOrders(limit int) ([]model.OrderAdminList, error) {
	rows, err := r.DB.Query(`
		SELECT
			o.id,
			u.name AS customer_name,
			o.order_type,
			o.status,
			o.total_price,
			o.created_at
		FROM orders o
		JOIN customers c ON c.id = o.customer_id
		JOIN users u ON u.id = c.user_id
		ORDER BY o.created_at DESC
		LIMIT $1
	`, limit)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.OrderAdminList

	for rows.Next() {
		var order model.OrderAdminList
		if err := rows.Scan(
			&order.ID,
			&order.CustomerName,
			&order.OrderType,
			&order.Status,
			&order.TotalPrice,
			&order.CreatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) GetUnpaidOrdersAdmin() ([]model.UnpaidOrderAdminView, error) {
	rows, err := r.DB.Query(`
		SELECT
			o.id,
			c.name AS customer_name,
			o.order_type,
			o.total_price
		FROM orders o
		LEFT JOIN customers c ON c.id = o.customer_id
		WHERE o.status = 'unpaid'
		ORDER BY o.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]model.UnpaidOrderAdminView, 0)

	for rows.Next() {
		var o model.UnpaidOrderAdminView
		if err := rows.Scan(
			&o.ID,
			&o.CustomerName,
			&o.OrderType,
			&o.TotalPrice,
		); err != nil {
			return nil, err
		}

		orders = append(orders, o)
	}

	return orders, nil
}

func (r *OrderRepository) GetRecentOrdersByCustomer(
	customerID int,
	limit int,
) ([]model.RecentOrderCustomerView, error) {

	rows, err := r.DB.Query(`
		SELECT
			o.id,
			o.order_type,
			o.status,
			o.total_price,
			TO_CHAR(o.created_at, 'YYYY-MM-DD')
		FROM orders o
		WHERE o.customer_id = $1
		ORDER BY o.created_at DESC
		LIMIT $2
	`, customerID, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]model.RecentOrderCustomerView, 0)

	for rows.Next() {
		var o model.RecentOrderCustomerView
		if err := rows.Scan(
			&o.OrderID,
			&o.OrderType,
			&o.Status,
			&o.TotalPrice,
			&o.CreatedAt,
		); err != nil {
			return nil, err
		}

		orders = append(orders, o)
	}

	return orders, nil
}