package wallet

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/pavel-pj/wallet/internal/client/db"
	"github.com/pavel-pj/wallet/internal/model"
	"github.com/pavel-pj/wallet/internal/repository/wallet/converter"
	repoWallet "github.com/pavel-pj/wallet/internal/repository/wallet/model"
)

func (r WalletRepository) Get(ctx context.Context, uuid string) (model.Wallet, error) {

	sqlStr, args, err := psql.Select().
		From(tableName).
		Columns(columnUUID, columnAmount, columnCreatedAt, columnUpdatedAt).
		Where(sq.Eq{columnUUID: uuid}).
		ToSql()

	if err != nil {
		return model.Wallet{}, err
	}

	q := db.Query{
		Name:     "wallet_repository.Get",
		QueryRow: sqlStr,
	}

	var repoWallet repoWallet.Wallet
	err = r.db.DB().ScanOneContext(ctx, &repoWallet, q, args...)
	if err != nil {
		var pgErr *pgconn.PgError

		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return model.Wallet{}, model.ErrWalletNotFound

		// Неверный uuid
		case errors.As(err, &pgErr) && pgErr.Code == "22P02":
			return model.Wallet{}, model.ErrInvalidUUID

		default:
			return model.Wallet{}, fmt.Errorf("get wallet: %w", err)
		}
	}

	wallet := converter.FromRepoWalletToService(repoWallet)

	return wallet, nil

}
