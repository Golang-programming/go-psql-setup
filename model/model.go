package model

type Stock struct {
	ID      int64   `json:"stock"`
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	Company string  `json:"company"`
}
