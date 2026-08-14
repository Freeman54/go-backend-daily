package workerreaper

import (
	"reflect"
	"testing"
	"time"
)

func TestSelectOldestIdleWorkersAndKeepMinimum(t *testing.T) {
	now := time.Unix(100, 0)
	workers := []Worker{
		{ID: "recent", LastUsed: now.Add(-time.Second)},
		{ID: "old", LastUsed: now.Add(-20 * time.Second)},
		{ID: "older", LastUsed: now.Add(-30 * time.Second)},
	}
	got, err := Select(now, workers, 10*time.Second, 2)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"older"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("Select() = %v, want %v", got, want)
	}
	if workers[0].ID != "recent" {
		t.Fatal("Select should not mutate input order")
	}
}

func TestSelectBoundaryAndValidation(t *testing.T) {
	now := time.Unix(100, 0)
	workers := []Worker{{ID: "idle", LastUsed: now.Add(-10 * time.Second)}}
	got, err := Select(now, workers, 10*time.Second, 0)
	if err != nil || !reflect.DeepEqual(got, []string{"idle"}) {
		t.Fatalf("got=%v err=%v", got, err)
	}
	got, err = Select(now, workers, 10*time.Second, 1)
	if err != nil || got != nil {
		t.Fatalf("minimum should retain worker: got=%v err=%v", got, err)
	}
	if _, err := Select(now, workers, 0, 0); err == nil {
		t.Fatal("zero timeout should fail")
	}
	if _, err := Select(now, workers, time.Second, -1); err == nil {
		t.Fatal("negative minimum should fail")
	}
}
