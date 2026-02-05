package repository

import "database/sql"
import "masakin-backend/app/model"

type OrderItemRepository struct {
	DB DBExecutor
}
func NewOrderItemRepository(db DBExecutor) *OrderItemRepository {
	return &OrderItemRepository{DB: db}
}
// CREATE
func (r *OrderItemRepository) Create(item *model.OrderItem) error {
	_, err := r.DB.Exec(`
		INSERT INTO order_items
		(order_id, menu_id, order_date, quantity)
		VALUES ($1, $2, $3, $4)
	`,
		item.OrderID,
		item.MenuID,
		item.OrderDate,
		item.Quantity,
	)
	return err
}

func (r *OrderItemRepository) GetByOrder(orderID int) (*sql.Rows, error) {
	return r.DB.Query(`
		SELECT
			m.name, oi.quantity, m.price, oi.order_date
		FROM order_items oi
		JOIN menus m ON m.id = oi.menu_id
		WHERE oi.order_id = $1
		ORDER BY oi.order_date
	`, orderID)
}

func (r *OrderItemRepository) GetDetail(orderID int) ([]map[string]interface{}, error) {
	rows, err := r.DB.Query(`
		SELECT
			m.name,
			oi.quantity,
			m.price,
			oi.order_date
		FROM order_items oi
		JOIN menus m ON m.id = oi.menu_id
		WHERE oi.order_id = $1
		ORDER BY oi.order_date
	`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []map[string]interface{}

	for rows.Next() {
		var name string
		var qty, price int
		var date string

		rows.Scan(&name, &qty, &price, &date)

		res = append(res, map[string]interface{}{
			"menu":  name,
			"qty":   qty,
			"price": price,
			"date":  date,
		})
	}

	return res, nil
}

