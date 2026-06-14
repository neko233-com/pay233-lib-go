package pay233

import "time"

type Money struct {
	Currency string `json:"currency"`
	Amount   int64  `json:"amount"`
}

type EnvType string

const (
	EnvTypeTest    EnvType = "test"
	EnvTypeRelease EnvType = "release"
)

type PaymentStatus string

const (
	StatusCreated  PaymentStatus = "created"
	StatusPending  PaymentStatus = "pending"
	StatusPaid     PaymentStatus = "paid"
	StatusClosed   PaymentStatus = "closed"
	StatusFailed   PaymentStatus = "failed"
	StatusRefunded PaymentStatus = "refunded"
	StatusLost     PaymentStatus = "lost"
)

type CallbackStatus string

const (
	CallbackNone    CallbackStatus = "none"
	CallbackPending CallbackStatus = "pending"
	CallbackSuccess CallbackStatus = "success"
	CallbackFailed  CallbackStatus = "failed"
	CallbackLost    CallbackStatus = "lost"
)

type CreatePaymentRequest struct {
	EnvType    EnvType           `json:"envType"`
	MerchantID string            `json:"merchant_id"`
	OutTradeNo string            `json:"out_trade_no"`
	Channel    string            `json:"channel"`
	Amount     Money             `json:"amount"`
	Subject    string            `json:"subject"`
	NotifyURL  string            `json:"notify_url,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

type Payment struct {
	ID               string            `json:"id"`
	EnvType          EnvType           `json:"env_type"`
	MerchantID       string            `json:"merchant_id"`
	OutTradeNo       string            `json:"out_trade_no"`
	Channel          string            `json:"channel"`
	Amount           Money             `json:"amount"`
	Subject          string            `json:"subject"`
	NotifyURL        string            `json:"notify_url,omitempty"`
	Status           PaymentStatus     `json:"status"`
	ProviderTrade    string            `json:"provider_trade,omitempty"`
	PayURL           string            `json:"pay_url,omitempty"`
	CallbackStatus   CallbackStatus    `json:"callback_status"`
	CallbackAttempts int               `json:"callback_attempts"`
	LastCallbackAt   *time.Time        `json:"last_callback_at,omitempty"`
	CallbackError    string            `json:"callback_error,omitempty"`
	FailureReason    string            `json:"failure_reason,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

type ChannelsResponse struct {
	Channels []string `json:"channels"`
}
