package wallet

import (
	"context"

	"github.com/pavel-pj/wallet/internal/model"
)

func (s WalletService) Store(ctx context.Context, amount int64) (model.Wallet, error) {

	wallet, err := s.Repository.Store(ctx, amount)
	if err != nil {
		return model.Wallet{}, err
	}

	return wallet, nil

}
