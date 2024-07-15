package ttlstorage

import (
	"context"
	"fmt"
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
				println(now.Format(time.RFC1123))
				ttlMap.data.Range(func(k, v interface{}) bool {
					if expEnt, ok := v.(expireEntry); ok {
						if expEnt.ExpiresAt.Before(now) {
							fmt.Printf("delte key %v\n", k)
							ttlMap.data.Delete(k)
						}
					} else {
						ttlMap.data.Delete(k) // unexpected data
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
	Value     interface{}
}

func (t *TTLMap) Store(key string, val interface{}) {
	t.data.Store(key, expireEntry{
		ExpiresAt: time.Now().Add(t.ttl),
		Value:     val,
	})
}

func (t *TTLMap) Delete(key string) {
	t.data.Delete(key)
}

func (t *TTLMap) Load(key string) (val interface{}) {
	entry, ok := t.data.Load(key)
	if !ok {
		return nil
	}

	expireEntry := entry.(expireEntry)
	return expireEntry.Value
}
