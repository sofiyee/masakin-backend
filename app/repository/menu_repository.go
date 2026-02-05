package repository

import (
	"masakin-backend/app/model"
)

type MenuRepository struct {
	DB DBExecutor
}

func NewMenuRepository(db DBExecutor) *MenuRepository {
	return &MenuRepository{DB: db}
}

// CREATE
func (r *MenuRepository) Create(m *model.Menu) error {
	_, err := r.DB.Exec(`
		INSERT INTO menus
		(name, description, price, image_url, menu_month, menu_year, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, true)
	`,
		m.Name,
		m.Description,
		m.Price,
		m.ImageURL,
		m.MenuMonth,
		m.MenuYear,
	)
	return err
}

// GET BY MONTH & YEAR (CUSTOMER + ADMIN)
func (r *MenuRepository) GetByMonthYear(month, year int) ([]model.Menu, error) {
	rows, err := r.DB.Query(`
		SELECT
			id, name, description, price,
			image_url, menu_month, menu_year,
			is_active, created_at
		FROM menus
		WHERE menu_month = $1
		  AND menu_year = $2
		  AND is_active = true
		ORDER BY id DESC
	`, month, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.Menu
	for rows.Next() {
		var m model.Menu
		if err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Description,
			&m.Price,
			&m.ImageURL,
			&m.MenuMonth,
			&m.MenuYear,
			&m.IsActive,
			&m.CreatedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

// UPDATE
func (r *MenuRepository) Update(id int, m *model.Menu) error {
	_, err := r.DB.Exec(`
		UPDATE menus
		SET name=$1,
		    description=$2,
		    price=$3,
		    image_url=$4,
		    menu_month=$5,
		    menu_year=$6
		WHERE id=$7
	`,
		m.Name,
		m.Description,
		m.Price,
		m.ImageURL,
		m.MenuMonth,
		m.MenuYear,
		id,
	)
	return err
}

// SOFT DELETE
func (r *MenuRepository) Disable(id int) error {
	_, err := r.DB.Exec(`
		UPDATE menus
		SET is_active = false
		WHERE id = $1
	`, id) // 🔥 INI YANG KURANG
	return err
}
// ENABLE

func (r *MenuRepository) Enable(id int) error {
	_, err := r.DB.Exec(`
		UPDATE menus
		SET is_active = true
		WHERE id = $1
	`, id)
	return err
}

func (r *MenuRepository) GetPriceByID(menuID int) (int, error) {
	var price int
	err := r.DB.QueryRow(`
		SELECT price
		FROM menus
		WHERE id = $1
		  AND is_active = true
	`, menuID).Scan(&price)

	return price, err
}

// GET BY MONTH & YEAR FOR ADMIN (ALL STATUS)
func (r *MenuRepository) GetByMonthYearAdmin(month, year int) ([]model.Menu, error) {
	rows, err := r.DB.Query(`
		SELECT
			id, name, description, price,
			image_url, menu_month, menu_year,
			is_active, created_at
		FROM menus
		WHERE menu_month = $1
		  AND menu_year = $2
		ORDER BY is_active DESC, id DESC
	`, month, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.Menu
	for rows.Next() {
		var m model.Menu
		if err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Description,
			&m.Price,
			&m.ImageURL,
			&m.MenuMonth,
			&m.MenuYear,
			&m.IsActive,
			&m.CreatedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}