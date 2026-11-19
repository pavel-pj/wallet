//go:build integration

package tests

import (
	"context"
	"testing"

	walletRepo "github.com/pavel-pj/wallet/internal/repository/wallet"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	repo := walletRepo.NewWalletRepository(testDB)

	w, err := repo.Store(context.Background(), 1000)
	require.NoError(t, err)
	require.NotEmpty(t, w.Uuid)
	require.Equal(t, int64(1000), w.Amount)
}
