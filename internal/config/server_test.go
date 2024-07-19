package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		assert.Equal(t, config.ServerDefaultAuthStageTimeout, conf.AuthStageTimeout)

		assert.Equal(t, config.ServerDefaultServerEncryptionKey, conf.ServerEncryptionKey)
		assert.Equal(t, config.ServerDefaultDomainName, conf.DomainName)

		assert.Equal(t, config.ServerDefaultSMTPHost, conf.SMTPHost)
		assert.Equal(t, config.ServerDefaultSMTPPort, conf.SMTPPort)
		assert.Equal(t, config.ServerDefaultServerEMail, conf.ServerEMail)
	})

	t.Run("env values", func(t *testing.T) {
		os.Setenv("PORT", ":9191")
		os.Setenv("TLS_KEY", "test.key")
		os.Setenv("TLS_CERT", "test.cert")

		os.Setenv("JWT_EXP", "1h")
		os.Setenv("JWT_SECRET", "pass")
		os.Setenv("AUTH_STAGE_TIMEOUT", "4m")

		os.Setenv("SERVER_ENCRYPTION_KEY", "key")
		os.Setenv("DOMAIN_NAME", "example.com")

		os.Setenv("SMTP_HOST", "127.0.0.1")
		os.Setenv("SMTP_PORT", "26")
		os.Setenv("SERVER_EMAIL", "gopheer@localhost")

		conf, err := config.LoadServConf()
		require.NoError(t, err)

		assert.Equal(t, ":9191", conf.Port)
		assert.Equal(t, "test.cert", conf.TLSCert)
		assert.Equal(t, "test.key", conf.TLSKey)

		assert.Equal(t, 1*time.Hour, conf.TokenExp)
		assert.Equal(t, "pass", conf.TokenSecret)
		assert.Equal(t, 4*time.Minute, conf.AuthStageTimeout)
		assert.Equal(t, "key", conf.ServerEncryptionKey)
		assert.Equal(t, "example.com", conf.DomainName)

		assert.Equal(t, "127.0.0.1", conf.SMTPHost)
		assert.Equal(t, 26, conf.SMTPPort)
		assert.Equal(t, "gopheer@localhost", conf.ServerEMail)
	})

	t.Run("err", func(t *testing.T) {
		os.Setenv("SMTP_PORT", "26asdasd")

		conf, err := config.LoadServConf()
		require.Nil(t, conf)
		require.Error(t, err)
	})
}
