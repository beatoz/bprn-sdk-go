package mock

import (
	"sync"

	"github.com/hyperledger/fabric-chaincode-go/v2/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
)

func NewTransactionContextMock(ledgerFake *LedgerFake) *TransactionContextMock {
	//chaincodeStub := &ChaincodeStub{}
	chaincodeStub2 := &ChaincodeStub2{}
	chaincodeStub2.CreateCompositeKeyCalls(ledgerFake.CreateCompositeKey)
	chaincodeStub2.PutStateCalls(ledgerFake.PutState)
	chaincodeStub2.GetStateCalls(ledgerFake.GetState)
	chaincodeStub2.DelStateCalls(ledgerFake.DeleteState)
	//chaincodeStub2.InvokeChaincodeCalls(InvokeChaincodeFake)
	chaincodeStub2.GetChannelIDCalls(func() string { return "TestChannel" })

	transactionContextMock := &TransactionContextMock{}
	transactionContextMock.GetStubReturns(chaincodeStub2)

	return transactionContextMock
}

type TransactionContextMock struct {
	GetClientIdentityStub        func() cid.ClientIdentity
	getClientIdentityMutex       sync.RWMutex
	getClientIdentityArgsForCall []struct {
	}
	getClientIdentityReturns struct {
		result1 cid.ClientIdentity
	}
	getClientIdentityReturnsOnCall map[int]struct {
		result1 cid.ClientIdentity
	}
	GetStubStub        func() shim.ChaincodeStubInterface
	getStubMutex       sync.RWMutex
	getStubArgsForCall []struct {
	}
	getStubReturns struct {
		result1 shim.ChaincodeStubInterface
	}
	getStubReturnsOnCall map[int]struct {
		result1 shim.ChaincodeStubInterface
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *TransactionContextMock) GetClientIdentity() cid.ClientIdentity {
	fake.getClientIdentityMutex.Lock()
	ret, specificReturn := fake.getClientIdentityReturnsOnCall[len(fake.getClientIdentityArgsForCall)]
	fake.getClientIdentityArgsForCall = append(fake.getClientIdentityArgsForCall, struct {
	}{})
	stub := fake.GetClientIdentityStub
	fakeReturns := fake.getClientIdentityReturns
	fake.recordInvocation("GetClientIdentity", []interface{}{})
	fake.getClientIdentityMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *TransactionContextMock) GetClientIdentityCallCount() int {
	fake.getClientIdentityMutex.RLock()
	defer fake.getClientIdentityMutex.RUnlock()
	return len(fake.getClientIdentityArgsForCall)
}

func (fake *TransactionContextMock) GetClientIdentityCalls(stub func() cid.ClientIdentity) {
	fake.getClientIdentityMutex.Lock()
	defer fake.getClientIdentityMutex.Unlock()
	fake.GetClientIdentityStub = stub
}

func (fake *TransactionContextMock) GetClientIdentityReturns(result1 cid.ClientIdentity) {
	fake.getClientIdentityMutex.Lock()
	defer fake.getClientIdentityMutex.Unlock()
	fake.GetClientIdentityStub = nil
	fake.getClientIdentityReturns = struct {
		result1 cid.ClientIdentity
	}{result1}
}

func (fake *TransactionContextMock) GetClientIdentityReturnsOnCall(i int, result1 cid.ClientIdentity) {
	fake.getClientIdentityMutex.Lock()
	defer fake.getClientIdentityMutex.Unlock()
	fake.GetClientIdentityStub = nil
	if fake.getClientIdentityReturnsOnCall == nil {
		fake.getClientIdentityReturnsOnCall = make(map[int]struct {
			result1 cid.ClientIdentity
		})
	}
	fake.getClientIdentityReturnsOnCall[i] = struct {
		result1 cid.ClientIdentity
	}{result1}
}

func (fake *TransactionContextMock) GetStub() shim.ChaincodeStubInterface {
	fake.getStubMutex.Lock()
	ret, specificReturn := fake.getStubReturnsOnCall[len(fake.getStubArgsForCall)]
	fake.getStubArgsForCall = append(fake.getStubArgsForCall, struct {
	}{})
	stub := fake.GetStubStub
	fakeReturns := fake.getStubReturns
	fake.recordInvocation("GetStub", []interface{}{})
	fake.getStubMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *TransactionContextMock) GetStubCallCount() int {
	fake.getStubMutex.RLock()
	defer fake.getStubMutex.RUnlock()
	return len(fake.getStubArgsForCall)
}

func (fake *TransactionContextMock) GetStubCalls(stub func() shim.ChaincodeStubInterface) {
	fake.getStubMutex.Lock()
	defer fake.getStubMutex.Unlock()
	fake.GetStubStub = stub
}

func (fake *TransactionContextMock) GetStubReturns(result1 shim.ChaincodeStubInterface) {
	fake.getStubMutex.Lock()
	defer fake.getStubMutex.Unlock()
	fake.GetStubStub = nil
	fake.getStubReturns = struct {
		result1 shim.ChaincodeStubInterface
	}{result1}
}

func (fake *TransactionContextMock) GetStubReturnsOnCall(i int, result1 shim.ChaincodeStubInterface) {
	fake.getStubMutex.Lock()
	defer fake.getStubMutex.Unlock()
	fake.GetStubStub = nil
	if fake.getStubReturnsOnCall == nil {
		fake.getStubReturnsOnCall = make(map[int]struct {
			result1 shim.ChaincodeStubInterface
		})
	}
	fake.getStubReturnsOnCall[i] = struct {
		result1 shim.ChaincodeStubInterface
	}{result1}
}

func (fake *TransactionContextMock) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getClientIdentityMutex.RLock()
	defer fake.getClientIdentityMutex.RUnlock()
	fake.getStubMutex.RLock()
	defer fake.getStubMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *TransactionContextMock) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}
