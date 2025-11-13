package cache

import (
	"runtime"
	"sync"
	"time"

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

func (c *Cache[V]) Update(fn func(*V)) V {
	c.mu.Lock()
	defer c.mu.Unlock()
	fn(&c.value)
	return c.value
}

type State struct {
	cache *Cache[types.DeviceState]
}

func NewDeviceState(config types.Config) *State {
	state := types.DeviceState{
		Identifier:         config.Identifier,
		SupportedProtocols: config.SupportedProtocols,
		DeviceVersions:     config.DeviceVersions,
		Architecture:       runtime.GOARCH,
		OS:                 runtime.GOOS,
		DeviceStatus:       types.DeviceStatusBooting,
		StreamInterval:     config.StreamInterval,
		Updated:            time.Now(),
	}
	return &State{cache: NewCache(state)}
}

func (s *State) GetState() types.DeviceState {
	return s.cache.Get()
}

func (s *State) SetState(state types.DeviceState) {
	state.Updated = time.Now()
	s.cache.Set(state)
}

func (s *State) UpdateState(fn func(*types.DeviceState)) types.DeviceState {
	return s.cache.Update(func(state *types.DeviceState) {
		fn(state)
		state.Updated = time.Now()
	})
}
