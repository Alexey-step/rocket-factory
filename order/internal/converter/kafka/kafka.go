package kafka

import "github.com/Alexey-step/rocket-factory/order/internal/model"

type OrderAssembledDecoder interface {
	Decode(data []byte) (model.ShipAssembled, error)
}
