package repository

type DeliveryRepository struct {
	DB DBExecutor
}

func NewDeliveryRepository(db DBExecutor) *DeliveryRepository {
	return &DeliveryRepository{DB: db}
}

func (r *DeliveryRepository) Create(orderID, courierID int, proofImage string) error {
	_, err := r.DB.Exec(`
		INSERT INTO deliveries (order_id, courier_id, delivery_date, status, proof_image)
		VALUES ($1, $2, CURRENT_DATE, 'delivered', $3)
	`, orderID, courierID, proofImage)

	return err
}
