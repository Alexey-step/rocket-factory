package part

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Alexey-step/rocket-factory/inventory/internal/config"
	def "github.com/Alexey-step/rocket-factory/inventory/internal/repository"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *repository {
	collection := db.Collection("parts")

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "uuid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		panic("Failed to create indexes: " + err.Error())
	}

	s := &repository{
		collection: collection,
	}

	// Проверяем, нужно ли отключить инициализацию тестовых данных
	isDisabledPartsInit := config.AppConfig().Mongo.DisabledInitMockParts()
	if !isDisabledPartsInit {
		s.InitParts(ctx)
	}

	return s
}
