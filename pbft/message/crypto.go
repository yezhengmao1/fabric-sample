package message

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
)

func Hash(content []byte) string {
	h := sha256.New()
	h.Write(content)
	return hex.EncodeToString(content)
}

func Digest(obj interface{}) (string, error) {
	content, err := json.Marshal(obj)
	if err != nil {
		log.Printf("[Crypto] marshl the object error: %s", err)
		return "", err
	}
	return Hash(content), nil
}
