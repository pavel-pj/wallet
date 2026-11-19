package user

import (
	"context"

	converter "github.com/pavel-pj/wallet/internal/converter/wallet"

	desc "github.com/pavel-pj/wallet/pkg/wallet_v1"
)

func (s *WalletImplementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {

	uuid := req.GetUuid()
	wallet, err := s.serv.Get(ctx, uuid)

	if err != nil {

		return &desc.GetResponse{}, err
	}

	response := converter.ToProtoWalletFromService(wallet)

	return &desc.GetResponse{Wallet: response}, nil
}
