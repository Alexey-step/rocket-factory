package part

import (
	"context"
	"slices"

	"github.com/Alexey-step/rocket-factory/inventory/internal/model"
	repoConverter "github.com/Alexey-step/rocket-factory/inventory/internal/repository/converter"
)

func (r *repository) ListParts(_ context.Context, filter model.PartsFilter) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.Part
	for _, part := range r.data {
		modelPart := repoConverter.PartToModel(part)
		if matchesFilter(modelPart, filter) {
			result = append(result, modelPart)
		}
	}

	if len(result) == 0 {
		return nil, model.ErrPartsNotFound
	}

	return result, nil
}

func matchesFilter(part model.Part, filter model.PartsFilter) bool {
	if len(filter.Uuids) > 0 && len(part.UUID) > 0 && !slices.Contains(filter.Uuids, part.UUID) {
		return false
	}

	// Фильтрация по имени
	if len(filter.Names) > 0 && !slices.Contains(filter.Names, part.Name) {
		return false
	}

	// Фильтрация по категории
	if len(filter.Categories) > 0 && !slices.Contains(filter.Categories, part.Category) {
		return false
	}

	// Фильтрация по странам
	if len(filter.ManufacturerCountries) > 0 && !slices.Contains(filter.ManufacturerCountries, part.Manufacturer.Country) {
		return false
	}

	// Фильтрация по тегам (если хотя бы один тег совпадает)
	if len(filter.Tags) > 0 && !hasCommonElement(filter.Tags, part.Tags) {
		return false
	}

	return true
}

func hasCommonElement(a, b []string) bool {
	for _, v := range a {
		if slices.Contains(b, v) {
			return true
		}
	}
	return false
}
