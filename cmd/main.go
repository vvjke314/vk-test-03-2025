package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/vvjke314/vk-test-03-2025/config"
	"github.com/vvjke314/vk-test-03-2025/internal/logger"
	"github.com/vvjke314/vk-test-03-2025/internal/repository"
	"github.com/vvjke314/vk-test-03-2025/internal/usecases"
	api "github.com/vvjke314/vk-test-03-2025/pkg/routes"
)

func main() {
	// cfg init
	loader := config.NewLoader()
	loader.Load()

	// repository config init
	repoCfg := config.NewTnConfig()

	// logger init
	appLogger, err := logger.NewSimpleLogger("application.log")
	if err != nil {
		log.Fatalf("error while getting log while: %v", err)
	}
	defer appLogger.Close()

	// init tarantool repository
	repo := repository.NewTnRepository()
	ctx := context.Background()
	if err := repo.Init(ctx, repoCfg, appLogger); err != nil {
		appLogger.Error(fmt.Sprintf("error while initing repository: %v", err))
		log.Fatalf("error initing repository: %v", err)
	}
	defer repo.Close()

	// use cases initsialize
	uc := usecases.NewKeyValueUseCase(repo)

	// HTTP setting up
	r := api.SetupRoutes(uc)

	// server start
	server := &http.Server{
		Addr:    ":" + "8080",
		Handler: r,
	}

	appLogger.Info("server is up, on port :8080")
	log.Printf("server is up, on port :8080")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		appLogger.Error(fmt.Sprintf("server error: %v", err))
		log.Fatalf("server error: %v", err)
	}
}
