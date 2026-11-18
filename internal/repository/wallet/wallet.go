package wallet

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	db "github.com/pavel-pj/wallet/internal/client/db"
	"github.com/pavel-pj/wallet/internal/model"
	"github.com/pavel-pj/wallet/internal/repository/wallet/converter"
	repoWallet "github.com/pavel-pj/wallet/internal/repository/wallet/model"
)

func (r WalletRepository) Wallet(ctx context.Context, data model.UpdateWallet) (model.Wallet, error) {

	builder := psql.Update(tableName).
		Set(columnUpdatedAt, sq.Expr("NOW()")).
		Where(sq.Eq{columnUUID: data.Uuid})

	switch data.Type {
	case model.OperationDeposit:
		builder = builder.Set(columnAmount, sq.Expr("amount + ?", data.Amount))

	case model.OperationWithdraw:
		builder = builder.Set(columnAmount, sq.Expr("amount - ?", data.Amount)).
			Where(sq.Expr("amount >= ?", data.Amount))

	default:
		return model.Wallet{}, errors.New("unknown operations type")
	}

	sqlStr, args, err := builder.
		Suffix("RETURNING uuid::text, amount, created_at, updated_at").
		ToSql()

	if err != nil {
		return model.Wallet{}, err
	}

	q := db.Query{
		Name:     "wallet_repository.Wallet",
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
