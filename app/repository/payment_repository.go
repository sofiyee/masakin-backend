package repository

import "masakin-backend/app/model"

type PaymentRepository struct {
	DB DBExecutor // interface yg sama kayak repo lain (db / tx)
}

func NewPaymentRepository(db DBExecutor) *PaymentRepository {
	return &PaymentRepository{DB: db}
}

// CREATE PAYMENT (ADMIN)
func (r *PaymentRepository) Create(orderID int, proofImage string, adminID int) error {
	_, err := r.DB.Exec(`
		INSERT INTO payments (order_id, proof_image, verified_by)
		VALUES ($1, $2, $3)
	`, orderID, proofImage, adminID)

	return err
}

// ambil proof image lama
func (r *PaymentRepository) GetProofImageByOrderID(orderID int) (string, error) {
	var filename string
	err := r.DB.QueryRow(`
		SELECT proof_image
		FROM payments
		WHERE order_id = $1
	`, orderID).Scan(&filename)

	return filename, err
}

// update proof image
func (r *PaymentRepository) UpdateProofImage(orderID int, filename string) error {
	_, err := r.DB.Exec(`
		UPDATE payments
		SET proof_image = $1
		WHERE order_id = $2
	`, filename, orderID)

	return err
}

func (r *PaymentRepository) GetAllAdmin() ([]model.PaymentAdminView, error) {
		rows, err := r.DB.Query(`
		SELECT
		p.id,
		p.order_id,
		c.name AS customer_name,
		o.order_type,
		o.total_price,
		o.status AS order_status,
		p.proof_image,
		p.paid_at
		FROM payments p
		JOIN orders o ON o.id = p.order_id
		LEFT JOIN customers c ON c.id = o.customer_id
		ORDER BY p.paid_at DESC
		`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 🔑 INIT SLICE (ANTI null)
	payments := make([]model.PaymentAdminView, 0)

	for rows.Next() {
		var p model.PaymentAdminView
		if err := rows.Scan(
			&p.ID,
			&p.OrderID,
			&p.CustomerName,
			&p.OrderType,
			&p.TotalPrice,
			&p.OrderStatus,
			&p.ProofImage,
			&p.PaidAt,
		); err != nil {
			return nil, err
		}

		payments = append(payments, p)
	}

	return payments, nil
}


