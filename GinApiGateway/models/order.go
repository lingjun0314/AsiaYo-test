package models

type Order struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Address  Address `json:"address"`
	Price    string  `json:"price"`
	Currency string  `json:"currency"`
}
