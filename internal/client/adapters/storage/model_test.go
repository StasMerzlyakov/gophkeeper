package storage_test

import (
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/adapters/storage"
	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {

	t.Run("masterPassword", func(t *testing.T) {
		app := storage.NewStorage()
		mKey := "MasterPassword"
		app.SetMasterPassword(mKey)
		assert.Equal(t, mKey, app.GetMasterPassword())
	})
}
