//go:build autotest
// +build autotest

package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestRegistration(t *testing.T) {
	suite.Run(t, new(RegistrationTest))
}
