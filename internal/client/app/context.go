package app

import (
	"context"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

const SessionIDKey = domain.ContextKey("SessionID")
const JWTKey = domain.ContextKey("JWT")
const MasterKey = domain.ContextKey("MasterKey")

func SetContextValue(ctx context.Context, key domain.ContextKey, val string) context.Context {
	resultCtx := context.WithValue(ctx, key, val)
	return resultCtx
}

func GetContextValue(ctx context.Context, key domain.ContextKey) string {
	if val := ctx.Value(key); val != nil {
		if res, ok := val.(string); ok {
			return res
		}
	}
	return ""
}
