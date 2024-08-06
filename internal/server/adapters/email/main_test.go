package email_test

import (
	"log"
	"os"
	"testing"

	smtpmock "github.com/mocktools/go-smtp-mock/v2"
)

var mockServer *smtpmock.Server

func TestMain(m *testing.M) {
	mockServer = smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:       true,
		LogServerActivity: true,
	})

	if err := mockServer.Start(); err != nil {
		log.Fatalf(err.Error())
	}

	defer func() {
		if err := mockServer.Stop(); err != nil {
			log.Fatalf(err.Error())
		}
	}()

	code := m.Run()
	os.Exit(code)
}
