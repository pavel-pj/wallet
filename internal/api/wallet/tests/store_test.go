package tests

import (
	"context"
	"errors"

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

func TestStore(t *testing.T) {

	type walletServiceMockFunc func(mc *minimock.Controller) service.WalletService

	type args struct {
		ctx context.Context
		req desc.StoreRequest
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

		respSuccess = &desc.StoreResponse{
			Wallet: &desc.Wallet{
				Uuid:      uuid,
				Amount:    amount,
				CreatedAt: timestamppb.New(createdAt),
				UpdatedAt: timestamppb.New(updatedAt),
			},
		}

		errService = errors.New("Some error from service")
	)

	tests := []struct {
		name              string
		args              args
		want              *desc.StoreResponse
		err               error
		walletServiceMock walletServiceMockFunc
	}{
		{
			name: "success",
			args: args{
				ctx: ctx,
				req: desc.StoreRequest{
					Amount: amount,
				},
			},
			want: respSuccess,
			err:  nil,
			walletServiceMock: func(mc *minimock.Controller) service.WalletService {
				mock := serviceMocks.NewWalletServiceMock(mc)
				mock.StoreMock.Expect(ctx, amount).Return(wallet, nil)
				return mock
			},
		},

		{
			name: "error",
			args: args{
				ctx: ctx,
				req: desc.StoreRequest{
					Amount: amount,
				},
			},
			want: &desc.StoreResponse{},
			err:  errService,
			walletServiceMock: func(mc *minimock.Controller) service.WalletService {
				mock := serviceMocks.NewWalletServiceMock(mc)
				mock.StoreMock.Expect(ctx, amount).Return(model.Wallet{}, errService)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			t.Parallel()

			api := walletApi.NewWalletImplimentation(tt.walletServiceMock(mc))
			resp, err := api.Store(tt.args.ctx, &tt.args.req)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, resp)

		})

	}

}
