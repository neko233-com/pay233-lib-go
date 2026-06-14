package pay233

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientSignsCreatePayment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(headerTimestamp) == "" {
			t.Fatal("missing timestamp")
		}
		if r.Header.Get(headerSignature) == "" {
			t.Fatal("missing signature")
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(body), `"envType":"test"`) {
			t.Fatalf("missing envType in request body: %s", body)
		}
		_ = json.NewEncoder(w).Encode(Payment{ID: "pay_1", Status: StatusPending})
	}))
	defer server.Close()

	client := NewClient(server.URL, "secret")
	got, err := client.CreatePayment(context.Background(), CreatePaymentRequest{
		EnvType:    EnvTypeTest,
		MerchantID: "m1",
		OutTradeNo: "o1",
		Channel:    "mock",
		Amount:     Money{Currency: "CNY", Amount: 100},
		Subject:    "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "pay_1" {
		t.Fatalf("expected pay_1, got %s", got.ID)
	}
}

func TestClientChannels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/channels" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(ChannelsResponse{Channels: []string{"mock"}})
	}))
	defer server.Close()

	client := NewClient(server.URL, "secret")
	channels, err := client.Channels(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(channels) != 1 || channels[0] != "mock" {
		t.Fatalf("unexpected channels: %#v", channels)
	}
}

func TestClientGetAndClosePayment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/payments/pay_1":
			_ = json.NewEncoder(w).Encode(Payment{ID: "pay_1", Status: StatusPending})
		case r.Method == http.MethodPost && r.URL.Path == "/v1/payments/pay_1/close":
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}
			expected := Sign("secret", r.Header.Get(headerTimestamp), body)
			if r.Header.Get(headerSignature) != expected {
				t.Fatal("invalid close signature")
			}
			_ = json.NewEncoder(w).Encode(Payment{ID: "pay_1", Status: StatusClosed})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "secret")
	got, err := client.GetPayment(context.Background(), "pay_1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != StatusPending {
		t.Fatalf("expected pending, got %s", got.Status)
	}

	closed, err := client.ClosePayment(context.Background(), "pay_1")
	if err != nil {
		t.Fatal(err)
	}
	if closed.Status != StatusClosed {
		t.Fatalf("expected closed, got %s", closed.Status)
	}
}

func TestClientReturnsServerErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"boom"}`, http.StatusBadRequest)
	}))
	defer server.Close()

	client := NewClient(server.URL, "secret")
	_, err := client.CreatePayment(context.Background(), CreatePaymentRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "status 400") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClientUsesCustomHTTPClient(t *testing.T) {
	transport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"channels":["mock"]}`)),
			Request:    req,
		}, nil
	})

	client := NewClient("https://pay233.test", "", WithHTTPClient(&http.Client{Transport: transport}))
	channels, err := client.Channels(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprint(channels) != "[mock]" {
		t.Fatalf("unexpected channels: %#v", channels)
	}
}

func TestSign(t *testing.T) {
	got := Sign("secret", "2026-06-14T00:00:00Z", []byte(`{"ok":true}`))
	want := "db17c1b08becb83d277da933d16d8c84ce3811abc2bc1df2595aff504bf34709"
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
