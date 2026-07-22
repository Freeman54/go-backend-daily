package sqlcursorcodec

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type Cursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        int64     `json:"id"`
}

func Encode(cursor Cursor) (string, error) {
	if cursor.CreatedAt.IsZero() || cursor.ID <= 0 {
		return "", fmt.Errorf("invalid cursor")
	}
	payload, err := json.Marshal(cursor)
	if err != nil {
		return "", fmt.Errorf("marshal cursor: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(payload), nil
}

func Decode(token string) (Cursor, error) {
	payload, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return Cursor{}, fmt.Errorf("decode token: %w", err)
	}
	var cursor Cursor
	if err := json.Unmarshal(payload, &cursor); err != nil {
		return Cursor{}, fmt.Errorf("unmarshal cursor: %w", err)
	}
	if cursor.CreatedAt.IsZero() || cursor.ID <= 0 {
		return Cursor{}, fmt.Errorf("invalid cursor")
	}
	return cursor, nil
}
