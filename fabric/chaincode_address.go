package fabric

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func GenerateChaincodeAddress(channelName string, chaincodeName string) string {
	hashed := sha256.Sum256([]byte(channelName + "-" + chaincodeName))
	address := hex.EncodeToString(hashed[len(hashed)-20:])

	return strings.ToLower(address)
}
