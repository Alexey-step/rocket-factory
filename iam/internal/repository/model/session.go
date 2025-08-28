package model

import "time"

type Session struct {
	UserUUID           string               `bson:"user_uuid"`
	Login              string               `bson:"login"`
	Email              string               `bson:"email"`
	NotificationMethod []NotificationMethod `bson:"notification_method"`
	UUID               string               `bson:"uuid"`
	CreatedAt          time.Time            `bson:"created_at"`
	UpdatedAt          *time.Time           `bson:"updated_at,omitempty"`
	ExpiresAt          time.Time            `bson:"expires_at"`
}

type SessionRedisView struct {
	UserUUID           string `redis:"user_uuid"`
	Login              string `redis:"login"`
	Email              string `redis:"email"`
	NotificationMethod string `redis:"notification_method"`
	UserCreatedAt      int64  `redis:"user_created_at"`
	UserUpdatedAt      *int64 `redis:"user_updated_at,omitempty"`
	UUID               string `redis:"uuid"`
	CreatedAtNs        int64  `redis:"created_at"`
	UpdatedAtNs        *int64 `redis:"updated_at,omitempty"`
	ExpiresAt          int64  `redis:"expires_at"`
}
