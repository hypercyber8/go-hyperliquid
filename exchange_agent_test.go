package hyperliquid

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestApproveAgentAddressUsesCallerPersistedKey(t *testing.T) {
	const agent = "0x2222222222222222222222222222222222222222"
	var payload map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("decode request: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	ownerKey, err := crypto.HexToECDSA(strings.Repeat("1", 64))
	if err != nil {
		t.Fatal(err)
	}
	exchange := NewExchangeWithInfo(context.Background(), NewAccount(ownerKey), server.URL, "", "", nil)
	name := "hypercyber-trading-agent"
	result, err := exchange.ApproveAgentAddress(context.Background(), &name, agent)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "ok" {
		t.Fatalf("status = %q, want ok", result.Status)
	}
	action, ok := payload["action"].(map[string]any)
	if !ok {
		t.Fatalf("action = %#v", payload["action"])
	}
	if got := action["agentAddress"]; got != agent {
		t.Fatalf("agent address = %v, want %s", got, agent)
	}
	if got := action["agentName"]; got != name {
		t.Fatalf("agent name = %v, want %s", got, name)
	}
}

func TestApproveAgentAddressRejectsInvalidAddressBeforeSigning(t *testing.T) {
	ownerKey, err := crypto.HexToECDSA(strings.Repeat("1", 64))
	if err != nil {
		t.Fatal(err)
	}
	exchange := NewExchangeWithInfo(context.Background(), NewAccount(ownerKey), TestnetAPIURL, "", "", nil)
	if _, err = exchange.ApproveAgentAddress(context.Background(), nil, "not-an-address"); err == nil {
		t.Fatal("invalid agent address was accepted")
	}
}

func TestApproveAgentAddressReturnsExchangeRejection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"err","response":"agent name already used"}`))
	}))
	defer server.Close()
	ownerKey, err := crypto.HexToECDSA(strings.Repeat("1", 64))
	if err != nil {
		t.Fatal(err)
	}
	exchange := NewExchangeWithInfo(context.Background(), NewAccount(ownerKey), server.URL, "", "", nil)
	_, err = exchange.ApproveAgentAddress(context.Background(), nil, "0x2222222222222222222222222222222222222222")
	if err == nil || !strings.Contains(err.Error(), "agent name already used") {
		t.Fatalf("error = %v, want exchange rejection", err)
	}
}
