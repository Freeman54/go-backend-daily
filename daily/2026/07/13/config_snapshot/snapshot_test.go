package configsnapshot

import (
	"errors"
	"testing"
	"time"
)

func TestNewManagerRejectsInvalidConfig(t *testing.T) {
	_, err := NewManager(Config{MaxOpenConns: 0, Timeout: time.Second})
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("NewManager error = %v want ErrInvalidConfig", err)
	}
}

func TestReloadIncrementsVersion(t *testing.T) {
	manager, err := NewManager(Config{
		MaxOpenConns: 8,
		Timeout:      time.Second,
		Flags: map[string]bool{
			"shadow": false,
		},
	})
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	version, err := manager.Reload(Config{
		MaxOpenConns: 16,
		Timeout:      2 * time.Second,
		Flags: map[string]bool{
			"shadow": true,
		},
	})
	if err != nil {
		t.Fatalf("Reload failed: %v", err)
	}
	if version != 2 {
		t.Fatalf("version = %d want 2", version)
	}
}

func TestLoadReturnsDefensiveCopy(t *testing.T) {
	manager, err := NewManager(Config{
		MaxOpenConns: 4,
		Timeout:      time.Second,
		Flags: map[string]bool{
			"beta": true,
		},
	})
	if err != nil {
		t.Fatalf("NewManager failed: %v", err)
	}

	snapshot := manager.Load()
	snapshot.Config.Flags["beta"] = false

	reloaded := manager.Load()
	if !reloaded.Config.Flags["beta"] {
		t.Fatal("stored config should not be mutated by caller")
	}
}
