//go:build integration

package integration

import (
	"context"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	repoModel "github.com/Alexey-step/rocket-factory/inventory/internal/repository/model"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
)

// InsertTestPart — вставляет тестовое часть в коллекцию Mongo и возвращает ее UUID
func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	partUUID := gofakeit.UUID()

	partDoc := getDBPart(partUUID)

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertOne(ctx, partDoc)
	if err != nil {
		logger.Error(ctx, "Ошибка при вставке тестовой части", zap.Error(err))
		return "", err
	}

	logger.Info(ctx, "Тестовая часть успешно вставлена", zap.String("uuid", partUUID), zap.String("database", databaseName), zap.String("collection", partsCollectionName))
	return partUUID, nil
}

func (env *TestEnvironment) InsertTestParts(ctx context.Context, count int) ([]string, error) {
	var partDocs []interface{}
	var uuids []string

	for i := 0; i < count; i++ {
		partUUID := gofakeit.UUID()
		uuids = append(uuids, partUUID)

		partDoc := getDBPart(partUUID)
		partDocs = append(partDocs, partDoc)
	}

	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory"
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertMany(ctx, partDocs)
	if err != nil {
		logger.Error(ctx, "Ошибка при массовой вставке тестовых частей", zap.Error(err))
		return nil, err
	}

	logger.Info(ctx, "Тестовые части успешно вставлены", zap.Int("count", count), zap.String("database", databaseName), zap.String("collection", partsCollectionName))
	return uuids, nil
}

// GetFakePartUUIDS — возвращает сгенерированные UUID для тестовых частей
func (env *TestEnvironment) GetFakePartUUIDS() []string {
	return []string{gofakeit.UUID(), gofakeit.UUID(), gofakeit.UUID()}
}

func getDBPart(partUUID string) bson.M {
	now := time.Now()
	return bson.M{
		"uuid":           partUUID,
		"name":           gofakeit.Name(),
		"description":    gofakeit.Sentence(10),
		"price":          gofakeit.Float64Range(100, 10_000),
		"stock_quantity": int64(gofakeit.Number(1, 100)),
		"category":       repoModel.Category("ENGINE"),
		"dimensions": bson.M{
			"length": gofakeit.Float64Range(0.1, 10.0),
			"width":  gofakeit.Float64Range(0.1, 10.0),
			"height": gofakeit.Float64Range(0.1, 10.0),
			"weight": gofakeit.Float64Range(0.1, 10.0),
		},
		"manufacturer": bson.M{
			"name":    gofakeit.Company(),
			"country": gofakeit.Country(),
			"website": gofakeit.URL(),
		},
		"tags": []string{gofakeit.Word(), gofakeit.Word()},
		"metadata": map[string]repoModel.Metadata{
			"string": {
				StringValue: lo.ToPtr(gofakeit.Word()),
			},
			"int": {
				Int64Value: lo.ToPtr(gofakeit.Int64()),
			},
			"double": {
				DoubleValue: lo.ToPtr(gofakeit.Float64()),
			},
			"bool": {
				BoolValue: lo.ToPtr(gofakeit.Bool()),
			},
		},
		"created_at": primitive.NewDateTimeFromTime(now),
	}
}

// ClearPartsCollection — удаляет все записи из коллекции parts
func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}
