package generator

import (
	"encoding/hex"
	"fmt"

	"github.com/beatoz/beatoz-go/types/crypto"
	"github.com/beatoz/bprn-sdk-go/types"
)

type SigVerifier struct {
}

func NewSigVerifier() *SigVerifier {
	return &SigVerifier{}
}

func (sv *SigVerifier) VerifySignature2(sigMsg []byte, sig []byte) (*types.Address, error) {
	btzAddress, compressedPubKey, err := crypto.Sig2Addr(sigMsg, sig)
	if err != nil {
		return nil, err
	}
	_ = compressedPubKey

	fmt.Println("btzAddress: ", btzAddress)

	address, err2 := types.NewAddress(btzAddress.String())
	if err2 != nil {
		return nil, err2
	}

	return address, nil
}

func (sv *SigVerifier) VerifySignature(sigMsg []byte, sig string) (*types.Address, error) {
	//hexSig := strings.TrimPrefix(sig, "0x")

	sigBytes, err := hex.DecodeString(sig)
	if err != nil {
		return nil, err
	}
	fmt.Println("sigBytes: ", sigBytes)

	btzAddress, _, err := crypto.Sig2Addr(sigMsg, sigBytes)
	if err != nil {
		return nil, err
	}

	fmt.Println("btzAddress: ", btzAddress)

	//addr1 := AddressToHex(btzAddress)
	//fmt.Println("addr1: ", addr1)
	//
	//addr := btzAddress.String()
	//fmt.Println("addr: ", addr)
	//
	//_, _ = addr, addr1
	address, err := types.NewAddress(btzAddress.String())
	if err != nil {
		return nil, err
	}

	return address, nil
}
