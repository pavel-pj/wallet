package user

import (
	"context"

	converter "github.com/pavel-pj/wallet/internal/converter/wallet"
	"github.com/pavel-pj/wallet/internal/model"

	desc "github.com/pavel-pj/wallet/pkg/wallet_v1"
)

func (s *WalletImplementation) Wallet(ctx context.Context, req *desc.WalletRequest) (*desc.WalletResponse, error) {

	data := model.UpdateWallet{
		Uuid:   req.GetWalletId(),
		Type:   req.GetOperationType().String(),
		Amount: req.GetAmount(),
	}

	wallet, err := s.serv.Wallet(ctx, data)

	if err != nil {

		return &desc.WalletResponse{}, err
	}

	response := converter.ToProtoWalletFromService(wallet)

	return &desc.WalletResponse{Wallet: response}, nil
}
