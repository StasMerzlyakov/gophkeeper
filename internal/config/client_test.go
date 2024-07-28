package config_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientConf(t *testing.T) {

	t.Run("default values", func(t *testing.T) {
		flagSet := flag.NewFlagSet(t.Name(), errorHandling)
		conf, err := config.LoadClientConf(flagSet)
		require.NoError(t, err)

		assert.Equal(t, config.ClientDefaultServerAddres, conf.ServerAddress)
		assert.Equal(t, config.ClientDefaultCACert, conf.CACert)
	})

	t.Run("env values", func(t *testing.T) {

		err := os.Setenv("SERVER_ADDRESS", "http://test")
		defer os.Clearenv()
		require.NoError(t, err)

		err = os.Setenv("CA_CERT", "ca.cert")
		require.NoError(t, err)

		flagSet := flag.NewFlagSet(t.Name(), errorHandling)

		conf, err := config.LoadClientConf(flagSet)
		require.NoError(t, err)

		assert.Equal(t, "http://test", conf.ServerAddress)
		assert.Equal(t, "ca.cert", conf.CACert)
	})

	t.Run("config file", func(t *testing.T) {

		testConfPath := filepath.Join(TestFileDirectory, "clientConf.json")

		currentArgs := os.Args[:]
		defer func() {
			os.Args = currentArgs
		}()

		os.Args = []string{t.Name(), fmt.Sprintf("-config=%v", testConfPath)}
		flagSet := flag.NewFlagSet(t.Name(), errorHandling)
		conf, err := config.LoadClientConf(flagSet)

		require.NoError(t, err)

		assert.Equal(t, "localhost:9193", conf.ServerAddress)
	})

	t.Run("env priority", func(t *testing.T) {

		testConfPath := filepath.Join(TestFileDirectory, "clientConf.json")

		currentArgs := os.Args[:]
		defer func() {
			os.Args = currentArgs
		}()

		os.Args = []string{t.Name(), fmt.Sprintf("-config=%v", testConfPath)}
		flagSet := flag.NewFlagSet(t.Name(), errorHandling)

		err := os.Setenv("SERVER_ADDRESS", "http://test")
		defer os.Clearenv()
		require.NoError(t, err)

		conf, err := config.LoadClientConf(flagSet)

		require.NoError(t, err)

		assert.Equal(t, "http://test", conf.ServerAddress)
	})

	/*
		t.Run("err", func(t *testing.T) {
			os.Setenv("SERVER_ADDRESS", "http://test")
			os.Setenv("CA_CERT", "ca.cert")

			conf, err := config.LoadClientConf()
			require.Error(t, err)

		}) */

}
