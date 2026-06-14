package pay233

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL       string
	signingSecret string
	httpClient    *http.Client
}

type Option func(*Client)

func NewClient(baseURL string, signingSecret string, options ...Option) *Client {
	client := &Client{
		baseURL:       strings.TrimRight(baseURL, "/"),
		signingSecret: signingSecret,
		httpClient:    http.DefaultClient,
	}
	for _, option := range options {
		option(client)
	}
	return client
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(client *Client) {
		if httpClient != nil {
			client.httpClient = httpClient
		}
	}
}

func (c *Client) Channels(ctx context.Context) ([]string, error) {
	var resp ChannelsResponse
	if err := c.do(ctx, http.MethodGet, "/v1/channels", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Channels, nil
}

func (c *Client) CreatePayment(ctx context.Context, req CreatePaymentRequest) (Payment, error) {
	var payment Payment
	if err := c.do(ctx, http.MethodPost, "/v1/payments", req, &payment); err != nil {
		return Payment{}, err
	}
	return payment, nil
}

func (c *Client) GetPayment(ctx context.Context, id string) (Payment, error) {
	var payment Payment
	if err := c.do(ctx, http.MethodGet, "/v1/payments/"+id, nil, &payment); err != nil {
		return Payment{}, err
	}
	return payment, nil
}

func (c *Client) ClosePayment(ctx context.Context, id string) (Payment, error) {
	var payment Payment
	if err := c.do(ctx, http.MethodPost, "/v1/payments/"+id+"/close", struct{}{}, &payment); err != nil {
		return Payment{}, err
	}
	return payment, nil
}

func (c *Client) do(ctx context.Context, method string, path string, input any, output any) error {
	var body []byte
	var err error
	if input != nil {
		body, err = json.Marshal(input)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	if input != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.signingSecret != "" && method != http.MethodGet {
		timestamp := time.Now().UTC().Format(time.RFC3339)
		req.Header.Set(headerTimestamp, timestamp)
		req.Header.Set(headerSignature, Sign(c.signingSecret, timestamp, body))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pay233: request failed with status %d: %s", resp.StatusCode, string(respBody))
	}
	if output == nil || len(respBody) == 0 {
		return nil
	}
	return json.Unmarshal(respBody, output)
}
