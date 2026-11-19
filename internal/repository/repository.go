package repository

import (
	"context"

	"github.com/pavel-pj/wallet/internal/model"
)

type WalletRepository interface {
	Store(ctx context.Context, amount int64) (model.Wallet, error)
	Wallet(ctx context.Context, data model.UpdateWallet) (model.Wallet, error)
	Get(ctx context.Context, uuid string) (model.Wallet, error)
}
