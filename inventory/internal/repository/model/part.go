package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Part struct {
	ID            primitive.ObjectID  `bson:"_id,omitempty"`
	UUID          string              `bson:"uuid"`
	Name          string              `bson:"name"`
	Description   string              `bson:"description"`
	Price         float64             `bson:"price"`
	StockQuantity int64               `bson:"stock_quantity"`
	Category      Category            `bson:"category"`
	Dimensions    Dimensions          `bson:"dimensions"`
	Manufacturer  Manufacturer        `bson:"manufacturer"`
	Tags          []string            `bson:"tags"`
	Metadata      map[string]Metadata `bson:"metadata"`
	CreatedAt     time.Time           `bson:"created_at"`
	UpdatedAt     *time.Time          `bson:"updated_at,omitempty"`
}

type Category string

const (
	CategoryUnspecified Category = "UNKNOWN"  // Неизвестная категория
	CategoryEngine      Category = "ENGINE"   // Двигатель
	CategoryFuel        Category = "FUEL"     // Топливо
	CategoryPorthole    Category = "PORTHOLE" // Иллюминатор
	CategoryWing        Category = "WING"     // Крыло
)

// Dimensions - размеры деталей
type Dimensions struct {
	Length float64 // Длина в см
	Width  float64 // Ширина в см
	Height float64 // Высота в см
	Weight float64 // Вес в кг
}

type Manufacturer struct {
	Name    string // Название
	Country string // Страна производства
	Website string // Сайт производителя
}

type Metadata struct {
	StringValue *string
	Int64Value  *int64
	DoubleValue *float64
	BoolValue   *bool
}

type PartsFilter struct {
	Uuids                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}
