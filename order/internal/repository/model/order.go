package model

import "time"

type OrderData struct {
	UUID            string         `json:"uuid"`
	UserUUID        string         `json:"user_uuid"`
	PartUuids       []string       `json:"part_uuids"`
	TotalPrice      float64        `json:"total_price"`
	TransactionUUID *string        `json:"transaction_uuid,omitempty"`
	PaymentMethod   *PaymentMethod `json:"payment_method,omitempty"`
	Status          OrderStatus    `json:"status"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       *time.Time     `json:"updated_at,omitempty"`
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
	OrderStatusCompleted      OrderStatus = "COMPLETED"
)

type OrderCreationInfo struct {
	OrderUUID  string
	TotalPrice float64
}
