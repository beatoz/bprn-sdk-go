package types

import (
	"strconv"

	"github.com/holiman/uint256"
)

type EvmTransactionParam struct {
	Version  uint32      `json:"version"`
	Nonce    uint64      `json:"nonce"`
	Gas      uint64      `json:"gas"`
	GasPrice uint256.Int `json:"gasPrice"`
	//ChaincodeParam []string `json:"args,omitempty"`
}

func NewDefaultEvmTxParam() *EvmTransactionParam {
	return NewEvmTxParam(1, 0, 21000, *uint256.NewInt(100000000))
}

func NewEvmTxParamFromString(version, nonce, gas, gasPrice string) (*EvmTransactionParam, error) {
	versionUint32, err := strconv.ParseUint(version, 10, 32)
	if err != nil {
		return nil, err
	}
	nonceUint64, err := strconv.ParseUint(nonce, 10, 64)
	if err != nil {
		return nil, err
	}
	gasUint64, err := strconv.ParseUint(gas, 10, 64)
	if err != nil {
		return nil, err
	}

	gasPriceUint256, err := uint256.FromHex("0x" + gasPrice)
	if err != nil {
		return nil, err
	}

	return NewEvmTxParam(uint32(versionUint32), nonceUint64, gasUint64, *gasPriceUint256), nil
}

func NewEvmTxParam(version uint32, nonce, gas uint64, gasPrice uint256.Int) *EvmTransactionParam {
	return &EvmTransactionParam{
		Version:  version,
		Nonce:    nonce,
		Gas:      gas,
		GasPrice: gasPrice,
		//ChaincodeParam: chaincodeParam,
	}
}

func (e *EvmTransactionParam) ToArray() []interface{} {
	paramArray := []interface{}{
		e.Version,
		e.Nonce,
		e.Gas,
		e.GasPrice.Bytes(),
	}

	return paramArray

	//hexGasPrice := hex.EncodeToString(e.GasPrice.Bytes())
	//gasPriceHexBytes, err := utils.NormalizeHexBytesFromUint256(e.GasPrice)

	//serialized, err := rlp.EncodeToBytes(rlpInput)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return serialized, nil
}

func (e *EvmTransactionParam) From(params []string) []string {
	if len(params) < 4 {
		return []string{}
	}

	version, err := strconv.ParseUint(params[0], 10, 32)
	if err != nil {
		return []string{}
	}
	e.Version = uint32(version)

	nonce, err := strconv.ParseUint(params[1], 10, 64)
	if err != nil {
		return []string{}
	}
	e.Nonce = nonce

	gas, err := strconv.ParseUint(params[2], 10, 64)
	if err != nil {
		return []string{}
	}
	e.Gas = gas

	gasPrice, err := uint256.FromDecimal(params[3])
	if err != nil {
		return []string{}
	}
	e.GasPrice = *gasPrice

	//e.Version = int([]rune(params[0])[0])
	//e.Nonce = params[1]
	//e.Gas = params[2]
	//e.GasPrice = params[3]
	//e.ChaincodeParam = params[4:]

	return params[4:]
}

func (e *EvmTransactionParam) IsEqual(dest EvmTransactionParam) bool {
	return e.Version == dest.Version &&
		e.Nonce == dest.Nonce &&
		e.Gas == dest.Gas &&
		e.GasPrice == dest.GasPrice
}
