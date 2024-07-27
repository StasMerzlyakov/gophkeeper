package tui

import "github.com/StasMerzlyakov/gophkeeper/internal/domain"

//go:generate mockgen -destination "./generated_mocks_test.go" -package ${GOPACKAGE}_test . Controller

type Controller interface {
	Login(*domain.EMailData)
}
