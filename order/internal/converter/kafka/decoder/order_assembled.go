package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
	eventsV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/events/v1"
)

type decoder struct{}

func NewOrderAssembledDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.ShipAssembled, error) {
	var pb eventsV1.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembled{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.ShipAssembled{
		EventUUID:    pb.EventUuid,
		OrderUUID:    pb.OrderUuid,
		UserUUID:     pb.UserUuid,
		BuildTimeSec: pb.BuildTimeSec,
	}, nil
}
