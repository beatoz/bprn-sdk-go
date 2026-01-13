package types

import (
	"strings"

	"github.com/holiman/uint256"
)

type ChainId struct {
	chainId *uint256.Int
}

func NewChainID(hexOrDecimalChainId string) (*ChainId, error) {
	if strings.HasPrefix(hexOrDecimalChainId, "0x") || strings.HasPrefix(hexOrDecimalChainId, "0X") {
		return NewChainIDFromHex(hexOrDecimalChainId)
	}
	return NewChainIDFromDecimal(hexOrDecimalChainId)
}

func NewChainIDFromDecimal(decimal string) (*ChainId, error) {
	chainIdU256, err := uint256.FromDecimal(decimal)
	if err != nil {
		return nil, err
	}
	return &ChainId{chainId: chainIdU256}, nil
}

func NewChainIDFromHex(hex string) (*ChainId, error) {
	chainIdU256, err := uint256.FromHex(hex)
	if err != nil {
		return nil, err
	}
	return &ChainId{chainId: chainIdU256}, nil
}

func (c *ChainId) Uint256() *uint256.Int {
	return c.chainId
}

func (c *ChainId) Dec() string {
	return c.chainId.Dec()
}

func (c *ChainId) Hex() string {
	return c.chainId.Hex()
}

func (c *ChainId) String() string {
	return c.Dec()
}

func (c *ChainId) Equal(other *ChainId) bool {
	return c.chainId.Eq(other.chainId)
}

func (c *ChainId) MarshalJSON() ([]byte, error) {
	if c.chainId == nil {
		return []byte("null"), nil
	}
	return c.chainId.MarshalJSON()
}

func (c *ChainId) UnmarshalJSON(data []byte) error {
	if c.chainId == nil {
		c.chainId = new(uint256.Int)
	}
	return c.chainId.UnmarshalJSON(data)
}
