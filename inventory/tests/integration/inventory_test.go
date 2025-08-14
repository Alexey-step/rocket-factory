//go:build integration

package integration

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
	inventoryV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/inventory/v1"
)

var _ = Describe("InventoryService", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient inventoryV1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		// Создаём gRPC клиент
		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		inventoryClient = inventoryV1.NewInventoryServiceClient(conn)
	})

	AfterEach(func() {
		// Чистим коллекцию после теста
		err := env.ClearPartsCollection(ctx)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешную очистку коллекции parts")

		cancel()
	})

	Describe("GetPart", func() {
		var partUUID string

		BeforeEach(func() {
			// Вставляем тестовое наблюдение
			var err error
			partUUID, err = env.InsertTestPart(ctx)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестового детали в MongoDB")
		})

		It("должен успешно возвращать часть по UUID", func() {
			resp, err := inventoryClient.GetPart(ctx, &inventoryV1.GetPartRequest{
				Uuid: partUUID,
			})

			if err != nil {
				logger.Error(ctx, "Ошибка при вызове GetPart", zap.Error(err))
				// Не завершаем тест сразу, чтобы увидеть больше информации
			} else {
				logger.Info(ctx, "GetPart выполнен успешно")
			}

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetPart()).ToNot(BeNil())
			Expect(resp.GetPart().Uuid).To(Equal(partUUID))
			Expect(resp.GetPart().Name).ToNot(BeNil())
			Expect(resp.GetPart().Description).ToNot(BeEmpty())
			Expect(resp.GetPart().GetCreatedAt()).ToNot(BeNil())
		})
	})

	Describe("ListParts", func() {
		var partUUIDs []string

		BeforeEach(func() {
			// Вставляем тестовое детали в бд
			var err error
			partUUIDs, err = env.InsertTestParts(ctx, 3)
			Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестовых деталей в MongoDB")
		})

		It("должен успешно возвращать список деталей по переданным фильтрам", func() {
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Uuids: partUUIDs,
				},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.GetParts()).ToNot(BeNil())
			Expect(len(resp.GetParts())).Should(Equal(len(partUUIDs)))
		})

		It("должен возвращать пустой список деталей", func() {
			fakePartsUUIDs := env.GetFakePartUUIDS()
			resp, err := inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Uuids: fakePartsUUIDs,
				},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(len(resp.GetParts())).Should(Equal(0))
		})
	})
})
