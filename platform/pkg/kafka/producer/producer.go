package producer

import (
	"context"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type producer struct {
	syncProducer sarama.SyncProducer
	topic        string
	logger       Logger
}

func NewProducer(syncProducer sarama.SyncProducer, topic string, logger Logger) *producer {
	return &producer{
		syncProducer: syncProducer,
		topic:        topic,
		logger:       logger,
	}
}

func (p *producer) Send(ctx context.Context, key string, value []byte) error {
	partition, offset, err := p.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	})
	if err != nil {
		p.logger.Error(ctx, "failed to send message", zap.String("topic", p.topic), zap.Error(err))
		return err
	}

	p.logger.Info(ctx, "message sent successfully",
		zap.String("topic", p.topic),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
		zap.String("key", key),
		zap.ByteString("value", value),
	)
	return nil
}
