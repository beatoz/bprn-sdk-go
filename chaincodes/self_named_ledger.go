package chaincodes

import (
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

const (
	selfChaincodeNameKey = "selfChaincodeName"
)

type SelfNamedChaincodeLedger struct {
	Ledger *BaseLedger
}

func NewSelfNamedChaincodeLedger(ctx contractapi.TransactionContextInterface) *SelfNamedChaincodeLedger {
	baseLedger := NewBaseLedger(ctx)
	return NewSelfNamedLedgerFromBaseLedger(baseLedger)
}

func NewSelfNamedLedgerFromBaseLedger(ledger *BaseLedger) *SelfNamedChaincodeLedger {
	return &SelfNamedChaincodeLedger{
		Ledger: ledger,
	}
}

func (sl *SelfNamedChaincodeLedger) GetBaseChaincodeLedger() *SelfNamedChaincodeLedger {
	return sl
}

func (sl *SelfNamedChaincodeLedger) PutSelfChaincodeName(selfChaincodeName string) error {
	err := sl.Ledger.PutString(selfChaincodeNameKey, selfChaincodeName)
	if err != nil {
		return err
	}
	return nil
}

func (sl *SelfNamedChaincodeLedger) GetSelfChaincodeName() (string, error) {
	selfChaincodeNameBytes, err := sl.Ledger.Get(selfChaincodeNameKey)
	if err != nil {
		return "", err
	}
	return string(selfChaincodeNameBytes), nil
}
