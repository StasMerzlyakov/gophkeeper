package tui

import "github.com/StasMerzlyakov/gophkeeper/internal/domain"

type Controller interface {
	LoginEMail(data *domain.EMailData)
	LoginPassOTP(otpPass *domain.OTPPass)
	LoginCheckMasterKey(masterKeyPassword string)
	RegEMail(data *domain.EMailData)
	RegPassOTP(otpPass *domain.OTPPass)
	RegInitMasterKey(mKey *domain.UnencryptedMasterKeyData)
}
