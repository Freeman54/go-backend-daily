package cursorpagination

import (
	"testing"
	"time"
)

func TestPaginateReturnsStableNextCursor(t *testing.T) {
	t.Parallel()

	items := sampleItems()

	firstPage, next, err := Paginate(items, "", 2)
	if err != nil {
		t.Fatalf("Paginate() first page error = %v", err)
	}
	if len(firstPage) != 2 || firstPage[0].ID != 5 || firstPage[1].ID != 4 {
		t.Fatalf("first page = %#v, want IDs [5 4]", firstPage)
	}
	if next == "" {
		t.Fatal("next cursor is empty")
	}

	secondPage, tail, err := Paginate(items, next, 2)
	if err != nil {
		t.Fatalf("Paginate() second page error = %v", err)
	}
	if len(secondPage) != 2 || secondPage[0].ID != 3 || secondPage[1].ID != 2 {
		t.Fatalf("second page = %#v, want IDs [3 2]", secondPage)
	}
	if tail == "" {
		t.Fatal("tail cursor is empty")
	}
}

func TestPaginateReturnsEmptyWhenCursorPointsToEnd(t *testing.T) {
	t.Parallel()

	items := sampleItems()
	lastCursor, err := EncodeCursor(items[len(items)-1])
	if err != nil {
		t.Fatalf("EncodeCursor() error = %v", err)
	}

	page, next, err := Paginate(items, lastCursor, 2)
	if err != nil {
		t.Fatalf("Paginate() error = %v", err)
	}
	if len(page) != 0 {
		t.Fatalf("page length = %d, want 0", len(page))
	}
	if next != "" {
		t.Fatalf("next cursor = %q, want empty", next)
	}
}

func TestPaginateRejectsInvalidCursor(t *testing.T) {
	t.Parallel()

	if _, _, err := Paginate(sampleItems(), "not-base64", 2); err == nil {
		t.Fatal("Paginate() error = nil, want invalid cursor error")
	}
}

func sampleItems() []Item {
	base := time.Date(2026, 6, 29, 12, 0, 0, 0, time.UTC)
	return []Item{
		{ID: 5, CreatedAt: base.Add(-1 * time.Minute), Title: "e"},
		{ID: 4, CreatedAt: base.Add(-2 * time.Minute), Title: "d"},
		{ID: 3, CreatedAt: base.Add(-2 * time.Minute), Title: "c"},
		{ID: 2, CreatedAt: base.Add(-3 * time.Minute), Title: "b"},
		{ID: 1, CreatedAt: base.Add(-4 * time.Minute), Title: "a"},
	}
}
