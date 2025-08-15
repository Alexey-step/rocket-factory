package consumer

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/Alexey-step/rocket-factory/platform/pkg/kafka"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type consumer struct {
	logger      Logger
	topics      []string
	group       sarama.ConsumerGroup
	middlewares []Middleware
}

// NewConsumer — создаёт новый consumer.
func NewConsumer(logger Logger, group sarama.ConsumerGroup, topics []string, middlewares ...Middleware) *consumer {
	return &consumer{
		logger:      logger,
		group:       group,
		topics:      topics,
		middlewares: middlewares,
	}
}

// Consume запускает консьюмер для списка топиков.
func (c *consumer) Consume(ctx context.Context, handler kafka.MessageHandler) error {
	newGroupHandler := NewGroupHandler(handler, c.logger, c.middlewares...)

	for {
		if err := c.group.Consume(ctx, c.topics, newGroupHandler); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return nil
			}

			c.logger.Error(ctx, "Kafka consume error", zap.Error(err))
			return err
		}

		if ctx.Err() != nil {
			c.logger.Info(ctx, "Kafka consumer context cancelled")
			return ctx.Err()
		}

		c.logger.Info(ctx, "Kafka consumer group rebalancing...")
	}
}
