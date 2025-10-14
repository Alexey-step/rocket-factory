package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                  primitive.ObjectID `db:"_id,omitempty"`
	UUID                string             `db:"uuid"`
	Login               string             `db:"login"`                // Логин
	Email               string             `db:"email"`                // Email
	NotificationMethods []byte             `db:"notification_methods"` //	Каналы уведомлений
	CreatedAt           time.Time          `db:"created_at"`
	UpdatedAt           *time.Time         `db:"updated_at,omitempty"`
	Password            string             `db:"password_hash"`
}

type NotificationMethod struct {
	ProviderName string `json:"provider_name"` // Провайдер: telegram, email, push и т.д.
	Target       string `json:"target"`        // Адрес/идентификатор назначения (email, чат-id)
}
