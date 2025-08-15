package telegram

import (
	"bytes"
	"context"
	"embed"
	"html/template"
	"time"

	"go.uber.org/zap"

	"github.com/Alexey-step/rocket-factory/notification/internal/client/http"
	"github.com/Alexey-step/rocket-factory/notification/internal/model"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
)

const chatID = 1141481244

//go:embed templates/paid_notification.tmpl
var paidNotificationTemplateFS embed.FS

//go:embed templates/assembled_notification.tmpl
var assembledNotificationTemplateFS embed.FS

type PaidNotificationTemplateData struct {
	OrderUUID       string
	EventUUID       string
	PaymentMethod   string
	TransactionUUID string
	RegisteredAt    time.Time
}

type AssembledNotificationTemplateData struct {
	OrderUUID    string
	EventUUID    string
	BuildTimeSec int64
	RegisteredAt time.Time
}

var (
	paidNotificationTemplate      = template.Must(template.ParseFS(paidNotificationTemplateFS, "templates/paid_notification.tmpl"))
	assembledNotificationTemplate = template.Must(template.ParseFS(assembledNotificationTemplateFS, "templates/assembled_notification.tmpl"))
)

type service struct {
	telegramClient http.TelegramClient
}

func NewService(telegramClient http.TelegramClient) *service {
	return &service{
		telegramClient: telegramClient,
	}
}

func (s *service) SendPaidNotification(ctx context.Context, paid model.OrderPaid) error {
	message, err := buildPaidNotificationMessage(paid)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat",
		zap.Int("chat_id", chatID),
		zap.String("message", message),
	)

	return nil
}

func (s *service) SendAssembledNotification(ctx context.Context, assembled model.ShipAssembled) error {
	message, err := buildAssembledNotificationMessage(assembled)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}
	logger.Info(ctx, "Telegram message sent to chat",
		zap.Int("chat_id", chatID),
		zap.String("message", message),
	)

	return nil
}

func buildPaidNotificationMessage(paid model.OrderPaid) (string, error) {
	data := &PaidNotificationTemplateData{
		OrderUUID:       paid.OrderUUID,
		EventUUID:       paid.EventUUID,
		PaymentMethod:   paid.PaymentMethod,
		TransactionUUID: paid.TransactionUUID,
		RegisteredAt:    time.Now(),
	}

	var buf bytes.Buffer
	err := paidNotificationTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func buildAssembledNotificationMessage(assembled model.ShipAssembled) (string, error) {
	data := &AssembledNotificationTemplateData{
		OrderUUID:    assembled.OrderUUID,
		EventUUID:    assembled.EventUUID,
		BuildTimeSec: assembled.BuildTimeSec,
		RegisteredAt: time.Now(),
	}

	var buf bytes.Buffer
	err := assembledNotificationTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
