package types

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Address struct {
	address string
}

func NewAddress(hexAddr string) (*Address, error) {
	normalized, err := validateAddress(hexAddr)
	if err != nil {
		return nil, err
	}
	return &Address{address: strings.ToLower(normalized)}, nil
}

func NewChaincodeAddress(channelName string, chaincodeName string) (*Address, error) {
	hashed := sha256.Sum256([]byte(channelName + "-" + chaincodeName))
	address := hex.EncodeToString(hashed[len(hashed)-20:])
	return NewAddress(address)
}

func validateAddress(hexAddress string) (string, error) {
	if hexAddress == "" {
		return "", errors.New("Address cannot be empty")
	}

	lower := strings.ToLower(hexAddress)

	// 0x + 40 hex
	if matched, _ := regexp.MatchString(`^0x[a-f0-9]{40}$`, lower); matched {
		return lower[2:], nil // strip "0x"
	}

	// 40 hex (no prefix)
	if matched, _ := regexp.MatchString(`^[a-f0-9]{40}$`, lower); matched {
		return lower, nil
	}

	return "", fmt.Errorf("invalid address format: %s", hexAddress)
}

func (a *Address) String() string {
	return a.ToHexString()
}

func (a *Address) ToHexString() string {
	return a.address
}

func (a *Address) To0xHexString() string {
	return "0x" + a.address
}

func (a *Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.To0xHexString()) // 또는 a.ToHexString()
}

func (a *Address) UnmarshalJSON(data []byte) error {
	var hexStr string
	if err := json.Unmarshal(data, &hexStr); err != nil {
		return err
	}
	addr, err := NewAddress(hexStr)
	if err != nil {
		return err
	}
	*a = *addr
	return nil
}

func (a *Address) Equal(other *Address) bool {
	return a.address == other.address
}
