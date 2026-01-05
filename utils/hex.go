package utils

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/holiman/uint256"
)

func NormalizeHexString(hexStr string) string {
	hexStr = strings.ToLower(hexStr)
	hexStr = strings.TrimPrefix(hexStr, "0x")

	// 홀수 길이면 앞에 "0" 붙이기
	if len(hexStr)%2 == 1 {
		hexStr = "0" + hexStr
	}

	return hexStr
}

func NormalizeHexBytesFromUint256(v uint256.Int) ([]byte, error) {
	hexStr := v.Hex() // [주의] 0x 포함
	normalized := NormalizeHexString(hexStr)
	hexBytes, err := hex.DecodeString(normalized)
	if err != nil {
		return nil, err
	}

	return hexBytes, nil
}

func NormalizeHexBytesFromUint64(v uint64) []byte {
	hexStr := fmt.Sprintf("%x", v)
	normalized := NormalizeHexString(hexStr)
	b, _ := hex.DecodeString(normalized)
	return b
}
