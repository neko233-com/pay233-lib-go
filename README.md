# pay233-lib-go

Go client SDK for `pay233-server`.

```go
client := pay233.NewClient("http://localhost:5500", "dev-secret")
payment, err := client.CreatePayment(ctx, pay233.CreatePaymentRequest{
    EnvType: pay233.EnvTypeTest,
    MerchantID: "merchant_1",
    OutTradeNo: "order_10001",
    Channel: "mock",
    Amount: pay233.Money{Currency: "CNY", Amount: 100},
    Subject: "Test order",
})
```

API requests are signed with `X-Pay233-Timestamp` and `X-Pay233-Signature` when a signing secret is configured.

`EnvType` controls environment isolation. Use `pay233.EnvTypeTest` for sandbox traffic and `pay233.EnvTypeRelease` for formal payment traffic. The server defaults missing values to `test`, but production callers should set it explicitly.

## Test

```bash
go test -cover ./...
go vet ./...
```
