package user

import (
	"github.com/pavel-pj/wallet/internal/service"
	desc "github.com/pavel-pj/wallet/pkg/wallet_v1"
)

type WalletImplementation struct {
	desc.UnimplementedWalletV1Server
	serv service.WalletService
}

func NewWalletImplimentation(walletService service.WalletService) *WalletImplementation {
	return &WalletImplementation{
		serv: walletService,
	}
}
