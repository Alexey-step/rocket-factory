package app

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/Alexey-step/rocket-factory/assembly/internal/config"
	assemblyMetrics "github.com/Alexey-step/rocket-factory/assembly/internal/metrics"
	"github.com/Alexey-step/rocket-factory/platform/pkg/closer"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
	platformMetrics "github.com/Alexey-step/rocket-factory/platform/pkg/metrics"
)

type App struct {
	diContainer *diContainer
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	// Канал для ошибок от компонентов
	errCh := make(chan error, 2)

	// Контекст для остановки всех горутин
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Консьюмер
	go func() {
		if err := a.runConsumer(ctx); err != nil {
			errCh <- errors.Errorf("consumer crashed: %v", err)
		}
	}()

	// Ожидание либо ошибки, либо завершения контекста (например, сигнал SIGINT/SIGTERM)
	select {
	case <-ctx.Done():
		logger.Info(ctx, "Shutdown signal received")
	case err := <-errCh:
		logger.Error(ctx, "Component crashed, shutting down", zap.Error(err))
		// Триггерим cancel, чтобы остановить второй компонент
		cancel()
		// Дождись завершения всех задач (если есть graceful shutdown внутри)
		<-ctx.Done()
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initMetrics,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(config.AppConfig().Logger) //nolint:contextcheck
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 OrderPaid Kafka consumer running")

	err := a.diContainer.OrderConsumerService().RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initMetrics(ctx context.Context) error {
	err := platformMetrics.InitProvider(ctx, config.AppConfig().Metrics)
	if err != nil {
		panic(fmt.Sprintf("failed to init metrics provider: %v", err))
	}

	closer.AddNamed("Metrics Provider", platformMetrics.Shutdown)

	err = assemblyMetrics.InitMetrics(config.AppConfig().Metrics)
	if err != nil {
		panic(fmt.Sprintf("failed to init Assembly metrics: %v", err))
	}

	return nil
}
