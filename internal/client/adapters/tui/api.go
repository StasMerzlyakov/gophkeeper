package tui

import "github.com/StasMerzlyakov/gophkeeper/internal/domain"

//go:generate mockgen -destination "./generated_mocks_test.go" -package ${GOPACKAGE}_test . LoginController,RegController

type LoginController interface {
	Login(data *domain.EMailData)
	PassOTP(otpPass *domain.OTPPass)
	CheckMasterKey(masterKeyPassword string)
}

type RegController interface {
	Registrate(data *domain.EMailData)
	PassOTP(otpPass *domain.OTPPass)
	InitMasterKey(mKey *domain.UnencryptedMasterKeyData)
}
