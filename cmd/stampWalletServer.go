package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/StampWallet/backend/internal/api"
	handlers "github.com/StampWallet/backend/internal/api/handlers"
	"github.com/StampWallet/backend/internal/config"
	"github.com/StampWallet/backend/internal/database"
	accessors "github.com/StampWallet/backend/internal/database/accessors"
	"github.com/StampWallet/backend/internal/managers"
	"github.com/StampWallet/backend/internal/middleware"
	"github.com/StampWallet/backend/internal/services"
)

// Creates server from config
func createServer(config config.Config) (*api.APIServer, error) {
	db, err := services.GetDatabase(config)
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %+v", err)
	}

	logger := log.New(os.Stderr, "StampWalletServer", 0)
	logger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	baseServices := services.BaseServices{
		Logger:   logger,
		Database: db,
	}

	tokenService := services.CreateTokenServiceImpl(baseServices.NewPrefix("TokenServiceImpl"))
	emailService, err := services.CreateEmailServiceImpl(config.SmtpConfig,
		services.NewPrefix(logger, "EmailServiceImpl"))
	if err != nil {
		return nil, fmt.Errorf("failed to create emailService: %+v", err)
	}
	fileStorageService, err := services.CreateFileStorageServiceImpl(
		baseServices.NewPrefix("FileStorageService"),
		config.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create fileStorageService: %+v", err)
	}

	authMiddleware := middleware.CreateAuthMiddleware(services.NewPrefix(logger, "AuthMiddleware"), tokenService)
	requireValidEmailMiddleware := middleware.CreateRequireValidEmailMiddleware(
		services.NewPrefix(logger, "RequireValidEmailMiddleware"))

	authManager := managers.CreateAuthManagerImpl(baseServices, emailService, tokenService)
	virtualCardManager := managers.CreateVirtualCardManagerImpl(baseServices)
	itemDefinitionManager := managers.CreateItemDefinitionManagerImpl(baseServices, fileStorageService)
	localCardManager := managers.CreateLocalCardManagerImpl(baseServices)
	businessManager := managers.CreateBusinessManagerImpl(baseServices, fileStorageService)
	transactionManager := managers.CreateTransactionManagerImpl(baseServices)

	userAuthorizedAcessor := accessors.CreateUserAuthorizedAccessorImpl(baseServices.Database)
	businessAuthorizedAccessor := accessors.CreateBusinessAuthorizedAccessorImpl(baseServices.Database)
	authorizedTransactionAccessor := accessors.CreateAuthorizedTransactionAccessorImpl(baseServices.Database)

	//idk if this shouldnt be handled by CreateAPIServer
	handlers := api.APIHandlers{
		AuthHandlers: handlers.CreateAuthHandlers(authManager, services.NewPrefix(logger, "AuthHandlers")),
		BusinessHandlers: handlers.CreateBusinessHandlers(
			businessManager,
			transactionManager,
			itemDefinitionManager,

			userAuthorizedAcessor,
			businessAuthorizedAccessor,
			authorizedTransactionAccessor,

			services.NewPrefix(logger, "BusinessHandlers"),
		),
		UserHandlers: handlers.CreateUserHandlers(
			virtualCardManager,
			localCardManager,
			businessManager,
			transactionManager,
			itemDefinitionManager,
			userAuthorizedAcessor,
			authorizedTransactionAccessor,
			services.NewPrefix(logger, "UserHandlers"),
		),
		FileHandlers: handlers.CreateFileHandlers(
			fileStorageService,
			services.NewPrefix(logger, "FileHandlers"),
			userAuthorizedAcessor,
		),
	}

	server := api.CreateAPIServer(authMiddleware, requireValidEmailMiddleware, &handlers,
		services.NewPrefix(logger, "APIServer"), config)

	return server, nil
}

// Entrypoint of the backend app
func main() {
	// CLI framework configuration
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "config", Value: "stamp_wallet_server.yaml"},
		},
		Commands: []*cli.Command{
			{
				Name:  "start",
				Usage: "starts the server",
				Action: func(ctx *cli.Context) error {
					config, err := config.LoadConfig(ctx.String("config"))
					if err != nil {
						return fmt.Errorf("failed to load config: %+v", err)
					}

					server, err := createServer(config)
					if err != nil {
						return err
					}
					err = server.Start()
					if err != nil {
						return fmt.Errorf("failed to start server: %+v", err)
					}
					return nil
				},
			},
			{
				Name:  "automigrate",
				Usage: "automatically migrates schema",
				Action: func(ctx *cli.Context) error {
					config, err := config.LoadConfig(ctx.String("config"))
					if err != nil {
						return fmt.Errorf("failed to load config: %+v", err)
					}

					db, err := services.GetDatabase(config)
					if err != nil {
						return fmt.Errorf("failed to get database: %+v", err)
					}
					err = database.AutoMigrate(db)

					if err != nil {
						return fmt.Errorf("failed to automigrate: %+v", err)
					}
					return nil
				},
			},
			{
				Name:  "example-config",
				Usage: "creates/replaces config file with example values",
				Action: func(ctx *cli.Context) error {
					return config.SaveConfig(config.GetDefaultConfig(), ctx.String("config"))
				},
			},
			{
				Name: "send-email",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "to"},
					&cli.StringFlag{Name: "subject"},
					&cli.StringFlag{Name: "body"},
				},
				Usage: "sends an email",
				Action: func(ctx *cli.Context) error {
					config, err := config.LoadConfig(ctx.String("config"))
					emailService, err := services.CreateEmailServiceImpl(config.SmtpConfig, log.Default())
					if err != nil {
						return fmt.Errorf("failed to create email service %+v", err)
					}
					return emailService.Send(ctx.String("to"), ctx.String("subject"), ctx.String("body"))
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
