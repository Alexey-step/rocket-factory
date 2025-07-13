package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Alexey-step/rocket-factory/inventory/internal/model"
	inventory_v1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/inventory/v1"
)

func PartToProto(part model.Part) *inventory_v1.Part {
	var updatedAt *timestamppb.Timestamp
	if part.UpdatedAt != nil {
		updatedAt = timestamppb.New(*part.UpdatedAt)
	}

	return &inventory_v1.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		StockQuantity: part.StockQuantity,
		Price:         part.Price,
		Metadata:      metadataToProto(part.Metadata),
		Category:      CategoryToProto(part.Category),
		Manufacturer:  manufacturerToProto(part.Manufacturer),
		CreatedAt:     timestamppb.New(part.CreatedAt),
		UpdatedAt:     updatedAt,
	}
}

func CategoryToProto(category model.Category) inventory_v1.Category {
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

func manufacturerToProto(m model.Manufacturer) *inventory_v1.Manufacturer {
	return &inventory_v1.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func metadataToProto(meta model.Metadata) map[string]*inventory_v1.Value {
	var val *inventory_v1.Value
	switch {
	case meta.StringValue != nil:
		val = &inventory_v1.Value{
			Kind: &inventory_v1.Value_StringValue{StringValue: *meta.StringValue},
		}
	case meta.Int64Value != nil:
		val = &inventory_v1.Value{
			Kind: &inventory_v1.Value_Int64Value{Int64Value: *meta.Int64Value},
		}
	case meta.DoubleValue != nil:
		val = &inventory_v1.Value{
			Kind: &inventory_v1.Value_DoubleValue{DoubleValue: *meta.DoubleValue},
		}
	case meta.BoolValue != nil:
		val = &inventory_v1.Value{
			Kind: &inventory_v1.Value_BoolValue{BoolValue: *meta.BoolValue},
		}
	default:
		val = &inventory_v1.Value{}
	}
	return map[string]*inventory_v1.Value{"value": val}
}

func PartsFilterToModel(filter *inventory_v1.PartsFilter) model.PartsFilter {
	partsUUIDs := make([]string, 0, len(filter.Uuids))
	if len(filter.Uuids) > 0 {
		partsUUIDs = filter.Uuids
	}

	partsNames := make([]string, 0, len(filter.Names))
	if len(filter.Names) > 0 {
		partsUUIDs = filter.Names
	}

	partsTags := make([]string, 0, len(filter.Tags))
	if len(filter.Tags) > 0 {
		partsUUIDs = filter.Tags
	}

	partsManufacturerCountries := make([]string, 0, len(filter.ManufacturerCountries))
	if len(filter.ManufacturerCountries) > 0 {
		partsUUIDs = filter.ManufacturerCountries
	}

	partsCategories := make([]model.Category, 0, len(filter.Categories))
	if len(filter.Categories) > 0 {
		for _, c := range filter.Categories {
			partsCategories = append(partsCategories, model.Category(c.String()))
		}
	}

	return model.PartsFilter{
		Uuids:                 partsUUIDs,
		Names:                 partsNames,
		Categories:            partsCategories,
		ManufacturerCountries: partsManufacturerCountries,
		Tags:                  partsTags,
	}
}

func PartsToProto(parts []model.Part) []*inventory_v1.Part {
	result := make([]*inventory_v1.Part, 0, len(parts))
	for _, part := range parts {
		result = append(result, PartToProto(part))
	}
	return result
}

func CategoryToModel(category inventory_v1.Category) model.Category {
	switch category {
	case inventory_v1.Category_CATEGORY_ENGINE:
		return model.CategoryEngine
	case inventory_v1.Category_CATEGORY_FUEL:
		return model.CategoryFuel
	case inventory_v1.Category_CATEGORY_PORTHOLE:
		return model.CategoryPorthole
	case inventory_v1.Category_CATEGORY_WING:
		return model.CategoryWing
	default:
		return model.CategoryUnspecified
	}
}
