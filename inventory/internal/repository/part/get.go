package part

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/Alexey-step/rocket-factory/inventory/internal/model"
	"github.com/Alexey-step/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/Alexey-step/rocket-factory/inventory/internal/repository/model"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
)

func (r *repository) GetPart(ctx context.Context, orderUUID string) (model.Part, error) {
	var part repoModel.Part
	err := r.collection.FindOne(ctx, bson.M{"uuid": orderUUID}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Error(ctx, "Part not found in MongoDB", zap.String("uuid", orderUUID))
			return model.Part{}, model.ErrPartNotFound
		}
		logger.Error(ctx, "Error finding part in MongoDB", zap.String("uuid", orderUUID), zap.Error(err))
		return model.Part{}, err
	}

	logger.Info(ctx, "Part found in MongoDB", zap.String("uuid", orderUUID))
	return converter.PartToModel(part), nil
}
