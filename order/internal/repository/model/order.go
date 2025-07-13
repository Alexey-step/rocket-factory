package model

import "time"

type OrderData struct {
	UUID            string
	UserUUID        string
	PartUuids       []string
	TotalPrice      float64
	TransactionUUID *string
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}

type OrderUpdateInfo struct {
	TotalPrice      *float64
	TransactionUUID *string
	PaymentMethod   *PaymentMethod
	Status          *OrderStatus
}

type PaymentMethod string

const (
	PaymentMethodUnknown       PaymentMethod = "UNKNOWN"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCanceled       OrderStatus = "CANCELED"
)

type OrderCreationInfo struct {
	OrderUUID  string
	TotalPrice float64
}
