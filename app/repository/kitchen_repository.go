package repository
import (
	"database/sql"
	"masakin-backend/app/model"
)

type KitchenRepository struct {
	DB *sql.DB
}
func NewKitchenRepository(db *sql.DB) *KitchenRepository {
	return &KitchenRepository{DB: db}
}

// ==========================
// REKAP GUDANG HARIAN
// ==========================
func (r *KitchenRepository) GetDailySummary(date string) ([]model.KitchenSummary, error) {
	rows, err := r.DB.Query(`
		SELECT
			m.id,
			m.name,
			SUM(oi.quantity) AS total_qty
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		JOIN menus m ON m.id = oi.menu_id
		WHERE oi.order_date = $1
		AND o.status IN ('paid','delivered')
		GROUP BY m.id, m.name
		ORDER BY m.name
	`, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.KitchenSummary
	for rows.Next() {
		var k model.KitchenSummary
		if err := rows.Scan(&k.MenuID, &k.MenuName, &k.TotalQty); err != nil {
			return nil, err
		}
		res = append(res, k)
	}

	return res, nil
}