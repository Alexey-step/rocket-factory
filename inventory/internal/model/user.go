package model

import "time"

type User struct {
	UUID      string     // UUID пользователя
	Info      UserInfo   // Базовая информация
	CreatedAt time.Time  // Дата создания
	UpdatedAt *time.Time // Дата обновления
}

type UserInfo struct {
	Login               string               // Логин
	Email               string               // Email
	NotificationMethods []NotificationMethod //	Каналы уведомлений
}

type NotificationMethod struct {
	ProviderName string // Провайдер: telegram, email, push и т.д.
	Target       string // Адрес/идентификатор назначения (email, чат-id)
}
