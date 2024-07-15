package config_test

import (
	"os"
	"testing"
	"time"

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

		assert.Equal(t, config.ServerDefaultTokenExp, conf.TokenExp)
		assert.Equal(t, config.ServerDefaultTokenSecret, conf.TokenSecret)
		assert.Equal(t, config.ServerDefaultAuthTimeout, conf.AuthTimeout)

		assert.Equal(t, config.ServerDefaultMasterKey, conf.MasterKey)
		assert.Equal(t, config.ServerDefaultDomainName, conf.DomainName)
	})

	t.Run("env values", func(t *testing.T) {
		os.Setenv("PORT", ":9191")
		os.Setenv("TLS_KEY", "test.key")
		os.Setenv("TLS_CERT", "test.cert")

		os.Setenv("JWT_EXP", "1h")
		os.Setenv("JWT_SECRET", "pass")
		os.Setenv("AUTH_TIMEOUT", "5m")

		os.Setenv("MASTER_KEY", "key")
		os.Setenv("DOMAIN_NAME", "example.com")

		conf, err := config.LoadServConf()
		require.NoError(t, err)

		assert.Equal(t, ":9191", conf.Port)
		assert.Equal(t, "test.cert", conf.TLSCert)
		assert.Equal(t, "test.key", conf.TLSKey)

		assert.Equal(t, 1*time.Hour, conf.TokenExp)
		assert.Equal(t, "pass", conf.TokenSecret)
		assert.Equal(t, 5*time.Minute, conf.AuthTimeout)
		assert.Equal(t, "key", conf.MasterKey)
		assert.Equal(t, "example.com", conf.DomainName)
	})
}
