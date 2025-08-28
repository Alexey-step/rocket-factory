package model

type OrderPaid struct {
	EventUUID       string //	Уникальный идентификатор события (для идемпотентности)
	OrderUUID       string //	Идентификатор оплаченного заказа
	UserUUID        string //	Идентификатор пользователя
	PaymentMethod   string //	Способ оплаты (строкой, значение из PaymentMethod)
	TransactionUUID string //	Идентификатор транзакции, сгенерированный в результате оплаты
}

type ShipAssembled struct {
	EventUUID    string // Уникальный идентификатор события (для идемпотентности)
	OrderUUID    string // Идентификатор собранного заказа
	UserUUID     string // Идентификатор пользователя
	BuildTimeSec int64  // Время (в секундах), потраченное на сборку корабля
}
