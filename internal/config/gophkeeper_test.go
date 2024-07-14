package config_test

import (
	"os"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func TestLoadServConf(t *testing.T) {

	t.Run("default values", func(t *testing.T) {
		conf, err := config.LoadServConf()
		require.NoError(t, err)

		assert.Equal(t, config.ServerDefaultPort, conf.Port)
		assert.Equal(t, config.ServerDefaultTLSCert, conf.TLSCert)
		assert.Equal(t, config.ServerDefaultTLSKey, conf.TLSKey)
	})

	t.Run("env values", func(t *testing.T) {
		os.Setenv("PORT", ":9191")
		os.Setenv("TLS_KEY", "test.key")
		os.Setenv("TLS_CERT", "test.cert")

		conf, err := config.LoadServConf()
		require.NoError(t, err)

		assert.Equal(t, ":9191", conf.Port)
		assert.Equal(t, "test.cert", conf.TLSCert)
		assert.Equal(t, "test.key", conf.TLSKey)
	})
}
