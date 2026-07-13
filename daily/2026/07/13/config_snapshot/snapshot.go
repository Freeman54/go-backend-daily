package configsnapshot

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var ErrInvalidConfig = errors.New("invalid config")

type Config struct {
	MaxOpenConns int
	Timeout      time.Duration
	Flags        map[string]bool
}

type Snapshot struct {
	Version uint64
	Config  Config
}

type Manager struct {
	mu sync.Mutex
	v  atomic.Value
}

func NewManager(initial Config) (*Manager, error) {
	if err := validate(initial); err != nil {
		return nil, err
	}
	manager := &Manager{}
	manager.v.Store(Snapshot{
		Version: 1,
		Config:  cloneConfig(initial),
	})
	return manager, nil
}

func (m *Manager) Load() Snapshot {
	current := m.v.Load().(Snapshot)
	return Snapshot{
		Version: current.Version,
		Config:  cloneConfig(current.Config),
	}
}

func (m *Manager) Reload(next Config) (uint64, error) {
	if err := validate(next); err != nil {
		return 0, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	current := m.v.Load().(Snapshot)
	updated := Snapshot{
		Version: current.Version + 1,
		Config:  cloneConfig(next),
	}
	m.v.Store(updated)
	return updated.Version, nil
}

func validate(cfg Config) error {
	if cfg.MaxOpenConns <= 0 {
		return ErrInvalidConfig
	}
	if cfg.Timeout <= 0 {
		return ErrInvalidConfig
	}
	return nil
}

func cloneConfig(cfg Config) Config {
	flags := make(map[string]bool, len(cfg.Flags))
	for key, value := range cfg.Flags {
		flags[key] = value
	}
	return Config{
		MaxOpenConns: cfg.MaxOpenConns,
		Timeout:      cfg.Timeout,
		Flags:        flags,
	}
}
