package tests

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/pavel-pj/wallet/internal/model"
	"github.com/pavel-pj/wallet/internal/repository"
	repositoryMocks "github.com/pavel-pj/wallet/internal/repository/mocks"

	//"github.com/pavel-pj/wallet/internal/service"
	walletService "github.com/pavel-pj/wallet/internal/service/wallet"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {

	type walletRepositoryMockFunc func(mc *minimock.Controller) repository.WalletRepository

	type args struct {
		ctx context.Context
		req string
	}

	var (
		ctx     = context.Background()
		mc      = minimock.NewController(t)
		uuid    = "686879b8-0edd-44a0-bf18-3f3327736945"
		uuidErr = "686879b8-0edd-44a0-bf18"

		respSuccess = model.Wallet{
			Uuid:      uuid,
			Amount:    48489,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	)

	tests := []struct {
		name                 string
		args                 args
		err                  error
		want                 model.Wallet
		walletRepositoryMock walletRepositoryMockFunc
	}{
		{
			name: "success",
			args: args{
				ctx: ctx,
				req: uuid,
			},
			err:  nil,
			want: respSuccess,
			walletRepositoryMock: func(mc *minimock.Controller) repository.WalletRepository {
				mock := repositoryMocks.NewWalletRepositoryMock(mc)
				mock.GetMock.Expect(ctx, uuid).Return(respSuccess, nil)
				return mock
			},
		},

		{
			name: "error notFound",
			args: args{
				ctx: ctx,
				req: uuid,
			},
			err:  model.ErrWalletNotFound,
			want: model.Wallet{},
			walletRepositoryMock: func(mc *minimock.Controller) repository.WalletRepository {
				mock := repositoryMocks.NewWalletRepositoryMock(mc)
				mock.GetMock.Expect(ctx, uuid).Return(model.Wallet{}, model.ErrWalletNotFound)
				return mock
			},
		},

		{
			name: "errorUUID",
			args: args{
				ctx: ctx,
				req: uuidErr,
			},
			err:  model.ErrInvalidUUID,
			want: model.Wallet{},
			walletRepositoryMock: func(mc *minimock.Controller) repository.WalletRepository {
				mock := repositoryMocks.NewWalletRepositoryMock(mc)
				mock.GetMock.Expect(ctx, uuidErr).Return(model.Wallet{}, model.ErrInvalidUUID)
				return mock
			},
		},
	}

	for _, tt := range tests {

		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			service := walletService.NewWalletrService(tt.walletRepositoryMock(mc))

			res, err := service.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)

		})

	}

}
