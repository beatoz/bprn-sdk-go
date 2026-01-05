package types

import (
	"encoding/json"

	"github.com/holiman/uint256"
)

type Payload struct {
	TxHash  string      `json:"txHash"`
	Details interface{} `json:"details"`
	// Data    []byte      `json:"data"`
}

type TrxMetadata struct {
	TxHash   string       `json:"txHash,omitempty"`
	Version  uint32       `json:"version,omitempty"`
	Time     int64        `json:"time"`
	Nonce    uint64       `json:"nonce"`
	From     string       `json:"from"`
	Gas      uint64       `json:"gas"`
	GasPrice *uint256.Int `json:"gasPrice"`
	Type     string       `json:"types"`
	Payload  Payload      `json:"payload,omitempty"`
}
type InvokeResponse struct {
	TX      TrxMetadata       `json:"tx"`
	Channel string            `json:"channel"`
	Event   []json.RawMessage `json:"event"`
	//Event   []Event     `json:"event"`
}
type TransferPayload struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
}
