package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty"`
	UUID      string             `json:"uuid"`
	Info      UserInfo           `json:"info"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty"`
	Password  []byte             `json:"password_hash"`
}

type UserInfo struct {
	Login               string               `json:"login"`               // Логин
	Email               string               `json:"email"`               // Email
	NotificationMethods []NotificationMethod `json:"notification_method"` //	Каналы уведомлений
}

type NotificationMethod struct {
	ProviderName string `json:"provider_name"` // Провайдер: telegram, email, push и т.д.
	Target       string `json:"target"`        // Адрес/идентификатор назначения (email, чат-id)
}
