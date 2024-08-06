package app_test

import (
	"context"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/stretchr/testify/assert"
)

func TestSetValue(t *testing.T) {

	val := "jwt"
	ctx := context.Background()
	nCtx := app.SetContextValue(ctx, app.JWTKey, val)

	assert.Equal(t, app.GetContextValue(nCtx, app.JWTKey), val)
	assert.Equal(t, app.GetContextValue(nCtx, app.MasterKey), "")
}
