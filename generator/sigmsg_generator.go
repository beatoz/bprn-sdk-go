package generator

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
)

type SigMsgGenerator struct {
}

func NewSigMsgGenerator() *SigMsgGenerator {
	return &SigMsgGenerator{}
}

func (s *SigMsgGenerator) GenerateSigMsg(chaincodeName string, chaincodeFunctionName string, chaincodeParams []string) ([]byte, error) {
	sigMsgArray := s.toArray(chaincodeName, chaincodeFunctionName, chaincodeParams)
	sigMsg, err := s.encode(sigMsgArray)
	if err != nil {
		return nil, err
	}
	return sigMsg, nil
}

func (s *SigMsgGenerator) toArray(chaincodeName string, chaincodeFunctionName string, chaincodeParams []string) []interface{} {
	var sigMsgArr []interface{}
	sigMsgArr = append(sigMsgArr, []byte(chaincodeName))
	sigMsgArr = append(sigMsgArr, []byte(chaincodeFunctionName))
	for _, param := range chaincodeParams {
		sigMsgArr = append(sigMsgArr, []byte(param))
	}
	return sigMsgArr
}

func (s *SigMsgGenerator) encode(sigMsgArr []interface{}) ([]byte, error) {
	rlpSigMsg, err := rlp.EncodeToBytes(sigMsgArr)
	if err != nil {
		return nil, fmt.Errorf("failed to RLP encode toArray: %v", err)
	}
	return rlpSigMsg, nil
}
