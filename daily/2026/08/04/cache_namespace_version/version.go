package cache_namespace_version

import (
	"fmt"
	"sync/atomic"
)

// Namespace generates cache keys whose version can invalidate a whole namespace.
type Namespace struct {
	name    string
	version atomic.Uint64
}

func New(name string) *Namespace {
	n := &Namespace{name: name}
	n.version.Store(1)
	return n
}

func (n *Namespace) Key(id string) string {
	return fmt.Sprintf("%s:v%d:%s", n.name, n.version.Load(), id)
}

func (n *Namespace) Bump() uint64    { return n.version.Add(1) }
func (n *Namespace) Version() uint64 { return n.version.Load() }
