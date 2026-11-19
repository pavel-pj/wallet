package tests

import (
	"context"
	"errors"
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

func TestWallet(t *testing.T) {

	type walletRepositoryMockFunc func(mc *minimock.Controller) repository.WalletRepository

	type args struct {
		ctx context.Context
		req model.UpdateWallet
	}

	var (
		ctx    = context.Background()
		mc     = minimock.NewController(t)
		amount = int64(2000)
		uuid   = "686879b8-0edd-44a0-bf18-3f3327736945"

		requestSuccess = model.UpdateWallet{
			Uuid:   uuid,
			Type:   model.OperationDeposit,
			Amount: amount,
		}

		requestErrorType = model.UpdateWallet{
			Uuid:   uuid,
			Type:   "FDFFD",
			Amount: amount,
		}

		requestErrorAmount = model.UpdateWallet{
			Uuid:   uuid,
			Type:   model.OperationDeposit,
			Amount: -100,
		}

		respSuccess = model.Wallet{
			Uuid:      uuid,
			Amount:    amount,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		errorRepo = errors.New("Some error from repository")
	)

	tests := []struct {
		name                 string
		args                 args
		err                  error
		want                 model.Wallet
		walletRepositoryMock walletRepositoryMockFunc
	}{
		{
			name: "success Deposit",
			args: args{
				ctx: ctx,
				req: requestSuccess,
			},
			err:  nil,
			want: respSuccess,
			walletRepositoryMock: func(mc *minimock.Controller) repository.WalletRepository {
				mock := repositoryMocks.NewWalletRepositoryMock(mc)
				mock.WalletMock.Expect(ctx, requestSuccess).Return(respSuccess, nil)
				return mock
			},
		},
		{
			name: "Error operation type",
			args: args{
				ctx: ctx,
				req: requestErrorType,
			},
			err:  model.ErrInvalidOperationType,
			want: model.Wallet{},
			walletRepositoryMock: func(mc *minimock.Controller) repository.WalletRepository {
				mock := repositoryMocks.NewWalletRepositoryMock(mc)
				return mock
			},
		},

		{
			name: "Error InvalidAmount",
			args: args{
				ctx: ctx,
				req: requestErrorAmount,
			},
			err:  model.ErrInvalidAmount,
			want: model.Wallet{},
			walletRepositoryMock: func(mc *minimock.Controller) repository.WalletRepository {
				mock := repositoryMocks.NewWalletRepositoryMock(mc)
				return mock
			},
		},

		{
			name: "Error Repository",
			args: args{
				ctx: ctx,
				req: requestSuccess,
			},
			err:  errorRepo,
			want: model.Wallet{},
			walletRepositoryMock: func(mc *minimock.Controller) repository.WalletRepository {
				mock := repositoryMocks.NewWalletRepositoryMock(mc)
				mock.WalletMock.Expect(ctx, requestSuccess).Return(model.Wallet{}, errorRepo)
				return mock
			},
		},
	}

	for _, tt := range tests {

		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			service := walletService.NewWalletrService(tt.walletRepositoryMock(mc))

			res, err := service.Wallet(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, res)

		})

	}

}
