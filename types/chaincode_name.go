package types

import "github.com/hyperledger/fabric-contract-api-go/v2/contractapi"

type ChaincodeName struct {
	chaincodeName    string
	chaincodeAddress *Address
}

func (cn *ChaincodeName) ChaincodeName() string {
	return cn.chaincodeName
}

func (cn *ChaincodeName) ChaincodeAddress() *Address {
	return cn.chaincodeAddress
}

//func (cn *ChaincodeName) ChaincodeAddress2(ctx contractapi.TransactionContextInterface) (*Address, error) {
//	chaincodeAddress, err := NewChaincodeAddress(ctx.GetStub().GetChannelID(), cn.chaincodeName)
//	if err != nil {
//		return nil, err
//	}
//
//	return chaincodeAddress, nil
//}

func NewChaincodeName(ctx contractapi.TransactionContextInterface, chaincodeName string) (*ChaincodeName, error) {
	chaincodeAddress, err := NewChaincodeAddress(ctx.GetStub().GetChannelID(), chaincodeName)
	if err != nil {
		return nil, err
	}

	return &ChaincodeName{chaincodeName: chaincodeName, chaincodeAddress: chaincodeAddress}, nil
}
