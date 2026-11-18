package wallet

import (
	db "github.com/pavel-pj/wallet/internal/client/db"
	"github.com/pavel-pj/wallet/internal/repository"

	sq "github.com/Masterminds/squirrel"
)

type WalletRepository struct {
	db db.Client
}

const (
	tableName       = "wallets"
	columnUUID      = "uuid"
	columnAmount    = "amount"
	columnCreatedAt = "created_at"
	columnUpdatedAt = "updated_at"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func NewWalletRepository(db db.Client) repository.WalletRepository {
	return &WalletRepository{
		db: db,
	}
}
