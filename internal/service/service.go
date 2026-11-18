package service

import (
	"context"

	"github.com/pavel-pj/wallet/internal/model"
)

type WalletService interface {
	Store(ctx context.Context, amount int64) (model.Wallet, error)
	Wallet(ctx context.Context, wallet model.UpdateWallet) (model.Wallet, error)
}
