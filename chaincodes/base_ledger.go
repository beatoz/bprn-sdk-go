package chaincodes

import (
	"errors"
	"strconv"

	"github.com/holiman/uint256"
	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

var ErrNotFound = errors.New("not found")

type BaseLedger struct {
	stub shim.ChaincodeStubInterface
}

func NewBaseLedger(ctx contractapi.TransactionContextInterface) *BaseLedger {
	return &BaseLedger{
		stub: ctx.GetStub(),
	}
}

func (k *BaseLedger) Delete(key string) error {
	if err := k.stub.DelState(key); err != nil {
		return err
	}
	return nil
}

func (k *BaseLedger) PutUint256(key string, value *uint256.Int) error {
	return k.PutString(key, value.String())
}

func (k *BaseLedger) PutUint64(key string, value uint64) error {
	return k.PutString(key, strconv.FormatUint(value, 10))
}

func (k *BaseLedger) PutString(key string, value string) error {
	if err := k.stub.PutState(key, []byte(value)); err != nil {
		return err
	}
	return nil
}

func (k *BaseLedger) PutUint8(key string, value uint8) error {
	return k.stub.PutState(key, []byte{value})
}

func (k *BaseLedger) PutBytes(key string, value []byte) error {
	if err := k.stub.PutState(key, value); err != nil {
		return err
	}
	return nil
}

func (k *BaseLedger) Get(key string) ([]byte, error) {
	value, err := k.stub.GetState(key)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, ErrNotFound
	}
	return value, nil
}

func (k *BaseLedger) GetUint8(key string) (uint8, error) {
	value, err := k.Get(key)
	if err != nil {
		return 0, err
	}
	if value == nil {
		return 0, nil
	}
	return value[0], nil
}

func (k *BaseLedger) GetString(key string) (string, error) {
	value, err := k.Get(key)
	if err != nil {
		return "", err
	}
	if value == nil {
		return "", nil
	}
	return string(value), nil
}

func (k *BaseLedger) IsExist(key string) (bool, error) {
	value, err := k.Get(key)
	if err != nil {
		return false, err
	}
	return value != nil, nil
}

func (k *BaseLedger) GetUint64(key string) (uint64, error) {
	valueBytes, err := k.Get(key)
	if err != nil {
		return 0, err
	}
	if valueBytes == nil {
		return 0, ErrNotFound
	}

	return strconv.ParseUint(string(valueBytes), 10, 64)
}

func (k *BaseLedger) GetUint256(key string) (*uint256.Int, error) {
	valueBytes, err := k.Get(key)
	if err != nil {
		return nil, err
	}
	if valueBytes == nil {
		return nil, ErrNotFound
	}

	value, err := uint256.FromDecimal(string(valueBytes))
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (k *BaseLedger) CreateCompositeKey(objectType string, attributes []string) (string, error) {
	return k.stub.CreateCompositeKey(objectType, attributes)
}
