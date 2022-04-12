package generators

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"strconv"

	"github.com/itchyny/base58-go"
)

// Hashing raw data.
func sha256Of(input string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))
	return algorithm.Sum(nil)
}

// This binary to text encoding will be used to provide the final output of the process.
// Base58 reduces confusion in character output.
func base58Encoded(bytes []byte) (string, error) {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	return string(encoded), err
}

// Hashing initialUrl + userId url with sha256.
// Here userId is added to prevent providing similar shortened urls to separate users.
// Applying base58 on the derived big integer value and pick the first 8 characters.
func GenerateShortLink(initialLink string, userID uint32) (string, error) {
	userIDstr := strconv.Itoa(int(userID))
	urlHashBytes := sha256Of(initialLink + userIDstr)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString, err := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	return finalString[:8], err
}
