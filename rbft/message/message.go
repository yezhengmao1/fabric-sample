package message

import (
	"crypto/sha256"
	"encoding/hex"
)

type TimeStamp uint64 // 时间戳格式
type Identify uint64  // 客户端标识格式
type View Identify    // 视图
type Sequence int64   // 序号

func Hash(content []byte) string {
	return hex.EncodeToString(HashByte(content))
}

func HashByte(content []byte) []byte {
	h := sha256.New()
	h.Write(content)
	return h.Sum(nil)
}
