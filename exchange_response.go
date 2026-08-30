package hyperliquid

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type exchangeResponseEnvelope struct {
	Status   string          `json:"status"`
	Response json.RawMessage `json:"response"`
}

type exchangeResponseBody struct {
	Type string `json:"type"`
}

func validateExchangeResponse(data []byte, expectedType string) error {
	var envelope exchangeResponseEnvelope
	if err := jUnmarshal(data, &envelope); err != nil {
		return fmt.Errorf("invalid exchange response: %w", err)
	}

	switch envelope.Status {
	case "err":
		return fmt.Errorf("exchange action rejected: %s", exchangeResponseError(envelope.Response))
	case "ok":
	default:
		return fmt.Errorf("invalid exchange response status %q", envelope.Status)
	}

	response := bytes.TrimSpace(envelope.Response)
	if len(response) == 0 || bytes.Equal(response, []byte("null")) {
		return fmt.Errorf("exchange action returned status ok without a business response")
	}
	if expectedType == "" {
		return nil
	}

	var body exchangeResponseBody
	if err := jUnmarshal(response, &body); err != nil {
		return fmt.Errorf("invalid exchange business response: %w", err)
	}
	if body.Type != expectedType {
		return fmt.Errorf("unexpected exchange response type %q, want %q", body.Type, expectedType)
	}
	return nil
}

func decodeExchangeResponse(data []byte, expectedType string, result any) error {
	if err := validateExchangeResponse(data, expectedType); err != nil {
		return err
	}
	if err := jUnmarshal(data, result); err != nil {
		return fmt.Errorf("cannot decode exchange response: %w", err)
	}
	return nil
}

func exchangeResponseError(response json.RawMessage) string {
	response = bytes.TrimSpace(response)
	if len(response) == 0 || bytes.Equal(response, []byte("null")) {
		return "unknown error"
	}

	var message string
	if err := jUnmarshal(response, &message); err == nil && message != "" {
		return message
	}
	return string(response)
}
