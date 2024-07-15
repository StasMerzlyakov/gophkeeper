package ttlstorage

import (
	"context"
	"sync"
	"time"
)

const refreshPeriod = 1 * time.Second

type TTLMap struct {
	ttl  time.Duration
	data sync.Map
}

func NewTTLMap(ctx context.Context, ttl time.Duration) *TTLMap {
	ttlMap := &TTLMap{
		ttl:  ttl,
		data: sync.Map{},
	}

	go func() {
		ticker := time.NewTicker(refreshPeriod)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				now := time.Now()
				ttlMap.data.Range(func(k, v any) bool {
					if expEnt, ok := v.(expireEntry); ok {
						if expEnt.ExpiresAt.Before(now) {
							ttlMap.data.Delete(k)
						}
					}
					return true
				})
			}
		}
	}()

	return ttlMap
}

type expireEntry struct {
	ExpiresAt time.Time
	Value     any
}

func (t *TTLMap) Store(key string, val any) {
	t.data.Store(key, expireEntry{
		ExpiresAt: time.Now().Add(t.ttl),
		Value:     val,
	})
}

func (t *TTLMap) Delete(key string) {
	t.data.Delete(key)
}

func (t *TTLMap) LoadAndDelete(key string) (any, bool) {
	entry, ok := t.data.LoadAndDelete(key)
	if !ok {
		return nil, false
	}
	expireEntry := entry.(expireEntry)
	return expireEntry.Value, true
}

func (t *TTLMap) LoadOrStore(key string, val any) (any, bool) {
	entry, ok := t.data.LoadOrStore(key, expireEntry{
		ExpiresAt: time.Now().Add(t.ttl),
		Value:     val,
	})
	expireEntry := entry.(expireEntry)
	return expireEntry.Value, ok
}

func (t *TTLMap) Load(key string) (any, bool) {
	entry, ok := t.data.Load(key)
	if !ok {
		return nil, false
	}

	expireEntry := entry.(expireEntry)
	return expireEntry.Value, true
}
