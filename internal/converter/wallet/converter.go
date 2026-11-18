package user

import (
	"time"

	"github.com/pavel-pj/wallet/internal/model"
	desc "github.com/pavel-pj/wallet/pkg/wallet_v1"

	"google.golang.org/protobuf/types/known/timestamppb"
	//"google.golang.org/protobuf/types/known/wrapperspb"
)

func ToProtoWalletFromService(data model.Wallet) *desc.Wallet {

	return &desc.Wallet{
		Uuid:      data.Uuid,
		Amount:    data.Amount,
		CreatedAt: toProtoTimestamp(data.CreatedAt),
		UpdatedAt: toProtoTimestamp(data.UpdatedAt),
	}

}

func toProtoTimestamp(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}
