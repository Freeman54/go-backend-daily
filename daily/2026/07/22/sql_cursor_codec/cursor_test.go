package sqlcursorcodec

import (
	"testing"
	"time"
)

func TestCodecRoundTripsCursor(t *testing.T) {
	c := Cursor{CreatedAt: time.Date(2026, 7, 22, 3, 4, 5, 123, time.UTC), ID: 42}
	token, err := Encode(c)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	got, err := Decode(token)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got != c {
		t.Fatalf("cursor = %#v, want %#v", got, c)
	}
}

func TestCodecRejectsInvalidCursor(t *testing.T) {
	if _, err := Encode(Cursor{}); err == nil {
		t.Fatal("expected invalid cursor error")
	}
	for _, token := range []string{"", "not-base64"} {
		if _, err := Decode(token); err == nil {
			t.Fatalf("expected decode error for %q", token)
		}
	}
}
