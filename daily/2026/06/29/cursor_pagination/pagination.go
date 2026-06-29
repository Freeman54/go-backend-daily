package cursorpagination

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

type Item struct {
	ID        int64
	CreatedAt time.Time
	Title     string
}

type cursor struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func EncodeCursor(item Item) (string, error) {
	payload, err := json.Marshal(cursor{
		ID:        item.ID,
		CreatedAt: item.CreatedAt.UTC(),
	})
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(payload), nil
}

func DecodeCursor(raw string) (Item, error) {
	if raw == "" {
		return Item{}, nil
	}

	payload, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return Item{}, err
	}

	var c cursor
	if err := json.Unmarshal(payload, &c); err != nil {
		return Item{}, err
	}
	if c.CreatedAt.IsZero() {
		return Item{}, errors.New("cursor missing created_at")
	}

	return Item{ID: c.ID, CreatedAt: c.CreatedAt}, nil
}

func Paginate(items []Item, after string, limit int) ([]Item, string, error) {
	if limit <= 0 {
		return nil, "", errors.New("limit must be positive")
	}

	cursorItem, err := DecodeCursor(after)
	if err != nil {
		return nil, "", err
	}

	start := 0
	if after != "" {
		for i, item := range items {
			if isBefore(item, cursorItem) {
				start = i
				break
			}
			start = len(items)
		}
	}

	if start >= len(items) {
		return nil, "", nil
	}

	end := start + limit
	if end > len(items) {
		end = len(items)
	}

	page := append([]Item(nil), items[start:end]...)
	if end == len(items) {
		return page, "", nil
	}

	nextCursor, err := EncodeCursor(page[len(page)-1])
	if err != nil {
		return nil, "", err
	}
	return page, nextCursor, nil
}

func isBefore(item Item, after Item) bool {
	if item.CreatedAt.Before(after.CreatedAt) {
		return true
	}
	if item.CreatedAt.Equal(after.CreatedAt) && item.ID < after.ID {
		return true
	}
	return false
}
