// Package agentwallet defines the storage contract shared by services that
// provision and consume Hyperliquid agent keys.
package agentwallet

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const ciphertextPrefix = "kms-agent:v1:"

// EncryptionContext binds one ciphertext to exactly one trading account.
func EncryptionContext(tradingAddress string) (map[string]string, error) {
	if !common.IsHexAddress(tradingAddress) {
		return nil, fmt.Errorf("invalid trading address %q", tradingAddress)
	}
	return map[string]string{
		"purpose":         "hyperliquid-agent-wallet",
		"trading-address": strings.ToLower(common.HexToAddress(tradingAddress).Hex()),
	}, nil
}

// EncodeCiphertext serializes an opaque KMS ciphertext for database storage.
func EncodeCiphertext(ciphertext []byte) (string, error) {
	if len(ciphertext) == 0 {
		return "", errors.New("empty agent-wallet ciphertext")
	}
	return ciphertextPrefix + base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecodeCiphertext parses a database value without decrypting it.
func DecodeCiphertext(encoded string) ([]byte, error) {
	if !strings.HasPrefix(encoded, ciphertextPrefix) {
		return nil, errors.New("unsupported agent-wallet ciphertext format")
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(encoded, ciphertextPrefix))
	if err != nil {
		return nil, fmt.Errorf("decode agent-wallet ciphertext: %w", err)
	}
	if len(decoded) == 0 {
		return nil, errors.New("empty agent-wallet ciphertext")
	}
	return decoded, nil
}

// Generate creates a secp256k1 agent key and returns its address plus a
// fixed-width private-key byte slice. The caller must erase privateKey after it
// has been encrypted.
func Generate() (address common.Address, privateKey []byte, err error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("generate agent key: %w", err)
	}
	return crypto.PubkeyToAddress(key.PublicKey), crypto.FromECDSA(key), nil
}

// Zero erases a mutable private-key byte slice.
func Zero(privateKey []byte) {
	for i := range privateKey {
		privateKey[i] = 0
	}
}
