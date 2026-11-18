package wallet

import (
	"context"

	db "github.com/pavel-pj/wallet/internal/client/db"
	"github.com/pavel-pj/wallet/internal/model"
	"github.com/pavel-pj/wallet/internal/repository/wallet/converter"
	repoWallet "github.com/pavel-pj/wallet/internal/repository/wallet/model"
)

func (r WalletRepository) Store(ctx context.Context, amount int64) (model.Wallet, error) {

	sqlStr, args, err := psql.Insert(tableName).
		Columns(columnAmount).
		Values(amount).
		Suffix("RETURNING uuid::text, amount,created_at,updated_at").
		ToSql()

	if err != nil {
		return model.Wallet{}, err
	}

	q := db.Query{
		Name:     "wallet_repository.Store",
		QueryRow: sqlStr,
	}

	var repoWallet repoWallet.Wallet
	err = r.db.DB().ScanOneContext(ctx, &repoWallet, q, args...)
	if err != nil {
		return model.Wallet{}, err
	}

	wallet := converter.FromRepoWalletToService(repoWallet)

	return wallet, nil

}
