package repository

import (
	"database/sql"
	"masakin-backend/app/model"
)

type MenuScheduleRepository struct {
	DB *sql.DB
}

func NewMenuScheduleRepository(db *sql.DB) *MenuScheduleRepository {
	return &MenuScheduleRepository{DB: db}
}

// assign menu ke tanggal
func (r *MenuScheduleRepository) Create(menuID int, date string) error {
	_, err := r.DB.Exec(`
		INSERT INTO menu_schedules (menu_id, serve_date)
		VALUES ($1, $2)
	`, menuID, date)

	return err
}

// get menu by tanggal (customer)
func (r *MenuScheduleRepository) GetByDate(date string) ([]model.Menu, error) {
	rows, err := r.DB.Query(`
		SELECT
			m.id, m.name, m.description, m.price, m.image_url
		FROM menu_schedules ms
		JOIN menus m ON m.id = ms.menu_id
		WHERE ms.serve_date = $1
		  AND m.is_active = true
		ORDER BY m.name
	`, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menus []model.Menu
	for rows.Next() {
		var m model.Menu
		if err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Description,
			&m.Price,
			&m.ImageURL,
		); err != nil {
			return nil, err
		}
		menus = append(menus, m)
	}

	return menus, nil
}

// ================================
// ADMIN - get schedules by month
// ================================
func (r *MenuScheduleRepository) GetByMonth(
	year int,
	month int,
) ([]model.MenuScheduleView, error) {

	rows, err := r.DB.Query(`
		SELECT
			ms.serve_date,
			m.id,
			m.name,
			m.description,
			m.price,
			m.image_url
		FROM menu_schedules ms
		JOIN menus m ON m.id = ms.menu_id
		WHERE EXTRACT(YEAR FROM ms.serve_date) = $1
		  AND EXTRACT(MONTH FROM ms.serve_date) = $2
		ORDER BY ms.serve_date, m.name
	`, year, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.MenuScheduleView
	for rows.Next() {
		var v model.MenuScheduleView
		if err := rows.Scan(
			&v.ServeDate,
			&v.MenuID,
			&v.Name,
			&v.Description,
			&v.Price,
			&v.ImageURL,
		); err != nil {
			return nil, err
		}
		result = append(result, v)
	}

	return result, nil
}

// ================================
// ADMIN - delete one menu schedule
// ================================
func (r *MenuScheduleRepository) DeleteOne(
	menuID int,
	date string,
) error {
	_, err := r.DB.Exec(`
		DELETE FROM menu_schedules
		WHERE menu_id = $1
		  AND serve_date = $2
	`, menuID, date)

	return err
}

// ================================
// ADMIN - clear all menus in date
// ================================
func (r *MenuScheduleRepository) ClearByDate(date string) error {
	_, err := r.DB.Exec(`
		DELETE FROM menu_schedules
		WHERE serve_date = $1
	`, date)

	return err
}

