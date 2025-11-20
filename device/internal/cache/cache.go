package cache

import (
	"runtime"
	"sync"

	"github.com/emil-j-olsson/ubiquiti/device/internal/types"
)

type Cache[V any] struct {
	value V
	mu    sync.RWMutex
}

func NewCache[V any](value V) *Cache[V] {
	return &Cache[V]{value: value}
}

func (c *Cache[V]) Get() V {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}

func (c *Cache[V]) Set(value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = value
}

func (c *Cache[V]) Update(fn func(V)) V {
	c.mu.Lock()
	defer c.mu.Unlock()
	fn(c.value)
	return c.value
}

type State struct {
	cache *Cache[types.DeviceState]
}

func NewDeviceState(config types.Config) *State {
	state := types.DeviceState{
		Identifier:     config.Identifier,
		DeviceVersions: config.DeviceVersions,
		Architecture:   runtime.GOARCH,
		OS:             runtime.GOOS,
		DeviceStatus:   types.DeviceStatusBooting,
	}
	return &State{cache: NewCache(state)}
}

// depending on state, we will allow calls to grpc/http... set the status... etc...

// methods to update state
