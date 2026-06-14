package pay233

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

const (
	headerSignature = "X-Pay233-Signature"
	headerTimestamp = "X-Pay233-Timestamp"
)

func Sign(secret string, timestamp string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp))
	mac.Write([]byte("."))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}
