//go:build integration

package tests

import (
	"context"
	"testing"

	"github.com/pavel-pj/wallet/internal/model"
	walletRepo "github.com/pavel-pj/wallet/internal/repository/wallet"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	repo := walletRepo.NewWalletRepository(testDB)
	ctx := context.Background()

	created, err := repo.Store(ctx, 2000)
	require.NoError(t, err)

	got, err := repo.Get(ctx, created.Uuid)
	require.NoError(t, err)

	require.Equal(t, created.Uuid, got.Uuid)
	require.Equal(t, created.Amount, got.Amount)
	require.False(t, got.CreatedAt.IsZero())
	require.False(t, got.UpdatedAt.IsZero())
}

func TestGet_NotFound(t *testing.T) {
	repo := walletRepo.NewWalletRepository(testDB)
	ctx := context.Background()

	_, err := repo.Get(ctx, "00000000-0000-0000-0000-000000000000")
	require.ErrorIs(t, err, model.ErrWalletNotFound)
}

func TestGet_InvalidUUID(t *testing.T) {
	repo := walletRepo.NewWalletRepository(testDB)
	ctx := context.Background()

	_, err := repo.Get(ctx, "not-a-uuid")
	require.ErrorIs(t, err, model.ErrInvalidUUID)
}
