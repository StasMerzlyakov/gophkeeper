package storage_test

import (
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/adapters/storage"
	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {

	t.Run("masterKey", func(t *testing.T) {
		app := storage.NewStorage()
		mKey := "MasterKey"
		app.SetMasterKey(mKey)
		assert.Equal(t, mKey, app.GetMasterKey())
	})

}
