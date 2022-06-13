package command

import (
	"context"
	"fmt"
	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/config"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/graph"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/graph/generated"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/handler"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/infrastracture/loader"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/infrastracture/persistence"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/io"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/server"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/version"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	exitOk    = 0
	exitError = 1
)

func Run() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	// init logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup logger: %s\n", err)
		return exitError
	}
	defer logger.Sync()
	logger = logger.With(zap.String("version", version.Version))

	// load config
	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		logger.Error("failed to load config", zap.Error(err))
		return exitError
	}

	// init listener
	listener, err := net.Listen("tcp", cfg.Address())
	if err != nil {
		logger.Error("failed to listen port", zap.Int("port", cfg.Port), zap.Error(err))
		return exitError
	}
	logger.Info("server start listening", zap.Int("port", cfg.Port))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// init mysql
	logger.Info("connect to mysql ", zap.String("DSN", cfg.DB.DSN))
	sqlSetting := &config.SQLDBSettings{
		SqlDsn:              cfg.DB.DSN,
		SqlMaxOpenConns:     cfg.DB.MaxOpenConns,
		SqlMaxIdleConns:     cfg.DB.MaxIdleConns,
		SqlConnsMaxLifetime: cfg.DB.ConnsMaxLifetime,
	}

	mysqlDatabase, err := io.NewDatabase(sqlSetting)
	if err != nil {
		logger.Error("failed to create mysql db repository", zap.Error(err), zap.String("DSN", cfg.DB.DSN))
		return exitError
	}

	repositories, err := persistence.NewRepositories(mysqlDatabase)
	if err != nil {
		logger.Error("failed to new repositories", zap.Error(err))
		return exitError
	}

	// init loader
	loaders := loader.NewLoaders(repositories)
	// init to start http server
	// init gql server
	query := gqlhandler.NewDefaultServer(generated.NewExecutableSchema(
		generated.Config{
			Resolvers: &graph.Resolver{
				Repo:    repositories,
				Loaders: loaders,
			},
		}))
	registry := handler.NewHandler(logger, repositories, query, version.Version)
	httpServer := server.NewServer(registry, &server.Config{Log: logger})
	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(func() error {
		return httpServer.Serve(listener)
	})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	select {
	case <-sigCh:
	case <-ctx.Done():
	}

	if err := httpServer.GracefulShutdown(ctx); err != nil {
		return exitError
	}

	cancel()
	if err := wg.Wait(); err != nil {
		return exitError
	}

	return exitOk
}
