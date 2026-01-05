package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/holiman/uint256"
)

func addBigInt(a string, b string) string {
	aInt := new(big.Int)
	aInt.SetString(a, 10)

	bInt := new(big.Int)
	bInt.SetString(b, 10)

	result := new(big.Int).Add(aInt, bInt)
	return result.String()
}

func subBigInt(a string, b string) string {
	aInt := new(big.Int)
	aInt.SetString(a, 10)

	bInt := new(big.Int)
	bInt.SetString(b, 10)

	result := new(big.Int).Sub(aInt, bInt)
	return result.String()
}

func hasSufficientBalance(a string, b string) bool {
	aInt := new(big.Int)
	aInt.SetString(a, 10)

	bInt := new(big.Int)
	bInt.SetString(b, 10)

	return aInt.Cmp(bInt) >= 0
}

func uint256ToBytes32(v *uint256.Int) []byte {
	// uint256.Int.Bytes32()는 이미 Solidity와 동일한 32바이트 big-endian 반환
	b := v.Bytes32()
	return b[:] // [32]byte → []byte
}

func StringToUint256(s string) (*uint256.Int, error) {
	b := []byte(s)
	if len(b) > 32 {
		return nil, fmt.Errorf("string too long: %d bytes (max 32)", len(b))
	}

	buf := make([]byte, 32)
	copy(buf[32-len(b):], b) // 오른쪽 정렬

	n := new(uint256.Int)
	n.SetBytes32(buf) // 반드시 32바이트 필요
	return n, nil
}

func Uint256ToString(n *uint256.Int) string {
	buf := n.Bytes32() // 항상 32바이트

	// 앞쪽(왼쪽)의 0-padding 제거
	i := 0
	for i < 32 && buf[i] == 0 {
		i++
	}

	return string(buf[i:])
}

func Uint256FromHex(hexStr string) (*uint256.Int, error) {
	n := new(uint256.Int)
	err := n.SetFromHex(hexStr)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func Generate16BytesRandom() []byte {
	// Generate 16 random bytes (will become 32 hex characters)
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		// This should rarely happen, but if it does, we can't continue safely
		panic(fmt.Sprintf("failed to generate random ID: %v", err))
	}
	return randomBytes
}

func NewID() (string, error) {
	const size = 16 // 128bit
	b := make([]byte, size)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
