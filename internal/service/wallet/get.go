package wallet

import (
	"context"

	"github.com/pavel-pj/wallet/internal/model"
)

func (s WalletService) Get(ctx context.Context, uuid string) (model.Wallet, error) {

	wallet, err := s.Repository.Get(ctx, uuid)
	if err != nil {
		return model.Wallet{}, err
	}

	return wallet, nil

}
