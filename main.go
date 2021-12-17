package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rusrafkasimov/history/internal/ch_database"
	"github.com/rusrafkasimov/history/internal/config"
	"github.com/rusrafkasimov/history/internal/grpc_server"
	"github.com/rusrafkasimov/history/internal/logger"
	"github.com/rusrafkasimov/history/internal/queue"
	"github.com/rusrafkasimov/history/internal/trace"
	"github.com/rusrafkasimov/history/internal/vault"
	"github.com/rusrafkasimov/history/pkg/delivery/router"
	"github.com/rusrafkasimov/history/proto"
	"time"
)

// @title Swagger Account History Service
// @version 0.1
// @description Account History Microservice (Golang)

// @contact.name Ruslan Kasimov

// @host 127.0.0.1:8091
// @BasePath /

// @securityDefinitions.apikey TokenJWT
// @in header
// @name Authorization

const (
	Name           = "History"
	contextKeyName = "Name"
	serverTimeout  = 10 * time.Second
)

func main() {
	ctx := context.Background()

	id, err := gonanoid.New()
	if err != nil {
		fmt.Println("Can't generate new node ID")
		return
	}

	ctx = context.WithValue(
		ctx,
		contextKeyName,
		Name+"_"+id,
	)

	var env string

	flag.StringVar(&env, "env", ".env.local", "Environment Variables filename")
	flag.Parse()

	// Load service configuration from environment
	if err := config.LoadConfig(env); err != nil {
		fmt.Printf("Error: can't load env. %" + err.Error())
	}

	// Initialize vault
	vaultProvider := vault.NewVaultProvider()
	appConfig := config.NewConfig(vaultProvider)

	// Initialize logger
	loki, err := logger.NewLogger(Name, "api", appConfig)
	if err != nil {
		fmt.Println("Error while connect to loki")
	}

	// Initialize tracing
	closer, err := trace.InitJaegerTracing(ctx, contextKeyName, appConfig)
	if err != nil {
		loki.Errorf("Error while init tracing")
	}
	defer closer.Close()

	// Initialize Database
	chDB, err := ch_database.InitializeDB(appConfig, loki)
	if err != nil {
		loki.Errorf("Error init database")
	}

	// Initialize NATS Queue
	newQueue, err := queue.NewQueue(ctx, loki, appConfig)
	if err != nil {
		loki.Errorf("Error init new queue")
	}

	// Migrate model
	// ch_database.MigrateHistoryModels(chDB)

	// Build context
	repoCtx := router.BuildRepositoryContext(chDB)
	ucCtx := router.BuildUseCaseContext(newQueue, repoCtx, loki)
	appCtx := router.BuildApplicationContext(ucCtx, loki)

	// Subscribe on NATS events
	historyCh, err := ucCtx.HistoryUse.SubscribeOnEvents(ctx)
	if err != nil {
		loki.Errorf("subscribe on history events")
	}

	// Events receiver
	ucCtx.HistoryUse.RunReceiveEventsLoop(ctx, historyCh)

	// Load gRPC server configuration
	gRPCServer := grpc_server.NewGRPCServer(appConfig, loki)

	// Run gRPC server
	go gRPCServer.RunGrpcServer(func(s *grpc_server.GRPCServerConfig) {
		operation_history.RegisterOperationHistoryServer(s, appCtx.RPCHistoryController)
	})

	// Initialize gin routes and run server
	rGin := gin.Default()
	gin.ForceConsoleColor()
	router.MapUrl(rGin, appCtx)

	httpHost, err := appConfig.Get("HTTP_HOST")
	if err != nil {
		loki.Errorf(err.Error())
	}

	httpPort := ":8090"

	if err = rGin.Run(httpPort); err != nil {
		loki.Errorf("Error: can't start GIN router. %s", err.Error())
	}

	loki.Infof("Upstream started at %v", httpHost+httpPort)

	defer func(){
		loki.Infof("History service stopped")
		_ = newQueue.Close()
	}()
}
