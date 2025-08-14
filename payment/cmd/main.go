package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/Alexey-step/rocket-factory/payment/internal/app"
	"github.com/Alexey-step/rocket-factory/payment/internal/config"
	"github.com/Alexey-step/rocket-factory/platform/pkg/closer"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
)

const configPath = "./deploy/compose/payment/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ failed to initialize application",
			zap.Error(err))
		return
	}

	if err := a.Run(appCtx); err != nil {
		logger.Error(appCtx, "❌ failed to run application",
			zap.Error(err))
		return
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ failed to close resources",
			zap.Error(err))
	}
}
