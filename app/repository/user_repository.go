package repository

import (
	"database/sql"
	"masakin-backend/app/model"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) FindByPhone(phone string) (*model.User, error) {
	row := r.DB.QueryRow(`
		SELECT id, phone, password, role
		FROM users
		WHERE phone = $1
	`, phone)

	user := model.User{}
	err := row.Scan(&user.ID, &user.Phone, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Create(user *model.User) (int, error) {
	var id int
	err := r.DB.QueryRow(`
		INSERT INTO users (name, address, password, role, phone)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`,
		user.Name,
		user.Address,
		user.Password,
		user.Role,
		user.Phone,
	).Scan(&id)

	return id, err
}


func (r *UserRepository) FindByName(name string) (*model.User, error) {
	row := r.DB.QueryRow(`
		SELECT id, name, password, role
		FROM users
		WHERE name = $1
	`, name)

	var u model.User
	err := row.Scan(&u.ID, &u.Name, &u.Password, &u.Role)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
