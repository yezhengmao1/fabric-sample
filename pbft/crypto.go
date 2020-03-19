package pbft

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

func Hash(content []byte) string {
	h := sha256.New()
	h.Write(content)
	return hex.EncodeToString(content)
}

// 签名
func Digest(obj interface{}) (string, error) {
	msg, err := json.Marshal(obj)
	if err != nil {
		logger.Info("error to json marshal object")
		return "", err
	}
	return Hash(msg), nil
}
