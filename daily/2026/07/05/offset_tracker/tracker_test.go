package offsettracker

import "testing"

func TestAckAdvancesOnlyWhenOffsetsBecomeContiguous(t *testing.T) {
	tracker := New(10)

	if commitTo, advanced := tracker.Ack(12); advanced || commitTo != 9 {
		t.Fatalf("Ack(12) = (%d, %v), want (9, false)", commitTo, advanced)
	}

	if commitTo, advanced := tracker.Ack(10); !advanced || commitTo != 10 {
		t.Fatalf("Ack(10) = (%d, %v), want (10, true)", commitTo, advanced)
	}

	if commitTo, advanced := tracker.Ack(11); !advanced || commitTo != 12 {
		t.Fatalf("Ack(11) = (%d, %v), want (12, true)", commitTo, advanced)
	}
}

func TestAckIgnoresAlreadyCommittedOffsets(t *testing.T) {
	tracker := New(5)
	if _, advanced := tracker.Ack(5); !advanced {
		t.Fatalf("first Ack(5) should advance")
	}

	if commitTo, advanced := tracker.Ack(4); advanced || commitTo != 5 {
		t.Fatalf("Ack(4) = (%d, %v), want (5, false)", commitTo, advanced)
	}
}
