package chaincodes

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/beatoz/bprn-sdk-go/fabric"
	"github.com/beatoz/bprn-sdk-go/types"
	"github.com/holiman/uint256"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
)

type DefaultChaincodeContext struct {
}

func (bc *DefaultChaincodeContext) ChannelName(ctx contractapi.TransactionContextInterface) string {
	return ctx.GetStub().GetChannelID()
}

func (bc *DefaultChaincodeContext) GetChainId(ctx contractapi.TransactionContextInterface) (uint256.Int, error) {
	return fabric.GetChainId(ctx.GetStub())
}

func (bc *DefaultChaincodeContext) InvokeChaincode(ctx contractapi.TransactionContextInterface, chaincodeName string, methodName string, methodArgs []string) *peer.Response {
	stub := ctx.GetStub()

	args := [][]byte{[]byte(methodName)}
	for _, methodArg := range methodArgs {
		args = append(args, []byte(methodArg))
	}

	return stub.InvokeChaincode(chaincodeName, args, stub.GetChannelID())
}

func (bc *DefaultChaincodeContext) CallerChaincodeName(ctx contractapi.TransactionContextInterface) (string, error) {
	callerChaincodeName, err := fabric.NewFabricUtil(ctx.GetStub()).CallerChaincodeName()
	if err != nil {
		return "", fmt.Errorf("failed to get caller chaincode name: %w", err)
	}

	return callerChaincodeName, nil
}

func (bc *DefaultChaincodeContext) GetSignerAddress(ctx contractapi.TransactionContextInterface, selfNamedCc *SelfNamedChaincode, sig string, methodName string, methodParams []string) (*types.Address, error) {
	chaincodeName, err := selfNamedCc.SelfChaincodeName(ctx)
	if err != nil {
		return nil, err
	}

	return bc.SignerAddress(sig, chaincodeName, methodName, methodParams)
}

func (bc *DefaultChaincodeContext) SignerAddress(sig string, chaincodeName string, methodName string, methodParams []string) (*types.Address, error) {
	return fabric.SigVerifyAndSignerAddress(sig, chaincodeName, methodName, methodParams)
}

func (bc *DefaultChaincodeContext) SetEvent(ctx contractapi.TransactionContextInterface, event interface{}) error {
	eventType := reflect.TypeOf(event)
	eventName := eventType.Name()

	// 포인터인 경우 실제 타입 이름 가져오기
	if eventType.Kind() == reflect.Ptr {
		eventName = eventType.Elem().Name()
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to obtain JSON encoding: %v", err)
	}

	err = ctx.GetStub().SetEvent(eventName, eventJSON)
	if err != nil {
		return fmt.Errorf("failed to set event: %v", err)
	}

	return nil
}

func (bc *DefaultChaincodeContext) IsSameChainId(ctx contractapi.TransactionContextInterface, targetChainId string) error {
	dstChainId, err := uint256.FromDecimal(targetChainId)
	if err != nil {
		return fmt.Errorf("failed to convert fromChainId to uint64: %w", err)
	}
	return fabric.IsSameChainId(ctx, dstChainId)
}
