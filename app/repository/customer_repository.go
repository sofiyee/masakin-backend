package repository

import (
	"database/sql"
	"masakin-backend/app/model"
)

type CustomerRepository struct {
	DB *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{DB: db}
}

func (r *CustomerRepository) Create(customer *model.Customer) error {
	_, err := r.DB.Exec(`
		INSERT INTO customers (user_id, name, region, full_address)
		VALUES ($1, $2, $3, $4)
	`,
		customer.UserID,
		customer.Name,
		customer.Region,
		customer.FullAddress,
	)

	return err
}

func (r *UserRepository) FindByID(id int) (*model.User, error) {
	row := r.DB.QueryRow(`
		SELECT id, phone, password, role
		FROM users
		WHERE id = $1
	`, id)

	var u model.User
	err := row.Scan(&u.ID, &u.Phone, &u.Password, &u.Role)
	if err != nil {
		return nil, err
	}

	return &u, nil
}


func (r *CustomerRepository) FindByName(name string) (*model.Customer, error) {
	row := r.DB.QueryRow(`
		SELECT id, user_id, name, address
		FROM customers
		WHERE name = $1
	`, name)

	var c model.Customer
	err := row.Scan(
		&c.ID,
		&c.UserID,
		&c.Name,
		&c.Region,
		&c.FullAddress)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *CustomerRepository) FindByUserID(userID int) (*model.Customer, error) {
	var c model.Customer
	err := r.DB.QueryRow(`
		SELECT id, user_id, name, region, full_address
		FROM customers
		WHERE user_id = $1
	`, userID).Scan(
		&c.ID,
		&c.UserID,
		&c.Name,
		&c.Region,
		&c.FullAddress,
	)

	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CustomerRepository) GetProfileByUserID(userID int) (*model.CustomerProfile, error) {
	var c model.CustomerProfile

	err := r.DB.QueryRow(`
		SELECT
			c.id,
			u.name,
			u.phone,
			u.address,
			COALESCE(c.region, ''),
			COALESCE(c.full_address, '')
		FROM customers c
		JOIN users u ON u.id = c.user_id
		WHERE u.id = $1
		AND u.role = 'customer'
	`, userID).Scan(
		&c.ID,
		&c.Name,
		&c.Phone,
		&c.Address,
		&c.Region,
		&c.FullAddress,
	)

	if err != nil {
		return nil, err
	}

	return &c, nil
}