package message_partition_router

import "testing"

func TestRouter_StableAndBounded(t *testing.T) {
	r, err := New(8)
	if err != nil {
		t.Fatal(err)
	}
	first := r.Partition("tenant-42")
	if first < 0 || first >= 8 {
		t.Fatalf("partition = %d, want [0,8)", first)
	}
	for i := 0; i < 10; i++ {
		if got := r.Partition("tenant-42"); got != first {
			t.Fatalf("partition changed: %d -> %d", first, got)
		}
	}
}

func TestNew_RejectsNonPositivePartitionCount(t *testing.T) {
	if _, err := New(0); err == nil {
		t.Fatal("分区数为零应返回错误")
	}
}
