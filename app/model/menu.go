package model

import "time"

type Menu struct {
	ID          int
	Name        string
	Description string
	Price       int
	ImageURL    string
	MenuMonth   int
	MenuYear    int
	IsActive    bool
	CreatedAt   time.Time
}
