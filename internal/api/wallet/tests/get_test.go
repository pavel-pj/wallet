package tests

import (
	"context"

	"github.com/gojuno/minimock/v3"
	walletApi "github.com/pavel-pj/wallet/internal/api/wallet"
	"github.com/pavel-pj/wallet/internal/model"
	"github.com/pavel-pj/wallet/internal/service"
	serviceMocks "github.com/pavel-pj/wallet/internal/service/mocks"
	"google.golang.org/protobuf/types/known/timestamppb"

	"testing"
	"time"

	desc "github.com/pavel-pj/wallet/pkg/wallet_v1"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {

	type walletServiceMockFunc func(mc *minimock.Controller) service.WalletService

	type args struct {
		ctx context.Context
		req *desc.GetRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		amount = int64(3000)

		uuid      = "686879b8-0edd-44a0-bf18-3f3327736945"
		createdAt = time.Now()
		updatedAt = time.Now()

		wallet = model.Wallet{
			Uuid:      uuid,
			Amount:    amount,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		respSuccess = &desc.GetResponse{
			Wallet: &desc.Wallet{
				Uuid:      uuid,
				Amount:    amount,
				CreatedAt: timestamppb.New(createdAt),
				UpdatedAt: timestamppb.New(updatedAt),
			},
		}
	)

	tests := []struct {
		name              string
		args              args
		want              *desc.GetResponse
		err               error
		walletServiceMock walletServiceMockFunc
	}{
		{
			name: "success",
			args: args{
				ctx: ctx,
				req: &desc.GetRequest{
					Uuid: uuid,
				},
			},
			want: respSuccess,
			err:  nil,
			walletServiceMock: func(mc *minimock.Controller) service.WalletService {
				mock := serviceMocks.NewWalletServiceMock(mc)
				mock.GetMock.Expect(ctx, uuid).Return(wallet, nil)
				return mock
			},
		},
		{
			name: "err NoFound",
			args: args{
				ctx: ctx,
				req: &desc.GetRequest{
					Uuid: uuid,
				},
			},
			want: &desc.GetResponse{},
			err:  model.ErrWalletNotFound,
			walletServiceMock: func(mc *minimock.Controller) service.WalletService {
				mock := serviceMocks.NewWalletServiceMock(mc)
				mock.GetMock.Expect(ctx, uuid).Return(model.Wallet{}, model.ErrWalletNotFound)
				return mock
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			t.Parallel()

			api := walletApi.NewWalletImplimentation(tt.walletServiceMock(mc))
			resp, err := api.Get(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, resp)

		})

	}

}
