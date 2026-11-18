package user

import (
	"context"

	converter "github.com/pavel-pj/wallet/internal/converter/wallet"

	desc "github.com/pavel-pj/wallet/pkg/wallet_v1"
)

func (s *WalletImplementation) Store(ctx context.Context, req *desc.StoreRequest) (*desc.StoreResponse, error) {

	amount := req.GetAmount()
	wallet, err := s.serv.Store(ctx, amount)

	if err != nil {

		return &desc.StoreResponse{}, err
	}

	response := converter.ToProtoWalletFromService(wallet)

	return &desc.StoreResponse{Wallet: response}, nil
}
