package generator

import (
	"github.com/beatoz/bprn-sdk-go/types"
)

type ERC20ArgsGenerator struct {
}

type BeatozErc20Param struct {
	ChaincodeMethodName string                    `json:"chaincodeMethodName"`
	EvmTransactionParam types.EvmTransactionParam `json:"evmTransactionParam"`
	ChaincodeParam      []string                  `json:"chaincodeParam"`
	Signature           string                    `json:"signature"`
}

func NewBeatozErc20Param(chaincodeMethodNam string, evmTransactionParam *types.EvmTransactionParam, chaincodeParam []string, signature string) *BeatozErc20Param {
	return &BeatozErc20Param{
		ChaincodeMethodName: chaincodeMethodNam,
		EvmTransactionParam: *evmTransactionParam,
		ChaincodeParam:      chaincodeParam,
		Signature:           signature,
	}
}

func NewBeatozErc20ParamFromArray(beatozErc20Params []string) *BeatozErc20Param {
	beatozErc20Param := &BeatozErc20Param{}
	beatozErc20Param.from(beatozErc20Params)
	return beatozErc20Param
}

func (e *BeatozErc20Param) toArray() []string {
	//result := []string{e.ChaincodeMethodName}
	var result []string
	//result = append(result, e.EvmTransactionParam.ToArray()...)
	result = append(result, e.ChaincodeParam...)
	return append(result, e.Signature)
}

func (e *BeatozErc20Param) from(beatozErc20Params []string) []string {
	if len(beatozErc20Params) < 1 {
		return []string{}
	}

	e.ChaincodeMethodName = beatozErc20Params[0]
	params := e.EvmTransactionParam.From(beatozErc20Params[1:])
	e.ChaincodeParam = params[:len(params)-1]
	e.Signature = beatozErc20Params[len(beatozErc20Params)-1]

	return e.ChaincodeParam
}

func (e *BeatozErc20Param) isEqual(dest *BeatozErc20Param) bool {
	if e.ChaincodeMethodName != dest.ChaincodeMethodName {
		return false
	}

	if !e.EvmTransactionParam.IsEqual(dest.EvmTransactionParam) {
		return false
	}

	if e.Signature != dest.Signature {
		return false
	}

	if len(e.ChaincodeParam) != len(dest.ChaincodeParam) {
		return false
	}

	for i := range e.ChaincodeParam {
		if e.ChaincodeParam[i] != dest.ChaincodeParam[i] {
			return false
		}
	}

	return true
}
