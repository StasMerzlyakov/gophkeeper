package config_test

import (
	"os"
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientConf(t *testing.T) {

	t.Run("default values", func(t *testing.T) {
		conf, err := config.LoadClientConf()
		require.NoError(t, err)

		assert.Equal(t, config.ClientDefaultServerAddres, conf.ServerAddress)
		assert.Equal(t, config.ClientDefaultCACert, conf.CACert)
	})

	t.Run("env values", func(t *testing.T) {
		err := os.Setenv("SERVER_ADDRESS", "http://test")

		require.NoError(t, err)

		err = os.Setenv("CA_CERT", "ca.cert")
		require.NoError(t, err)

		conf, err := config.LoadClientConf()
		require.NoError(t, err)

		assert.Equal(t, "http://test", conf.ServerAddress)
		assert.Equal(t, "ca.cert", conf.CACert)
	})

	/*
		t.Run("err", func(t *testing.T) {
			os.Setenv("SERVER_ADDRESS", "http://test")
			os.Setenv("CA_CERT", "ca.cert")

			conf, err := config.LoadClientConf()
			require.Error(t, err)

		}) */

}
