package model

import "time"

type Wallet struct {
	Uuid      string
	Amount    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
