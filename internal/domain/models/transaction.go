package models

import "time"

// Transaction represent balance transferring between wallets
// or trying to do this, which indicates by Successful field
type Transaction struct {
	ID          int
	FromAddress string // Source Wallet Address
	ToAddress   string // Target Wallet Address
	Amount      Balance
	Timestamp   time.Time
	Successful  bool
}
