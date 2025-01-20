package dtos

import "time"

type GetLastRequest struct {
	Count int `json:"count" validate:"gte=0"`
}

type GetLastResponse struct {
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	ID          int       `json:"id"`
	FromAddress string    `json:"from"`
	ToAddress   string    `json:"to"`
	Amount      string    `json:"amount"`
	Timestamp   time.Time `json:"timestamp"`
}

type SendRequest struct {
	FromAddress string `json:"from"`
	ToAddress   string `json:"to"`
	Amount      string `json:"amount"`
}

type SendResponse struct {
}
