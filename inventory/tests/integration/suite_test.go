//go:build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
)

const testTimeout = 5 * time.Minute

var (
	env *TestEnvironment

	suiteCtx    context.Context
	suiteCancel context.CancelFunc
)

// testLoggerConfig implements logger.LoggerConfig for integration tests
type testLoggerConfig struct {
	level                  string
	asJson                 bool
	otelEnabled            bool
	otelServiceName        string
	otelServiceEnvironment string
	otelEndpoint           string
}

func (c *testLoggerConfig) Level() string                  { return c.level }
func (c *testLoggerConfig) AsJson() bool                   { return c.asJson }
func (c *testLoggerConfig) OtelEnabled() bool              { return c.otelEnabled }
func (c *testLoggerConfig) OtelServiceName() string        { return c.otelServiceName }
func (c *testLoggerConfig) OtelServiceEnvironment() string { return c.otelServiceEnvironment }
func (c *testLoggerConfig) OtelEndpoint() string           { return c.otelEndpoint }

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Inventory Service Integration Test Suite")
}

var _ = BeforeSuite(func() {
	testLoggerConfig := &testLoggerConfig{
		level:                  loggerLevelValue,
		asJson:                 true,
		otelEnabled:            false,
		otelServiceName:        "inventory-integration-tests",
		otelServiceEnvironment: "test",
		otelEndpoint:           "",
	}

	err := logger.Init(testLoggerConfig)
	if err != nil {
		panic(fmt.Sprintf("не удалось инициализировать логгер: %v", err))
	}

	suiteCtx, suiteCancel = context.WithTimeout(context.Background(), testTimeout)

	// Загружаем .env файл и устанавливаем переменные в окружение
	envVars, err := godotenv.Read(filepath.Join("..", "..", "..", "deploy", "compose", "inventory", ".env"))
	if err != nil {
		logger.Error(suiteCtx, "Не удалось загрузить .env файл", zap.Error(err))
	}

	// Устанавливаем переменные в окружение процесса
	for key, value := range envVars {
		_ = os.Setenv(key, value)
	}

	logger.Info(suiteCtx, "Запуск тестового окружения...")
	env = setupTestEnvironment(suiteCtx)
})

var _ = AfterSuite(func() {
	logger.Info(context.Background(), "Завершение набора тестов")
	if env != nil {
		teardownTestEnvironment(suiteCtx, env)
	}
	suiteCancel()
})
