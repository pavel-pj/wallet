package wallet

import (
	"context"

	"github.com/pavel-pj/wallet/internal/model"
)

func (s WalletService) Wallet(ctx context.Context, data model.UpdateWallet) (model.Wallet, error) {

	switch data.Type {
	case model.OperationDeposit, model.OperationWithdraw:
	default:
		return model.Wallet{}, model.ErrInvalidOperationType
	}
	if data.Amount <= 0 {
		return model.Wallet{}, model.ErrInvalidAmount
	}

	wallet, err := s.Repository.Wallet(ctx, data)
	if err != nil {
		return model.Wallet{}, err
	}

	return wallet, nil

}
