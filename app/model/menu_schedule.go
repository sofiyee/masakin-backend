package model

import "time"

type MenuSchedule struct {
	ID        int
	MenuID    int
	ServeDate time.Time
}

type MenuScheduleView struct {
  ServeDate   time.Time `json:"serve_date"`
  MenuID      int       `json:"menu_id"`
  Name        string    `json:"name"`
  Description string    `json:"description"`
  Price       int       `json:"price"`
  ImageURL    string    `json:"image_url"`
}
