package repository

import (
	"database/sql"
	"masakin-backend/app/model"
)

type CourierRepository struct {
	DB *sql.DB
}

func NewCourierRepository(db *sql.DB) *CourierRepository {
	return &CourierRepository{DB: db}
}

// CREATE
func (r *CourierRepository) Create(userID int, region string) error {
	_, err := r.DB.Exec(`
		INSERT INTO couriers (user_id, region)
		VALUES ($1, $2)
	`, userID, region)

	return err
}

// LIST
func (r *CourierRepository) GetAll() ([]model.Courier, error) {
	rows, err := r.DB.Query(`
		SELECT
			c.id,
			u.id,
			u.name,
			u.phone,
			c.region,
			c.active
		FROM couriers c
		JOIN users u ON u.id = c.user_id
		ORDER BY c.id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Courier
	for rows.Next() {
		var c model.Courier
		if err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.Name,
			&c.Phone,
			&c.Region,
			&c.Active,
		); err != nil {
			return nil, err
		}
		result = append(result, c)
	}

	return result, nil
}

// UPDATE REGION
func (r *CourierRepository) UpdateRegion(userID int, region string) error {
	_, err := r.DB.Exec(`
		UPDATE couriers
		SET region = $1
		WHERE user_id = $2
	`, region, userID)

	return err
}

// ENABLE / DISABLE
func (r *CourierRepository) SetActive(userID int, active bool) error {
	_, err := r.DB.Exec(`
		UPDATE couriers
		SET active = $1
		WHERE user_id = $2
	`, active, userID)

	return err
}
