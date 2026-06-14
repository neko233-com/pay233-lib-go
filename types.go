package pay233

import "time"

type Money struct {
	Currency string `json:"currency"`
	Amount   int64  `json:"amount"`
}

type PaymentStatus string

const (
	StatusCreated PaymentStatus = "created"
	StatusPending PaymentStatus = "pending"
	StatusPaid    PaymentStatus = "paid"
	StatusClosed  PaymentStatus = "closed"
	StatusFailed  PaymentStatus = "failed"
)

type CreatePaymentRequest struct {
	MerchantID string            `json:"merchant_id"`
	OutTradeNo string            `json:"out_trade_no"`
	Channel    string            `json:"channel"`
	Amount     Money             `json:"amount"`
	Subject    string            `json:"subject"`
	NotifyURL  string            `json:"notify_url,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

type Payment struct {
	ID            string            `json:"id"`
	MerchantID    string            `json:"merchant_id"`
	OutTradeNo    string            `json:"out_trade_no"`
	Channel       string            `json:"channel"`
	Amount        Money             `json:"amount"`
	Subject       string            `json:"subject"`
	NotifyURL     string            `json:"notify_url,omitempty"`
	Status        PaymentStatus     `json:"status"`
	ProviderTrade string            `json:"provider_trade,omitempty"`
	PayURL        string            `json:"pay_url,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type ChannelsResponse struct {
	Channels []string `json:"channels"`
}
