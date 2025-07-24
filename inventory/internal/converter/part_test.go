package converter

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Alexey-step/rocket-factory/inventory/internal/model"
	inventory_v1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/inventory/v1"
)

func TestPartToProto(t *testing.T) {
	part := getMockPart()

	protoPart := PartToProto(part)

	assert.Equal(t, part.UUID, protoPart.Uuid)
	assert.Equal(t, part.Name, protoPart.Name)
	assert.Equal(t, part.Description, protoPart.Description)
	assert.Equal(t, part.StockQuantity, protoPart.StockQuantity)
	assert.Equal(t, part.Price, protoPart.Price)
	assert.Equal(t, CategoryToProto(part.Category), protoPart.Category)
	assert.Equal(t, part.Manufacturer.Name, protoPart.Manufacturer.Name)
	assert.Equal(t, part.Manufacturer.Country, protoPart.Manufacturer.Country)
	assert.Equal(t, part.Manufacturer.Website, protoPart.Manufacturer.Website)
	assert.Equal(t, timestamppb.New(part.CreatedAt), protoPart.CreatedAt)
	assert.NotNil(t, protoPart.Metadata)
}

func TestCategoryToProto(t *testing.T) {
	assert.Equal(t, inventory_v1.Category_CATEGORY_ENGINE, CategoryToProto(model.CategoryEngine))
	assert.Equal(t, inventory_v1.Category_CATEGORY_FUEL, CategoryToProto(model.CategoryFuel))
	assert.Equal(t, inventory_v1.Category_CATEGORY_PORTHOLE, CategoryToProto(model.CategoryPorthole))
	assert.Equal(t, inventory_v1.Category_CATEGORY_WING, CategoryToProto(model.CategoryWing))
	assert.Equal(t, inventory_v1.Category_CATEGORY_UNSPECIFIED, CategoryToProto("UNKNOWN"))
}

func TestMetadataToProto(t *testing.T) {
	strVal := "test"
	meta := map[string]model.Metadata{"value": {StringValue: &strVal}}
	result := metadataToProto(meta)
	assert.Equal(t, &inventory_v1.Value{Kind: &inventory_v1.Value_StringValue{StringValue: strVal}}, result["value"])

	intVal := int64(42)
	meta = map[string]model.Metadata{"value": {Int64Value: &intVal}}
	result = metadataToProto(meta)
	assert.Equal(t, &inventory_v1.Value{Kind: &inventory_v1.Value_Int64Value{Int64Value: intVal}}, result["value"])

	doubleVal := 3.14
	meta = map[string]model.Metadata{"value": {DoubleValue: &doubleVal}}
	result = metadataToProto(meta)
	assert.Equal(t, &inventory_v1.Value{Kind: &inventory_v1.Value_DoubleValue{DoubleValue: doubleVal}}, result["value"])

	boolVal := true
	meta = map[string]model.Metadata{"value": {BoolValue: &boolVal}}
	result = metadataToProto(meta)
	assert.Equal(t, &inventory_v1.Value{Kind: &inventory_v1.Value_BoolValue{BoolValue: boolVal}}, result["value"])

	meta = map[string]model.Metadata{"value": {}}
	result = metadataToProto(meta)
	assert.Equal(t, &inventory_v1.Value{}, result["value"])
}

func getMockPart() model.Part {
	return model.Part{
		UUID:          gofakeit.UUID(),
		Name:          gofakeit.Name(),
		Description:   "Primary propulsion unit",
		Price:         gofakeit.Float64Range(100, 10_000),
		StockQuantity: int64(gofakeit.Number(1, 100)),
		Category:      "ENGINE",
		Dimensions: model.Dimensions{
			Width:  gofakeit.Float64Range(0.1, 10.0),
			Height: gofakeit.Float64Range(0.1, 10.0),
			Length: gofakeit.Float64Range(0.1, 10.0),
			Weight: gofakeit.Float64Range(0.1, 10.0),
		},
		Manufacturer: model.Manufacturer{
			Name:    gofakeit.Name(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags: []string{gofakeit.EmojiTag(), gofakeit.EmojiTag()},
		Metadata: map[string]model.Metadata{
			"string": {StringValue: lo.ToPtr(gofakeit.Word())},
			"int":    {Int64Value: lo.ToPtr(gofakeit.Int64())},
			"double": {DoubleValue: lo.ToPtr(gofakeit.Float64())},
			"bool":   {BoolValue: lo.ToPtr(gofakeit.Bool())},
		},
		CreatedAt: timestamppb.Now().AsTime(),
	}
}
