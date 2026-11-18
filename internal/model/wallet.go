package model

import (
	"time"
)

type Wallet struct {
	Uuid      string
	Amount    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateWallet struct {
	Uuid   string
	Type   OperationType
	Amount int64
}

type OperationType = string

const (
	OperationDeposit  OperationType = "DEPOSIT"
	OperationWithdraw OperationType = "WITHDRAW"
)
