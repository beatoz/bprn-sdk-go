package mock

import (
	"github.com/beatoz/bprn-sdk-go/chaincodes"
	"github.com/beatoz/bprn-sdk-go/types"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
)

type InvokeChaincodeFunc func(ctx contractapi.TransactionContextInterface, chaincodeName string, methodName string, methodArgs []string) *peer.Response

type MockChaincodeContext struct {
	callerChaincodeName string
	chainId             *types.ChainId
	invokeChaincodeFn   InvokeChaincodeFunc
	chaincodes.DefaultChaincodeContext
}

func NewMockChaincodeContext(callerChaincodeName string, chainId *types.ChainId) *MockChaincodeContext {
	return &MockChaincodeContext{callerChaincodeName: callerChaincodeName, chainId: chainId}
}

func (mc *MockChaincodeContext) GetChainId(ctx contractapi.TransactionContextInterface) (*types.ChainId, error) {
	return mc.chainId, nil
}

func (mc *MockChaincodeContext) CallerChaincodeName(ctx contractapi.TransactionContextInterface) (string, error) {
	return mc.callerChaincodeName, nil
}

func (mc *MockChaincodeContext) ChangeCallerChaincodeName(callerChaincodeName string) {
	mc.callerChaincodeName = callerChaincodeName
}

func (mc *MockChaincodeContext) ChangeInvokeChaincode(fn InvokeChaincodeFunc) {
	mc.invokeChaincodeFn = fn
}

func (mc *MockChaincodeContext) InvokeChaincode(ctx contractapi.TransactionContextInterface, chaincodeName string, methodName string, methodArgs []string) *peer.Response {
	return mc.invokeChaincodeFn(ctx, chaincodeName, methodName, methodArgs)
}
