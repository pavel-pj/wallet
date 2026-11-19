package app

import (
	//"github.com/pavel-pj/auth/internal/closer"
	"context"
	"log"

	"github.com/pavel-pj/wallet/internal/config"
	env "github.com/pavel-pj/wallet/internal/config/env"

	grpcWalletApi "github.com/pavel-pj/wallet/internal/api/wallet"

	db "github.com/pavel-pj/wallet/internal/client/db"
	pg "github.com/pavel-pj/wallet/internal/client/db/pg"

	//env "github.com/pavel-pj/auth/internal/config/env"

	repository "github.com/pavel-pj/wallet/internal/repository"
	walletRepo "github.com/pavel-pj/wallet/internal/repository/wallet"

	service "github.com/pavel-pj/wallet/internal/service"
	walletServ "github.com/pavel-pj/wallet/internal/service/wallet"

	"github.com/pavel-pj/wallet/internal/closer"
)

type serviceProvider struct {
	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	httpConfig    config.HTTPConfig
	swaggerConfig config.SwaggerConfig

	pgClient db.Client

	walletRepository repository.WalletRepository
	walletService    service.WalletService
	walletImpl       *grpcWalletApi.WalletImplementation
}

func NewServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GetPGConfig() config.PGConfig {

	if s.pgConfig == nil {
		conf, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = conf
	}

	return s.pgConfig

}

func (s *serviceProvider) GetGRPCConfig() config.GRPCConfig {

	if s.grpcConfig == nil {
		conf, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}
		s.grpcConfig = conf
	}

	return s.grpcConfig
}

func (s *serviceProvider) GetHTTPConfig() config.HTTPConfig {

	if s.httpConfig == nil {
		conf, err := env.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = conf
	}

	return s.httpConfig

}

func (s *serviceProvider) GetSwaggerConfig() config.SwaggerConfig {

	if s.swaggerConfig == nil {
		conf, err := env.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}

		s.swaggerConfig = conf
	}

	return s.swaggerConfig

}

func (s *serviceProvider) GetDBClient(ctx context.Context) db.Client {
	if s.pgClient == nil {

		client, err := pg.New(context.Background(), s.GetPGConfig().DSN())

		if err != nil {
			log.Fatalf("failed to connect to database: %s", err.Error())
		}

		// Проверяем подключение
		err = client.DB().Ping(context.Background())
		if err != nil {
			log.Fatalf("failed to ping database: %s", err.Error())
		}

		closer.Add(func() error {
			// тк pool.Close не возвращает error - создаем обертку
			return client.Close()

		})

		s.pgClient = client

	}

	return s.pgClient
}

/*
	func (s *serviceProvider) GetTxManager(ctx context.Context) db.TxManager {
		if s.txManager == nil {
			s.txManager = transaction.NewTransationManager(s.GetDBClient(ctx).DB())
		}

		return s.txManager
	}
*/
func (s *serviceProvider) GetWalletRepository(ctx context.Context) repository.WalletRepository {

	if s.walletRepository == nil {
		s.walletRepository = walletRepo.NewWalletRepository(s.GetDBClient(ctx))
	}

	return s.walletRepository
}

func (s *serviceProvider) GetWalletService(ctx context.Context) service.WalletService {
	if s.walletService == nil {
		s.walletService = walletServ.NewWalletrService(
			s.GetWalletRepository(ctx),
		)
	}

	return s.walletService
}

func (s *serviceProvider) GetWalletImplimentation(ctx context.Context) *grpcWalletApi.WalletImplementation {
	if s.walletImpl == nil {
		s.walletImpl = grpcWalletApi.NewWalletImplimentation(s.GetWalletService(ctx))
	}

	return s.walletImpl
}
