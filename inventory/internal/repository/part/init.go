package part

import (
	"context"
	"log"
	"math"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	repoModel "github.com/Alexey-step/rocket-factory/inventory/internal/repository/model"
)

func (r *repository) InitParts(ctx context.Context) {
	parts := generateParts()

	docs := make([]interface{}, len(parts))
	for i, part := range parts {
		docs[i] = part
	}

	_, err := r.collection.InsertMany(ctx, docs)
	if err != nil {
		panic("failed to insert parts: " + err.Error())
	}
	log.Println("✅ Inventory parts collection initialized")
}

func generateParts() []repoModel.Part {
	names := []string{
		"Main Engine",
		"Reserve Engine",
		"Thruster",
		"Fuel Tank",
		"Left Wing",
		"Right Wing",
		"Window A",
		"Window B",
		"Control Module",
		"Stabilizer",
	}

	descriptions := []string{
		"Primary propulsion unit",
		"Backup propulsion unit",
		"Thruster for fine adjustments",
		"Main fuel tank",
		"Left aerodynamic wing",
		"Right aerodynamic wing",
		"Front viewing window",
		"Side viewing window",
		"Flight control module",
		"Stabilization fin",
	}

	var parts []repoModel.Part
	for i := 0; i < gofakeit.Number(1, 50); i++ {
		idx := gofakeit.Number(0, len(names)-1)
		parts = append(parts, repoModel.Part{
			UUID:          uuid.NewString(),
			Name:          names[idx],
			Description:   descriptions[idx],
			Price:         roundTo(gofakeit.Float64Range(100, 10_000)),
			StockQuantity: int64(gofakeit.Number(1, 100)),
			Category:      repoModel.Category(gofakeit.RandomString([]string{"UNKNOWN", "ENGINE", "FUEL", "PORTHOLE", "WING"})),
			Dimensions:    generateDimensions(),
			Manufacturer:  generateManufacturer(),
			Tags:          generateTags(),
			Metadata:      generateMetadata(),
			CreatedAt:     timestamppb.Now().AsTime(),
		})
	}

	return parts
}

func generateDimensions() repoModel.Dimensions {
	return repoModel.Dimensions{
		Length: roundTo(gofakeit.Float64Range(1, 1000)),
		Width:  roundTo(gofakeit.Float64Range(1, 1000)),
		Height: roundTo(gofakeit.Float64Range(1, 1000)),
		Weight: roundTo(gofakeit.Float64Range(1, 1000)),
	}
}

func generateManufacturer() repoModel.Manufacturer {
	return repoModel.Manufacturer{
		Name:    gofakeit.Name(),
		Country: gofakeit.Country(),
		Website: gofakeit.URL(),
	}
}

func generateTags() []string {
	var tags []string
	for i := 0; i < gofakeit.Number(1, 10); i++ {
		tags = append(tags, gofakeit.EmojiTag())
	}

	return tags
}

func generateMetadata() map[string]repoModel.Metadata {
	return map[string]repoModel.Metadata{
		"string": {
			StringValue: lo.ToPtr(gofakeit.Word()),
		},
		"int": {
			Int64Value: lo.ToPtr(gofakeit.Int64()),
		},
		"double": {
			DoubleValue: lo.ToPtr(gofakeit.Float64()),
		},
		"bool": {
			BoolValue: lo.ToPtr(gofakeit.Bool()),
		},
	}
}

func roundTo(x float64) float64 {
	return math.Round(x*100) / 100
}
