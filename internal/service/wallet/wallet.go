package wallet

import (
	"context"

	"github.com/pavel-pj/wallet/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s WalletService) Wallet(ctx context.Context, data model.UpdateWallet) (model.Wallet, error) {

	switch data.Type {
	case model.OperationDeposit, model.OperationWithdraw:
	default:
		return model.Wallet{}, status.Error(codes.InvalidArgument, "invalid operation type")
	}
	if data.Amount <= 0 {
		return model.Wallet{}, status.Error(codes.InvalidArgument, "amount must be positive")
	}

	wallet, err := s.Repository.Wallet(ctx, data)
	if err != nil {
		return model.Wallet{}, err
	}

	return wallet, nil

}
