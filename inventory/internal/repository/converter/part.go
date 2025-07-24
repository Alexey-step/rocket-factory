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

func PartMetadataToModel(metadata repoModel.Metadata) model.Metadata {
	return model.Metadata{
		StringValue: metadata.StringValue,
		Int64Value:  metadata.Int64Value,
		DoubleValue: metadata.DoubleValue,
		BoolValue:   metadata.BoolValue,
	}
}
