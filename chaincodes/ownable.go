package chaincodes

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

const ownerKey = "owner"

type OwnableUnauthorizedAccount struct {
	Account string
}

func (e *OwnableUnauthorizedAccount) Error() string {
	return fmt.Sprintf("ownable: unauthorized account %s", e.Account)
}

type OwnableInvalidOwner struct {
	Owner string
}

func (e *OwnableInvalidOwner) Error() string {
	return fmt.Sprintf("ownable: invalid owner %s", e.Owner)
}

type OwnershipTransferredEvent struct {
	PreviousOwner string `json:"previousOwner"`
	NewOwner      string `json:"newOwner"`
}

type OwnableContract struct {
	//contractapi.Contract
}

func (c *OwnableContract) InitOwnable(ctx contractapi.TransactionContextInterface, initialOwner string) error {
	if initialOwner == "" {
		return &OwnableInvalidOwner{Owner: initialOwner}
	}
	return c.transferOwnership(ctx, initialOwner)
}

func (c *OwnableContract) Owner(ctx contractapi.TransactionContextInterface) (string, error) {
	ownerBytes, err := ctx.GetStub().GetState(ownerKey)
	if err != nil {
		return "", fmt.Errorf("failed to get owner: %v", err)
	}
	if ownerBytes == nil {
		return "", nil
	}
	return string(ownerBytes), nil
}

func (c *OwnableContract) CheckOwner(ctx contractapi.TransactionContextInterface, caller string) error {
	fmt.Println("[OwnableContract.CheckOwner] caller:", caller)
	owner, err := c.Owner(ctx)
	if err != nil {
		return err
	}

	if owner != caller {
		fmt.Println("[OwnableContract.CheckOwner] owner:", owner, " caller:", caller)
		return &OwnableUnauthorizedAccount{Account: caller}
	}

	return nil
}

func (c *OwnableContract) RenounceOwnership(ctx contractapi.TransactionContextInterface) error {
	if err := c.CheckOwner(ctx, ""); err != nil {
		return err
	}
	return c.transferOwnership(ctx, "")
}

func (c *OwnableContract) TransferOwnership(ctx contractapi.TransactionContextInterface, newOwner string) error {
	if err := c.CheckOwner(ctx, ""); err != nil {
		return err
	}
	if newOwner == "" {
		return &OwnableInvalidOwner{Owner: newOwner}
	}
	return c.transferOwnership(ctx, newOwner)
}

func (c *OwnableContract) transferOwnership(ctx contractapi.TransactionContextInterface, newOwner string) error {
	oldOwner, err := c.Owner(ctx)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(ownerKey, []byte(newOwner))
	if err != nil {
		return fmt.Errorf("failed to set owner: %v", err)
	}

	// Emit OwnershipTransferred TransferEvent
	event := OwnershipTransferredEvent{
		PreviousOwner: oldOwner,
		NewOwner:      newOwner,
	}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal TransferEvent: %v", err)
	}

	err = ctx.GetStub().SetEvent("OwnershipTransferred", eventBytes)
	if err != nil {
		return fmt.Errorf("failed to set TransferEvent: %v", err)
	}

	return nil
}
