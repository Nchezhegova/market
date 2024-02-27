package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"github.com/Nchezhegova/market/internal/config"
)

func CalculateHash(pass string) []byte {
	hmacKey := []byte(config.PASSWORDHASH)
	h := hmac.New(sha256.New, hmacKey)
	h.Write([]byte(pass))
	return h.Sum(nil)
}
