package wallet

import (
	"github.com/pavel-pj/wallet/internal/repository"
	"github.com/pavel-pj/wallet/internal/service"
)

type WalletService struct {
	Repository repository.WalletRepository
}

func NewWalletrService(
	repository repository.WalletRepository,

) service.WalletService {
	return WalletService{
		Repository: repository,
	}
}
