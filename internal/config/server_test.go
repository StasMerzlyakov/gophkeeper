package config_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TestFileDirectory = "../../testdata"

const errorHandling = flag.ErrorHandling(4) // хак для игнорирования флагов запуска тестов

func TestLoadServConf(t *testing.T) {

	t.Run("default values", func(t *testing.T) {
		flagSet := flag.NewFlagSet(t.Name(), errorHandling)
		conf, err := config.LoadServConf(flagSet)
		require.NoError(t, err)

		assert.Equal(t, config.ServerDefaultPort, conf.Port)
		assert.Equal(t, config.ServerDefaultTLSCert, conf.TLSCert)
		assert.Equal(t, config.ServerDefaultTLSKey, conf.TLSKey)

		assert.Equal(t, config.ServerDefaultTokenExp, conf.TokenExp)
		assert.Equal(t, config.ServerDefaultTokenSecret, conf.TokenSecret)
		assert.Equal(t, config.ServerDefaultAuthStageTimeout, conf.AuthStageTimeout)

		assert.Equal(t, config.ServerDefaultServerSecret, conf.ServerSecret)
		assert.Equal(t, config.ServerDefaultDomainName, conf.DomainName)

		assert.Equal(t, config.ServerDefaultSMTPHost, conf.SMTPHost)
		assert.Equal(t, config.ServerDefaultSMTPPort, conf.SMTPPort)

		assert.Equal(t, config.ServerDefaultSMTPUsername, conf.SMTPUsername)
		assert.Equal(t, config.ServerDefaultSMTPPassword, conf.SMTPPassword)

		assert.Equal(t, config.ServerDefaultServerEMail, conf.ServerEMail)

		assert.Equal(t, config.ServerDefaultMaxConns, conf.MaxConns)
		assert.Equal(t, config.ServerDefaultMaxConnLifetime, conf.MaxConnLifetime)
		assert.Equal(t, config.ServerDefaultMaxConnIdleTime, conf.MaxConnIdleTime)
	})

	t.Run("env values durations", func(t *testing.T) {
		defer os.Clearenv()

		err := os.Setenv("JWT_EXP", "1h")
		require.NoError(t, err)

		flagSet := flag.NewFlagSet(t.Name(), errorHandling)
		conf, err := config.LoadServConf(flagSet)

		require.NoError(t, err)

		assert.Equal(t, 1*time.Hour, conf.TokenExp)
	})

	t.Run("env values", func(t *testing.T) {
		defer os.Clearenv()
		err := os.Setenv("PORT", ":9191")
		require.NoError(t, err)

		err = os.Setenv("TLS_KEY", "test.key")
		require.NoError(t, err)

		err = os.Setenv("TLS_CERT", "test.cert")
		require.NoError(t, err)

		err = os.Setenv("JWT_EXP", "1h")
		require.NoError(t, err)

		err = os.Setenv("JWT_SECRET", "pass")
		require.NoError(t, err)

		err = os.Setenv("AUTH_STAGE_TIMEOUT", "4m")
		require.NoError(t, err)

		err = os.Setenv("SERVER_SECRET", "key")
		require.NoError(t, err)

		err = os.Setenv("DOMAIN_NAME", "example.com")
		require.NoError(t, err)

		err = os.Setenv("SMTP_HOST", "127.0.0.1")
		require.NoError(t, err)

		err = os.Setenv("SMTP_USERNAME", "john.doe")
		require.NoError(t, err)

		err = os.Setenv("SMTP_PASSWORD", "s3cr3t")
		require.NoError(t, err)

		err = os.Setenv("SMTP_PORT", "26")
		require.NoError(t, err)

		err = os.Setenv("SERVER_EMAIL", "gopheer@localhost")
		require.NoError(t, err)

		err = os.Setenv("DATABASE_DN", "db_uri")
		require.NoError(t, err)

		err = os.Setenv("DATABASE_MAX_CONNS", "3")
		require.NoError(t, err)

		err = os.Setenv("DATABASE_MAX_CONN_LIFE_TIME", "1m")
		require.NoError(t, err)

		err = os.Setenv("DATABASE_MAX_CONN_IDLE_TIME", "2m")
		require.NoError(t, err)

		flagSet := flag.NewFlagSet(t.Name(), errorHandling)
		conf, err := config.LoadServConf(flagSet)

		require.NoError(t, err)

		assert.Equal(t, ":9191", conf.Port)
		assert.Equal(t, "test.cert", conf.TLSCert)
		assert.Equal(t, "test.key", conf.TLSKey)

		assert.Equal(t, 1*time.Hour, conf.TokenExp)
		assert.Equal(t, "pass", conf.TokenSecret)
		assert.Equal(t, 4*time.Minute, conf.AuthStageTimeout)
		assert.Equal(t, "key", conf.ServerSecret)
		assert.Equal(t, "example.com", conf.DomainName)

		assert.Equal(t, "127.0.0.1", conf.SMTPHost)
		assert.Equal(t, 26, conf.SMTPPort)
		assert.Equal(t, "gopheer@localhost", conf.ServerEMail)

		assert.Equal(t, "john.doe", conf.SMTPUsername)
		assert.Equal(t, "s3cr3t", conf.SMTPPassword)

		assert.Equal(t, "db_uri", conf.DatabaseDN)
		assert.Equal(t, 3, conf.MaxConns)
		assert.Equal(t, 1*time.Minute, conf.MaxConnLifetime)
		assert.Equal(t, 2*time.Minute, conf.MaxConnIdleTime)
	})

	t.Run("env rewrite", func(t *testing.T) {
		err := os.Setenv("PORT", ":9191")
		require.NoError(t, err)

		defer os.Clearenv()
		testConfPath := filepath.Join(TestFileDirectory, "serverConf.json")

		currentArgs := os.Args[:]
		defer func() {
			os.Args = currentArgs
		}()

		os.Args = []string{t.Name(), fmt.Sprintf("-config=%v", testConfPath)}
		flagSet := flag.NewFlagSet(t.Name(), errorHandling)
		conf, err := config.LoadServConf(flagSet)

		require.NoError(t, err)

		assert.Equal(t, ":9191", conf.Port)
		assert.Equal(t, 45*time.Second, conf.AuthStageTimeout)
	})

	t.Run("err", func(t *testing.T) {
		err := os.Setenv("SMTP_PORT", "26asdasd")
		require.NoError(t, err)

		defer os.Clearenv()

		flagSet := flag.NewFlagSet(t.Name(), errorHandling)
		conf, err := config.LoadServConf(flagSet)

		require.Nil(t, conf)
		require.Error(t, err)
	})
}

func TestFlag1(t *testing.T) {
	testConfPath := filepath.Join(TestFileDirectory, "serverConf.json")

	currentArgs := os.Args[:]
	defer func() {
		os.Args = currentArgs
	}()

	os.Args = []string{t.Name(), fmt.Sprintf("-config=%v", testConfPath)}
	flagSet := flag.NewFlagSet(t.Name(), errorHandling)
	conf, err := config.LoadServConf(flagSet)

	require.NoError(t, err)

	assert.Equal(t, ":9192", conf.Port)
}

func TestFlag2(t *testing.T) {

	err := os.Setenv("PORT", ":9191")
	require.NoError(t, err)

	defer os.Clearenv()
	testConfPath := filepath.Join(TestFileDirectory, "serverConf.json")

	currentArgs := os.Args[:]
	defer func() {
		os.Args = currentArgs
	}()

	os.Args = []string{t.Name(), fmt.Sprintf("-config=%v", testConfPath)}
	flagSet := flag.NewFlagSet(t.Name(), errorHandling)
	conf, err := config.LoadServConf(flagSet)

	require.NoError(t, err)

	assert.Equal(t, ":9191", conf.Port)
}
