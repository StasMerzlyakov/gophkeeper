package ttlstorage

import (
	"context"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/server/domain"
)

func NewMemStorage(ctx context.Context, conf *config.ServerConf) *memStorage {
	return &memStorage{
		ttlMap: NewTTLMap(ctx, conf.AuthStageTimeout),
	}
}

type memStorage struct {
	ttlMap *TTLMap
}

// Create create new key in storage.
// Returns:
//
//	[domain.ErrDublicateKeyViolation] if key exists.
func (mSt *memStorage) Create(ctx context.Context, sessionID domain.SessionID, data any) error {
	if _, ok := mSt.ttlMap.LoadOrStore(string(sessionID), data); ok {
		return fmt.Errorf("%w - key %s exists", domain.ErrDublicateKeyViolation, sessionID)
	}
	return nil
}

// Load value if exists.
// Returns:
//
//	nil, [domain.ErrDataNotExists] if no value fund.
func (mSt *memStorage) Load(ctx context.Context, sessionID domain.SessionID) (any, error) {
	val, ok := mSt.ttlMap.Load(string(sessionID))
	if !ok {
		return nil, fmt.Errorf("%w - now data found by key %s", domain.ErrDataNotExists, sessionID)
	}
	return val, nil
}

// DeleteAndCreate delete value by sessionID and create new value.
// Returns:
//
//	[domain.ErrDataNotExists] if no value fund by oldSessionID.
//	[domain.ErrDublicateKeyViolation] if key sessionID exists.
func (mSt *memStorage) DeleteAndCreate(ctx context.Context,
	oldSessionID domain.SessionID,
	sessionID domain.SessionID,
	data any,
) error {
	if _, ok := mSt.ttlMap.LoadAndDelete(string(oldSessionID)); !ok {
		return fmt.Errorf("%w - now data found by key %s", domain.ErrDataNotExists, oldSessionID)
	}

	if _, ok := mSt.ttlMap.LoadOrStore(string(sessionID), data); ok {
		return fmt.Errorf("%w - key %s exists", domain.ErrDublicateKeyViolation, sessionID)
	}

	return nil
}

// LoadAndDelete delete and return value by sessionID.
// Returns:
//
//	[domain.ErrDataNotExists] if no value fund by sessionID.
func (mSt *memStorage) LoadAndDelete(ctx context.Context,
	sessionID domain.SessionID,
) (any, error) {
	if val, ok := mSt.ttlMap.LoadAndDelete(string(sessionID)); !ok {
		return nil, fmt.Errorf("%w - now data found by key %s", domain.ErrDataNotExists, sessionID)
	} else {
		return val, nil
	}
}
