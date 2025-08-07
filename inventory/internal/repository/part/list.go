package part

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/Alexey-step/rocket-factory/inventory/internal/model"
	"github.com/Alexey-step/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/Alexey-step/rocket-factory/inventory/internal/repository/model"
)

func (r *repository) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	mongoFilter := buildMongoFilter(filter)

	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, err
	}

	defer func() {
		cerr := cursor.Close(ctx)
		if cerr != nil {
			log.Printf("failed to close cursor: %v\n", cerr)
		}
	}()

	var result []repoModel.Part

	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, err
	}

	ok := matchesFilter(result, filter)
	if !ok {
		return nil, err
	}

	return converter.PartsToModel(result), nil
}

func buildMongoFilter(filter model.PartsFilter) bson.M {
	m := bson.M{}
	if len(filter.Uuids) > 0 {
		m["uuid"] = bson.M{"$in": filter.Uuids}
	}
	if len(filter.Names) > 0 {
		m["name"] = bson.M{"$in": filter.Names}
	}
	if len(filter.Categories) > 0 {
		m["category"] = bson.M{"$in": filter.Categories}
	}
	if len(filter.ManufacturerCountries) > 0 {
		m["manufacturer.country"] = bson.M{"$in": filter.ManufacturerCountries}
	}
	if len(filter.Tags) > 0 {
		m["tags"] = bson.M{"$all": filter.Tags}
	}
	return m
}

func matchesFilter(result []repoModel.Part, filter model.PartsFilter) bool {
	if len(filter.Uuids) > 0 && len(result) != len(filter.Uuids) {
		return false
	}

	if len(filter.Names) > 0 && len(result) != len(filter.Names) {
		return false
	}

	if len(filter.ManufacturerCountries) > 0 && len(result) != len(filter.ManufacturerCountries) {
		return false
	}

	if len(filter.Categories) > 0 && len(result) != len(filter.Categories) {
		return false
	}

	return true
}
