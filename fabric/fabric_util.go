package fabric

import (
	"encoding/hex"
	"fmt"

	"github.com/beatoz/bprn-sdk-go/generator"
	"github.com/beatoz/bprn-sdk-go/types"
	"github.com/beatoz/bprn-sdk-go/utils"
	"github.com/golang/protobuf/proto"
	"github.com/holiman/uint256"
	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
)

func ConvertChainId(chainId *uint256.Int) string {
	return utils.Uint256ToString(chainId)
}

func GetChainId(stub shim.ChaincodeStubInterface) (*types.ChainId, error) {
	channelId := stub.GetChannelID()
	chainIdInt, err := utils.StringToUint256(channelId)
	if err != nil {
		return nil, fmt.Errorf("failed to convert channel ID to chain ID: %v", err)
	}

	chainId, err := types.NewChainIDFromDecimal(chainIdInt.Dec())
	if err != nil {
		return nil, fmt.Errorf("failed to create chain ID: %v", err)
	}

	return chainId, nil
}

func IsSameChainId(ctx contractapi.TransactionContextInterface, targetChainId *types.ChainId) error {
	chainId, err := GetChainId(ctx.GetStub())
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	if targetChainId.Equal(chainId) {
		return nil
	} else {
		return fmt.Errorf("chainId is not equal")
	}
}

func InvokeChaincode(stub shim.ChaincodeStubInterface, chaincodeName string, methodName string, methodArgs []string) *peer.Response {
	args := [][]byte{[]byte(methodName)}
	for _, methodArg := range methodArgs {
		args = append(args, []byte(methodArg))
	}
	return stub.InvokeChaincode(chaincodeName, args, stub.GetChannelID())
}

type FabricUtil interface {
	CallerChaincodeName() (string, error)
}

type MockFabricUtil struct {
	callerChaincodeName string
}

func NewMockFabricUtil(callerChaincodeName string) *MockFabricUtil {
	return &MockFabricUtil{
		callerChaincodeName: callerChaincodeName,
	}
}

func (mfu *MockFabricUtil) CallerChaincodeName() (string, error) {
	return mfu.callerChaincodeName, nil
}

type FabricUtilImpl struct {
	stub shim.ChaincodeStubInterface
}

func NewFabricUtil(stub shim.ChaincodeStubInterface) *FabricUtilImpl {
	return &FabricUtilImpl{
		stub: stub,
	}
}

func (fu *FabricUtilImpl) CallerChaincodeName() (string, error) {
	sp, err := fu.stub.GetSignedProposal()
	if err != nil {
		fmt.Println("failed to get signed proposal")
		return "", err
	}

	prop := &peer.Proposal{}
	if err := proto.Unmarshal(sp.ProposalBytes, prop); err != nil {
		_ = fmt.Errorf("!!! failed to unmarshal proposal: %v", err)
		fmt.Printf("failed to unmarshal proposal: %v\n", err)
		return "", err
	}

	// Proposal.Payload -> ChaincodeProposalPayload -> Input -> ChaincodeInvocationSpec
	cpp := &peer.ChaincodeProposalPayload{}
	if err := proto.Unmarshal(prop.Payload, cpp); err != nil {
		fmt.Println("failed to unmarshal proposalPayload")
		return "", err
	}

	cis := &peer.ChaincodeInvocationSpec{}
	if err := proto.Unmarshal(cpp.Input, cis); err != nil {
		fmt.Println("failed to unmarshal input")
		return "", err
	}

	fmt.Println("originalChaincode Name: ", cis.ChaincodeSpec.ChaincodeId.Name)
	fmt.Println("originalChaincode Version: ", cis.ChaincodeSpec.ChaincodeId.Version)
	fmt.Println("originalChaincode Path: ", cis.ChaincodeSpec.ChaincodeId.Path)

	return cis.ChaincodeSpec.ChaincodeId.Name, nil
}

func CallerChaincodeName(stub shim.ChaincodeStubInterface) (string, error) {
	sp, err := stub.GetSignedProposal()
	if err != nil {
		fmt.Println("failed to get signed proposal")
		return "", err
	}

	prop := &peer.Proposal{}
	if err := proto.Unmarshal(sp.ProposalBytes, prop); err != nil {
		_ = fmt.Errorf("!!! failed to unmarshal proposal: %v", err)
		fmt.Printf("failed to unmarshal proposal: %v\n", err)
		return "", err
	}

	// Proposal.Payload -> ChaincodeProposalPayload -> Input -> ChaincodeInvocationSpec
	cpp := &peer.ChaincodeProposalPayload{}
	if err := proto.Unmarshal(prop.Payload, cpp); err != nil {
		fmt.Println("failed to unmarshal proposalPayload")
		return "", err
	}

	cis := &peer.ChaincodeInvocationSpec{}
	if err := proto.Unmarshal(cpp.Input, cis); err != nil {
		fmt.Println("failed to unmarshal input")
		return "", err
	}

	fmt.Println("originalChaincode Name: ", cis.ChaincodeSpec.ChaincodeId.Name)
	fmt.Println("originalChaincode Version: ", cis.ChaincodeSpec.ChaincodeId.Version)
	fmt.Println("originalChaincode Path: ", cis.ChaincodeSpec.ChaincodeId.Path)

	return cis.ChaincodeSpec.ChaincodeId.Name, nil
}

func SigVerifyAndSignerAddress(srcHexSig string, chaincodeName string, chaincodeMethodName string, args []string) (*types.Address, error) {

	fmt.Println("srcHexSig: ", srcHexSig)
	fmt.Println("chaincodeName: ", chaincodeName)
	fmt.Println("chaincodeMethodName: ", chaincodeMethodName)
	for _, arg := range args {
		fmt.Println("arg: ", arg)
	}

	sigMsgGenerator := generator.NewSigMsgGenerator()
	sigMsg, err := sigMsgGenerator.GenerateSigMsg(chaincodeName, chaincodeMethodName, args)
	if err != nil {
		return nil, err
	}

	sigMsgHex := hex.EncodeToString(sigMsg)
	fmt.Println(sigMsgHex)

	sigVerifier := generator.NewSigVerifier()
	address, err := sigVerifier.VerifySignature(sigMsg, srcHexSig)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("address: ", address)

	return address, nil
}

func SigVerifyAndSignerAddressFromTxContext(ctx contractapi.TransactionContextInterface, sig string, chaincodeMethodName string, args []string) (*types.Address, error) {
	callerChaincodeName, err := NewFabricUtil(ctx.GetStub()).CallerChaincodeName()
	if err != nil {
		return nil, fmt.Errorf("failed to get caller chaincode name: %w", err)
	}

	return SigVerifyAndSignerAddress(sig, callerChaincodeName, chaincodeMethodName, args)
}
