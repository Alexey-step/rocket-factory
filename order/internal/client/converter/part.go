package converter

import (
	"log"
	"time"

	"github.com/samber/lo"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
	inventory_v1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/inventory/v1"
)

func PartListToModel(parts []*inventory_v1.Part) []model.Part {
	modelParts := make([]model.Part, len(parts))
	for _, part := range parts {
		modelParts = append(modelParts, PartToModel(part))
	}
	return modelParts
}

func PartToModel(part *inventory_v1.Part) model.Part {
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

func DimensionsToModel(dimensions *inventory_v1.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
	}
}

func ManufacturerToModel(manufacturer *inventory_v1.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func MetadataToModel(metadata map[string]*inventory_v1.Value) model.Metadata {
	res := model.Metadata{}

	for _, value := range metadata {
		if value == nil {
			continue
		}

		switch v := value.Kind.(type) {
		case *inventory_v1.Value_StringValue:
			res.StringValue = lo.ToPtr(v.StringValue)
		case *inventory_v1.Value_Int64Value:
			res.Int64Value = lo.ToPtr(v.Int64Value)
		case *inventory_v1.Value_BoolValue:
			res.BoolValue = lo.ToPtr(v.BoolValue)
		case *inventory_v1.Value_DoubleValue:
			res.DoubleValue = lo.ToPtr(v.DoubleValue)
		default:
			log.Printf("unknown metadata metadata type: %T", value)
		}
	}
	return res
}

func PartsFilterToProto(filter model.PartsFilter) *inventory_v1.PartsFilter {
	categories := make([]inventory_v1.Category, len(filter.Categories))

	if len(filter.Categories) > 0 {
		for _, category := range filter.Categories {
			categories = append(categories, categoryToProto(category))
		}
	}

	return &inventory_v1.PartsFilter{
		Uuids:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.ManufacturerCountries,
	}
}

func categoryToProto(category model.Category) inventory_v1.Category {
	switch category {
	case model.CategoryEngine:
		return inventory_v1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return inventory_v1.Category_CATEGORY_FUEL
	case model.CategoryPorthole:
		return inventory_v1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return inventory_v1.Category_CATEGORY_WING
	default:
		return inventory_v1.Category_CATEGORY_UNSPECIFIED
	}
}
