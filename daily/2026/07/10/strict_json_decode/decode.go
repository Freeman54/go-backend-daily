package strictjsondecode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	Name  string `json:"name"`
	Limit int    `json:"limit"`
}

func DecodeRequest(body []byte) (Request, error) {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return Request{}, fmt.Errorf("empty body")
	}

	decoder := json.NewDecoder(bytes.NewReader(trimmed))
	decoder.DisallowUnknownFields()

	var req Request
	if err := decoder.Decode(&req); err != nil {
		return Request{}, normalizeDecodeError(err)
	}

	if err := ensureSingleValue(decoder); err != nil {
		return Request{}, err
	}

	if strings.TrimSpace(req.Name) == "" {
		return Request{}, fmt.Errorf("name is required")
	}
	if req.Limit <= 0 {
		return Request{}, fmt.Errorf("limit must be positive")
	}

	return req, nil
}

func ensureSingleValue(decoder *json.Decoder) error {
	var extra json.RawMessage
	if err := decoder.Decode(&extra); err != nil {
		if err == io.EOF {
			return nil
		}
		return normalizeDecodeError(err)
	}
	return fmt.Errorf("body must contain a single JSON object")
}

func normalizeDecodeError(err error) error {
	switch {
	case strings.Contains(err.Error(), "unknown field"):
		return fmt.Errorf("unknown field: %w", err)
	case strings.Contains(err.Error(), "cannot unmarshal"):
		return fmt.Errorf("type mismatch: %w", err)
	default:
		return err
	}
}
