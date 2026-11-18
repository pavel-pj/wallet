package converter

import (
	"github.com/pavel-pj/wallet/internal/model"
	repoModel "github.com/pavel-pj/wallet/internal/repository/wallet/model"
)

func FromRepoWalletToService(wallet repoModel.Wallet) model.Wallet {

	return model.Wallet{
		Uuid:      wallet.Uuid,
		Amount:    wallet.Amount,
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
	}

}
