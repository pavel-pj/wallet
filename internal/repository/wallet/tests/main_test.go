//go:build integration

package tests

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	db "github.com/pavel-pj/wallet/internal/client/db"
	pgclient "github.com/pavel-pj/wallet/internal/client/db/pg"
)

var testDB db.Client

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := postgres.Run(ctx, "postgres:18",
		postgres.WithDatabase("wallet_test"),
		postgres.WithUsername("wallet"),
		postgres.WithPassword("wallet"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Fatalf("start postgres container: %v", err)
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = container.Terminate(ctx)
		log.Fatalf("get dsn: %v", err)
	}

	if err := applyMigrations(ctx, dsn); err != nil {
		_ = container.Terminate(ctx)
		log.Fatalf("apply migrations: %v", err)
	}

	testDB, err = pgclient.New(ctx, dsn)
	if err != nil {
		_ = container.Terminate(ctx)
		log.Fatalf("connect to database: %v", err)
	}

	code := m.Run()

	if testDB != nil {
		_ = testDB.Close()
	}
	_ = container.Terminate(ctx)

	os.Exit(code)
}

func applyMigrations(ctx context.Context, dsn string) error {
	path, err := filepath.Abs("../../../../migrations/20261118092059_create_wallets_table.sql")
	if err != nil {
		return err
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Берём только Up-часть, иначе DROP TABLE из Down удалит таблицу.
	up := string(raw)
	if i := strings.Index(up, "-- +goose Down"); i >= 0 {
		up = up[:i]
	}

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, up)
	return err
}
