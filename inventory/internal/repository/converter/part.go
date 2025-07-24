package converter

import (
	"github.com/Alexey-step/rocket-factory/inventory/internal/model"
	repoModel "github.com/Alexey-step/rocket-factory/inventory/internal/repository/model"
)

func PartToModel(part repoModel.Part) model.Part {
	return model.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		StockQuantity: part.StockQuantity,
		Price:         part.Price,
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
		Category:      model.Category(part.Category),
		Dimensions:    PartDimensionsToModel(part.Dimensions),
		Manufacturer:  PartManufacturerToModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      PartMetadataToModel(part.Metadata),
		Description:   part.Description,
	}
}

func PartDimensionsToModel(dimensions repoModel.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func PartManufacturerToModel(manufacturer repoModel.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func PartMetadataToModel(metadata map[string]repoModel.Metadata) map[string]model.Metadata {
	result := make(map[string]model.Metadata)
	for k, v := range metadata {
		result[k] = model.Metadata{
			StringValue: v.StringValue,
			Int64Value:  v.Int64Value,
			DoubleValue: v.DoubleValue,
			BoolValue:   v.BoolValue,
		}
	}
	return result
}

func PartsToModel(parts []repoModel.Part) []model.Part {
	result := make([]model.Part, len(parts))
	for i, part := range parts {
		result[i] = PartToModel(part)
	}
	return result
}
