package part

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Alexey-step/rocket-factory/inventory/internal/model"
	"github.com/Alexey-step/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/Alexey-step/rocket-factory/inventory/internal/repository/model"
)

func (r *repository) GetPart(ctx context.Context, orderUUID string) (model.Part, error) {
	var part repoModel.Part
	err := r.collection.FindOne(ctx, bson.M{"uuid": orderUUID}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Part{}, model.ErrPartNotFound
		}
		return model.Part{}, err
	}

	return converter.PartToModel(part), nil
}
