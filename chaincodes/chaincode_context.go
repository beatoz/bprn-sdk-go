package chaincodes

import (
	"github.com/beatoz/bprn-sdk-go/types"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
)

type ChaincodeContextInterface interface {
	ChannelName(ctx contractapi.TransactionContextInterface) string
	GetChainId(ctx contractapi.TransactionContextInterface) (*types.ChainId, error)
	InvokeChaincode(ctx contractapi.TransactionContextInterface, chaincodeName string, methodName string, methodArgs []string) *peer.Response
	CallerChaincodeName(ctx contractapi.TransactionContextInterface) (string, error)
	GetSignerAddress(ctx contractapi.TransactionContextInterface, selfNamedCc *SelfNamedChaincode, sig string, methodName string, methodParams []string) (*types.Address, error)
	SignerAddress(txid string, sig string, chaincodeName string, methodName string, methodParams []string) (*types.Address, error)
	SetEvent(ctx contractapi.TransactionContextInterface, event interface{}) error
	IsSameChainId(ctx contractapi.TransactionContextInterface, targetChainId string) error
}
