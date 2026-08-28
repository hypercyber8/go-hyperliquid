package agentwallet

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestCiphertextAndContextContract(t *testing.T) {
	const trading = "0x1111111111111111111111111111111111111111"
	ctx, err := EncryptionContext(trading)
	if err != nil {
		t.Fatal(err)
	}
	if ctx["trading-address"] != trading {
		t.Fatalf("trading context = %q", ctx["trading-address"])
	}
	encoded, err := EncodeCiphertext([]byte("opaque"))
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeCiphertext(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(decoded, []byte("opaque")) {
		t.Fatalf("decoded = %q", decoded)
	}
}

func TestGenerateReturnsMatchingAddressAndErasableBytes(t *testing.T) {
	address, keyBytes, err := Generate()
	if err != nil {
		t.Fatal(err)
	}
	key, err := crypto.ToECDSA(keyBytes)
	if err != nil {
		t.Fatal(err)
	}
	if got := crypto.PubkeyToAddress(key.PublicKey); got != address {
		t.Fatalf("address = %s, want %s", got, address)
	}
	Zero(keyBytes)
	if !bytes.Equal(keyBytes, make([]byte, 32)) {
		t.Fatal("private key bytes were not erased")
	}
}
