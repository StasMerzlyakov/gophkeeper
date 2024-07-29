package tui

import "github.com/StasMerzlyakov/gophkeeper/internal/domain"

type Controller interface {
	LoginEMail(data *domain.EMailData)
	LoginPassOTP(otpPass *domain.OTPPass)
	LoginCheckMasterKey(masterPassword string)
	RegEMail(data *domain.EMailData)
	RegPassOTP(otpPass *domain.OTPPass)
	RegInitMasterKey(mKey *domain.UnencryptedMasterKeyData)

	ShowBankCardList()
	ShowBankCard(num string)
	AddBankCard(bankCard *domain.BankCardView)
	UpdateBankCard(bankCard *domain.BankCardView)
	DeleteBankCard(number string)
	GetBankCard(number string)

	ShowUserPasswordDataList()
	ShowUserPasswordData(hint string)
	AddUserPasswordData(data *domain.UserPasswordData)
	UpdatePasswordData(data *domain.UserPasswordData)
	DeleteUpdatePasswordData(hint string)
	GetUpdatePasswordData(hint string)
}
