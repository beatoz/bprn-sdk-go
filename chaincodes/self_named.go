package chaincodes

import (
	"fmt"

	"github.com/beatoz/bprn-sdk-go/types"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SelfNamedChaincode struct {
	chaincodeName    string
	chaincodeAddress *types.Address
	chaincodeContext ChaincodeContextInterface
}

func NewSelfNamedChaincode(chaincodeContext ChaincodeContextInterface) *SelfNamedChaincode {
	return &SelfNamedChaincode{
		chaincodeContext: chaincodeContext,
	}
}

// TODO : 다른 체인코드에서 InitLedger를 호출했을때 결과 확인
func (sn *SelfNamedChaincode) InitSelf(ctx contractapi.TransactionContextInterface) error {
	ledger := NewSelfNamedChaincodeLedger(ctx)

	if selfChaincode, err := ledger.GetSelfChaincodeName(); err == nil {
		return fmt.Errorf("self chaincode name already exists: %s", selfChaincode)
	}

	selfChaincodeName, err := sn.chaincodeContext.CallerChaincodeName(ctx)
	if err != nil {
		return fmt.Errorf("failed to get caller chaincode name: %w", err)
	}
	if selfChaincodeName == "" {
		return fmt.Errorf("caller chaincode name is empty")
	}

	sn.chaincodeName = selfChaincodeName
	chaincodeAddress, err := types.NewChaincodeAddress(ctx.GetStub().GetChannelID(), selfChaincodeName)
	if err != nil {
		return fmt.Errorf("failed to create chaincode address: %w", err)
	}
	sn.chaincodeAddress = chaincodeAddress

	if err := ledger.PutSelfChaincodeName(selfChaincodeName); err != nil {
		return fmt.Errorf("failed to put self chaincode name: %w", err)
	}

	fmt.Println("[SelfNamedChaincode.InitLedger] selfChaincodeName: ", sn.chaincodeName)

	return nil
}

func (sn *SelfNamedChaincode) ChaincodeAddress(ctx contractapi.TransactionContextInterface) (*types.Address, error) {
	if sn.chaincodeAddress == nil {
		selfChaincodeName, err := sn.SelfChaincodeName(ctx)
		if err != nil {
			return nil, err
		}

		chaincodeAddress, err := types.NewChaincodeAddress(ctx.GetStub().GetChannelID(), selfChaincodeName)
		if err != nil {
			return nil, err
		}
		sn.chaincodeAddress = chaincodeAddress
	}

	return sn.chaincodeAddress, nil
}

func (sn *SelfNamedChaincode) SelfChaincodeName(ctx contractapi.TransactionContextInterface) (string, error) {
	if sn.chaincodeName == "" {
		ledger := NewSelfNamedChaincodeLedger(ctx)
		selfChaincodeName, err := sn.chaincodeNameFromLedger(ledger)
		if err != nil {
			return "", err
		}
		sn.chaincodeName = selfChaincodeName
	}

	return sn.chaincodeName, nil
}

func (sn *SelfNamedChaincode) chaincodeNameFromLedger(ledger *SelfNamedChaincodeLedger) (string, error) {
	selfChaincodeName, err := ledger.GetSelfChaincodeName()
	if err != nil {
		return "", fmt.Errorf("failed to get self chaincode name: %w", err)
	}
	return selfChaincodeName, nil
}
