package kafka

import "github.com/Alexey-step/rocket-factory/assembly/internal/model"

type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaid, error)
}
