package message_batch_barrier

import "testing"

func TestBarrier_AdvancesOnlyAcrossContiguousAcknowledgements(t *testing.T) {
	b := New(10)
	if got := b.Ack(12); got != 9 {
		t.Fatalf("Ack(12) = %d, want 9", got)
	}
	if got := b.Ack(10); got != 10 {
		t.Fatalf("Ack(10) = %d, want 10", got)
	}
	if got := b.Ack(11); got != 12 {
		t.Fatalf("Ack(11) = %d, want 12", got)
	}
}

func TestBarrier_IgnoresOldAndDuplicateAcknowledgements(t *testing.T) {
	b := New(5)
	b.Ack(5)
	if got := b.Ack(5); got != 5 {
		t.Fatalf("duplicate Ack() = %d, want 5", got)
	}
	if got := b.Ack(4); got != 5 {
		t.Fatalf("old Ack() = %d, want 5", got)
	}
}

func TestBarrier_PendingReportsGaps(t *testing.T) {
	b := New(20)
	b.Ack(22)
	b.Ack(24)
	if got := b.Pending(); got != 2 {
		t.Fatalf("Pending() = %d, want 2", got)
	}
}
