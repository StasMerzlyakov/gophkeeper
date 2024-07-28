package storage_test

import (
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/adapters/storage"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestApp(t *testing.T) {

	t.Run("masterKey", func(t *testing.T) {
		app := storage.NewStorage()
		mKey := "MasterKey"
		app.SetMasterKey(mKey)
		assert.Equal(t, mKey, app.GetMasterKey())
	})

	t.Run("status", func(t *testing.T) {
		app := storage.NewStorage()
		assert.Equal(t, domain.ClientStatusOnline, app.GetStatus())
		app.SetStatus(domain.ClientStatusOffline)
		assert.Equal(t, domain.ClientStatusOffline, app.GetStatus())
	})

}
