package model


type Customer struct {
	ID          int
	UserID      int
	Name        string
	Region      string
	FullAddress string
}

type CustomerProfile struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	Region      string `json:"region"`
	FullAddress string `json:"full_address"`
}
