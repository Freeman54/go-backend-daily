package logsamplingpolicy

import "testing"

func TestPolicy_AlwaysKeepsImportantEvents(t *testing.T) {
	p, err := New(0)
	if err != nil {
		t.Fatal(err)
	}
	for _, event := range []Event{
		{RequestID: "a", StatusCode: 500},
		{RequestID: "b", DurationMillis: 1000},
		{RequestID: "c", Force: true},
	} {
		if !p.Keep(event) {
			t.Fatalf("重要事件被丢弃: %+v", event)
		}
	}
	if p.Keep(Event{RequestID: "ordinary"}) {
		t.Fatal("0 采样率不应保留普通事件")
	}
}

func TestPolicy_DeterministicByRequestID(t *testing.T) {
	p, _ := New(0.25)
	first := p.Keep(Event{RequestID: "same"})
	for i := 0; i < 20; i++ {
		if got := p.Keep(Event{RequestID: "same"}); got != first {
			t.Fatal("相同请求 ID 的采样决策必须稳定")
		}
	}
	all, _ := New(1)
	if !all.Keep(Event{RequestID: "ordinary"}) {
		t.Fatal("100% 采样率应保留普通事件")
	}
	if _, err := New(1.1); err == nil {
		t.Fatal("非法采样率应报错")
	}
}
