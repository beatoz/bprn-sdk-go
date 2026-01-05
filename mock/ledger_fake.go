package mock

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
)

type LedgerFake struct {
	ledger  map[string]*Stack
	ledger2 map[string][]byte
}

func NewLedgerFake() *LedgerFake {
	return &LedgerFake{
		ledger:  make(map[string]*Stack),
		ledger2: make(map[string][]byte),
	}
}

func (l *LedgerFake) PutState(key string, stateJSON []byte) error {
	l.ledger2[key] = stateJSON
	return nil
}

func (l *LedgerFake) PutState2(key string, stateJSON []byte) error {
	stack, ok := l.ledger[key]
	if !ok {
		stack = &Stack{
			items: make([]interface{}, 0),
		}
		l.ledger[key] = stack
	}
	stack.Push(stateJSON)
	return nil
}

func (l *LedgerFake) GetState(key string) ([]byte, error) {
	state, ok := l.ledger2[key]
	if !ok {
		return nil, nil
	}

	return state, nil
}

func (l *LedgerFake) DeleteState(key string) error {
	delete(l.ledger2, key)
	return nil
}

func (l *LedgerFake) GetState2(key string) ([]byte, error) {
	stack, ok := l.ledger[key]
	if !ok {
		return nil, nil
	}

	stateJSON := stack.Pop()
	if stateJSON == nil {
		return nil, nil
	}
	retval, ok := stateJSON.([]byte)
	if !ok {
		return nil, errors.New("type assertion to []byte failed")
	}
	return retval, nil
}

func (l *LedgerFake) CreateCompositeKey(delimiter string, args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("empty string array")
	}

	compositeKey := strings.Join(args, delimiter)
	return compositeKey, nil
}

func (l *LedgerFake) getTxID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)

	hash := hex.EncodeToString(bytes)
	return hash
}
