package mock

var TransactionContextMockInstance *TransactionContextMock

type MockInstances struct {
	TransactionContextMock *TransactionContextMock
	LedgerFake             *LedgerFake
}

func NewMockInstances() *MockInstances {
	inMemoryLedgerFake := NewLedgerFake()
	transactionContextMock := NewTransactionContextMock(inMemoryLedgerFake)

	return &MockInstances{
		TransactionContextMock: transactionContextMock,
		LedgerFake:             inMemoryLedgerFake,
	}
}

func InitMockInstances() *MockInstances {
	mockInstances := NewMockInstances()
	TransactionContextMockInstance = mockInstances.TransactionContextMock

	return mockInstances
}
