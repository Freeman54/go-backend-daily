package entrycostbudget

import (
	"sync"
	"testing"
)

func TestBudgetReserveAndRelease(t *testing.T) {
	b := New(10)
	if !b.Reserve(7) || b.Reserve(4) || b.Reserve(0) {
		t.Fatal("unexpected admission decision")
	}
	b.Release(3)
	if !b.Reserve(4) || b.Used() != 8 {
		t.Fatalf("used = %d", b.Used())
	}
	b.Release(100)
	if b.Used() != 0 {
		t.Fatalf("used after clamp = %d", b.Used())
	}
}

func TestBudgetConcurrentAdmission(t *testing.T) {
	b := New(50)
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() { defer wg.Done(); b.Reserve(1) }()
	}
	wg.Wait()
	if b.Used() != 50 {
		t.Fatalf("used = %d", b.Used())
	}
}

func TestNegativeBudgetRejects(t *testing.T) {
	if New(-1).Reserve(1) {
		t.Fatal("negative budget admitted an entry")
	}
}
