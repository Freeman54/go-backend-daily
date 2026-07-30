package messageacklease

import (
	"testing"
	"time"
)

func TestLease_ExtendsBeforeSafetyMargin(t *testing.T) {
	now := time.Unix(1000, 0)
	l, err := New(now, 30*time.Second, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if l.ShouldExtend(now.Add(24 * time.Second)) {
		t.Fatal("尚未进入安全边界")
	}
	if !l.ShouldExtend(now.Add(25 * time.Second)) {
		t.Fatal("进入安全边界后应续租")
	}
	l.Extended(now.Add(25*time.Second), 30*time.Second)
	if l.ShouldExtend(now.Add(49 * time.Second)) {
		t.Fatal("续租后截止时间应向后移动")
	}
}

func TestLease_StopsAfterAckOrInvalidInput(t *testing.T) {
	now := time.Unix(1000, 0)
	l, _ := New(now, time.Second, 100*time.Millisecond)
	l.Ack()
	if l.ShouldExtend(now.Add(time.Second)) {
		t.Fatal("确认后不应再续租")
	}
	if _, err := New(now, time.Second, time.Second); err == nil {
		t.Fatal("安全边界必须小于租期")
	}
}
