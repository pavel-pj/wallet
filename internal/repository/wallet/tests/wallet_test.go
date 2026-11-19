//go:build integration

package tests

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/pavel-pj/wallet/internal/model"
	walletRepo "github.com/pavel-pj/wallet/internal/repository/wallet"
	"github.com/stretchr/testify/require"
)

func TestWallet_Deposit(t *testing.T) {
	repo := walletRepo.NewWalletRepository(testDB)
	ctx := context.Background()

	created, err := repo.Store(ctx, 1000)
	require.NoError(t, err)

	updated, err := repo.Wallet(ctx, model.UpdateWallet{
		Uuid:   created.Uuid,
		Type:   model.OperationDeposit,
		Amount: 500,
	})
	require.NoError(t, err)
	require.Equal(t, int64(1500), updated.Amount)

	after, err := repo.Get(ctx, created.Uuid)
	require.NoError(t, err)
	require.Equal(t, int64(1500), after.Amount)
}

func TestWallet_Withdraw(t *testing.T) {
	repo := walletRepo.NewWalletRepository(testDB)
	ctx := context.Background()

	created, err := repo.Store(ctx, 1000)
	require.NoError(t, err)

	updated, err := repo.Wallet(ctx, model.UpdateWallet{
		Uuid:   created.Uuid,
		Type:   model.OperationWithdraw,
		Amount: 400,
	})
	require.NoError(t, err)
	require.Equal(t, int64(600), updated.Amount)

	after, err := repo.Get(ctx, created.Uuid)
	require.NoError(t, err)
	require.Equal(t, int64(600), after.Amount)
}

func TestWallet_WithdrawInsufficient(t *testing.T) {
	repo := walletRepo.NewWalletRepository(testDB)
	ctx := context.Background()

	created, err := repo.Store(ctx, 1000)
	require.NoError(t, err)

	_, err = repo.Wallet(ctx, model.UpdateWallet{
		Uuid:   created.Uuid,
		Type:   model.OperationWithdraw,
		Amount: 1500,
	})
	require.Error(t, err)

	after, err := repo.Get(ctx, created.Uuid)
	require.NoError(t, err)
	require.Equal(t, int64(1000), after.Amount)
}

func TestWallet_UnknownOperation(t *testing.T) {
	repo := walletRepo.NewWalletRepository(testDB)
	ctx := context.Background()

	created, err := repo.Store(ctx, 1000)
	require.NoError(t, err)

	_, err = repo.Wallet(ctx, model.UpdateWallet{
		Uuid:   created.Uuid,
		Type:   "UNKNOWN",
		Amount: 100,
	})
	require.Error(t, err)
}

func TestWallet_NotFound(t *testing.T) {
	repo := walletRepo.NewWalletRepository(testDB)
	ctx := context.Background()

	_, err := repo.Wallet(ctx, model.UpdateWallet{
		Uuid:   "00000000-0000-0000-0000-000000000000",
		Type:   model.OperationDeposit,
		Amount: 100,
	})
	require.Error(t, err)
}

func TestWallet_InvalidUUID(t *testing.T) {
	repo := walletRepo.NewWalletRepository(testDB)
	ctx := context.Background()

	_, err := repo.Wallet(ctx, model.UpdateWallet{
		Uuid:   "not-a-uuid",
		Type:   model.OperationDeposit,
		Amount: 100,
	})
	require.Error(t, err)
}

func TestWallet_ConcurrentWithdraw(t *testing.T) {
	repo := walletRepo.NewWalletRepository(testDB)
	ctx := context.Background()

	created, err := repo.Store(ctx, 1000)
	require.NoError(t, err)

	const workers = 100
	const amount = int64(10)

	var wg sync.WaitGroup
	var success atomic.Int64
	var failed atomic.Int64

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			_, err := repo.Wallet(ctx, model.UpdateWallet{
				Uuid:   created.Uuid,
				Type:   model.OperationWithdraw,
				Amount: amount,
			})
			if err != nil {
				failed.Add(1)
				return
			}
			success.Add(1)
		}()
	}
	wg.Wait()

	after, err := repo.Get(ctx, created.Uuid)
	require.NoError(t, err)

	require.Equal(t, int64(0), after.Amount)
	require.Equal(t, int64(100), success.Load())
	require.Equal(t, int64(0), failed.Load())
}
