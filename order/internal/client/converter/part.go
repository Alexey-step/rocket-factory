package converter

import (
	"log"
	"time"

	"github.com/samber/lo"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
	inventoryV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/inventory/v1"
)

func PartListToModel(parts []*inventoryV1.Part) []model.Part {
	var modelParts []model.Part
	for _, part := range parts {
		modelPart := PartToModel(part)
		modelParts = append(modelParts, modelPart)
	}
	return modelParts
}

func PartToModel(part *inventoryV1.Part) model.Part {
	var updatedAt *time.Time
	if part.UpdatedAt != nil {
		updatedAt = lo.ToPtr(part.UpdatedAt.AsTime())
	}

	return model.Part{
		UUID:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      model.Category(part.Category),
		Dimensions:    DimensionsToModel(part.Dimensions),
		Manufacturer:  ManufacturerToModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      MetadataToModel(part.Metadata),
		CreatedAt:     part.CreatedAt.AsTime(),
		UpdatedAt:     updatedAt,
	}
}

func DimensionsToModel(dimensions *inventoryV1.Dimensions) model.Dimensions {
	if dimensions == nil {
		return model.Dimensions{}
	}
	return model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ManufacturerToModel(manufacturer *inventoryV1.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func MetadataToModel(metadata map[string]*inventoryV1.Value) model.Metadata {
	res := model.Metadata{}

	for _, value := range metadata {
		if value == nil {
			continue
		}

		switch v := value.Kind.(type) {
		case *inventoryV1.Value_StringValue:
			res.StringValue = lo.ToPtr(v.StringValue)
		case *inventoryV1.Value_Int64Value:
			res.Int64Value = lo.ToPtr(v.Int64Value)
		case *inventoryV1.Value_BoolValue:
			res.BoolValue = lo.ToPtr(v.BoolValue)
		case *inventoryV1.Value_DoubleValue:
			res.DoubleValue = lo.ToPtr(v.DoubleValue)
		default:
			log.Printf("unknown metadata metadata type: %T", value)
		}
	}
	return res
}

func PartsFilterToProto(filter model.PartsFilter) *inventoryV1.PartsFilter {
	categories := make([]inventoryV1.Category, len(filter.Categories))

	if len(filter.Categories) > 0 {
		for _, category := range filter.Categories {
			categories = append(categories, categoryToProto(category))
		}
	}

	filters := &inventoryV1.PartsFilter{
		Uuids:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}

	return filters
}

func categoryToProto(category model.Category) inventoryV1.Category {
	switch category {
	case model.CategoryEngine:
		return inventoryV1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return inventoryV1.Category_CATEGORY_FUEL
	case model.CategoryPorthole:
		return inventoryV1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return inventoryV1.Category_CATEGORY_WING
	default:
		return inventoryV1.Category_CATEGORY_UNSPECIFIED
	}
}
