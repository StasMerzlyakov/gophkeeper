package app

import (
	"context"
	"crypto/rand"
	"fmt"
	"strconv"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewAppController(conf *config.ClientConf) *appController {
	helper := NewHelper(rand.Read)
	cntr := &appController{
		conf:        conf,
		loginer:     NewLoginer().LoginHelper(helper),
		registrator: NewRegistrator().RegHelper(helper),
	}

	cntr.SetDomainHelper(helper)
	return cntr
}

func (ac *appController) Start() {
	if ac.server == nil {
		panic("appController is not initialized - server is nil")
	}
	ac.server.Start()
}

func (ac *appController) Stop() {
	if ac.server == nil {
		panic("appController is not initialized - server is nil")
	}
	ac.server.Stop()
}

func (ac *appController) SetInfoView(view AppView) *appController {
	ac.appView = view
	return ac
}

func (ac *appController) SetDomainHelper(helper DomainHelper) *appController {
	ac.helper = helper
	return ac
}

func (ac *appController) SetAppStorage(storage AppStorage) *appController {
	ac.storage = storage
	return ac
}

func (ac *appController) SetServer(server AppServer) *appController {
	ac.server = server
	ac.loginer.LoginSever(server)
	ac.registrator.RegServer(server)
	return ac
}

type appController struct {
	conf        *config.ClientConf
	appView     AppView
	helper      DomainHelper // не нужен
	server      AppServer
	loginer     *loginer
	registrator *registrator
	storage     AppStorage
}

func (ac *appController) invokeFn(fn func(ctx context.Context) error, successFn func()) {
	go func() {
		ctx := context.Background()
		if err := fn(ctx); err != nil {
			ac.appView.ShowMsg(err.Error())
		} else {
			successFn()
		}
	}()
}

func (ac *appController) LoginEMail(data *domain.EMailData) {
	ac.invokeFn(
		func(ctx context.Context) error {
			return ac.loginer.Login(ctx, data)
		},
		func() {
			ac.appView.ShowLogOTPView()
		})
}
func (ac *appController) LoginPassOTP(otpPass *domain.OTPPass) {
	ac.invokeFn(
		func(ctx context.Context) error {
			return ac.loginer.PassOTP(ctx, otpPass)
		},
		func() {
			ac.appView.ShowMasterKeyView("")
		})
}
func (ac *appController) LoginCheckMasterKey(masterKeyPassword string) {
	ac.invokeFn(
		func(ctx context.Context) error {
			err, hint := ac.loginer.CheckMasterKey(ctx, masterKeyPassword)
			if err != nil {
				ac.appView.ShowMasterKeyView(hint)
			}
			return err
		},
		func() {
			ac.storage.SetMasterPassword(masterKeyPassword)
			ac.appView.ShowDataAccessView()
		})
}
func (ac *appController) RegEMail(data *domain.EMailData) {
	ac.invokeFn(
		func(ctx context.Context) error {
			return ac.registrator.Registrate(ctx, data)
		},
		func() {
			ac.appView.ShowRegOTPView()
		})
}
func (ac *appController) RegPassOTP(otpPass *domain.OTPPass) {
	ac.invokeFn(
		func(ctx context.Context) error {
			return ac.registrator.PassOTP(ctx, otpPass)
		},
		func() {
			ac.appView.ShowRegMasterKeyView()
		})
}
func (ac *appController) RegInitMasterKey(mKey *domain.UnencryptedMasterKeyData) {
	ac.invokeFn(func(ctx context.Context) error {
		return ac.registrator.InitMasterKey(ctx, mKey)
	}, func() {
		ac.appView.ShowLoginView()
	})
}

func (ac *appController) AddBankCard(bankCardView *domain.BankCardView) {
	go func() {
		bankCard := &domain.BankCard{
			Number: bankCardView.Number,
			CVV:    bankCardView.CVV,
		}

		expMonth, err := strconv.Atoi(bankCardView.ExpiryMonth)
		if err != nil {
			ac.appView.ShowMsg("Wrong month value")
			return
		}
		bankCard.ExpiryMonth = expMonth

		expEear, err := strconv.Atoi(bankCardView.ExpiryYear)
		if err != nil {
			ac.appView.ShowMsg("Wrong year value")
			return
		}
		if expEear < 100 {
			expEear += 2000
		}
		bankCard.ExpiryYear = expEear

		if err := ac.helper.CheckBankCardData(bankCard); err != nil {
			ac.appView.ShowMsg(fmt.Sprintf("Wrong card data %v", err.Error()))
			return
		}

		if err := ac.storage.AddBankCard(bankCard); err != nil {
			panic(err)
		}
		ac.appView.ShowBankCardListView(ac.storage.GetBankCardNumberList())
	}()
}
func (ac *appController) UpdateBankCard(bankCardView *domain.BankCardView) {
	go func() {
		bankCard := &domain.BankCard{
			Number: bankCardView.Number,
			CVV:    bankCardView.CVV,
		}

		expMonth, err := strconv.Atoi(bankCardView.ExpiryMonth)
		if err != nil {
			ac.appView.ShowMsg("Wrong month value")
			return
		}
		bankCard.ExpiryMonth = expMonth

		expEear, err := strconv.Atoi(bankCardView.ExpiryYear)
		if err != nil {
			ac.appView.ShowMsg("Wrong year value")
			return
		}
		if expEear < 100 {
			expEear += 2000
		}
		bankCard.ExpiryYear = expEear

		if err := ac.helper.CheckBankCardData(bankCard); err != nil {
			ac.appView.ShowMsg(fmt.Sprintf("Wrong card data %v", err.Error()))
			return
		}

		if err := ac.storage.UpdateBankCard(bankCard); err != nil {
			panic(err)
		}
		ac.appView.ShowBankCardListView(ac.storage.GetBankCardNumberList())
	}()
}

func (ac *appController) DeleteBankCard(number string) {
	go func() {
		if err := ac.storage.DeleteBankCard(number); err != nil {
			panic(err)
		}
		ac.appView.ShowBankCardListView(ac.storage.GetBankCardNumberList())
	}()
}

func (ac *appController) GetBankCard(number string) {
	go func() {
		if card, err := ac.storage.GetBankCard(number); err != nil {
			panic(err)
		} else {
			ac.appView.ShowBankCardView(card)
		}
	}()
}

func (ac *appController) ShowBankCard(num string) {
	go func() {
		if num == "" {
			ac.appView.ShowBankCardView(nil)
		} else {
			if data, err := ac.storage.GetBankCard(num); err != nil {
				panic(err)
			} else {
				ac.appView.ShowBankCardView(data)
			}
		}
	}()
}

func (ac *appController) ShowBankCardList() {
	go func() {
		ac.appView.ShowBankCardListView(ac.storage.GetBankCardNumberList())
	}()
}

func (ac *appController) ShowUserPasswordData(hint string) {
	panic("TODO")
}

func (ac *appController) ShowUserPasswordDataList() {
	panic("TODO")
}

func (ac *appController) AddUserPasswordData(data *domain.UserPasswordData) {
	if err := ac.storage.AddUserPasswordData(data); err != nil {
		panic(err)
	}
}
func (ac *appController) UpdatePasswordData(data *domain.UserPasswordData) {
	if err := ac.storage.UpdatePasswordData(data); err != nil {
		panic(err)
	}
}
func (ac *appController) DeleteUpdatePasswordData(hint string) {
	if err := ac.storage.DeleteUpdatePasswordData(hint); err != nil {
		panic(err)
	}
}
func (ac *appController) GetUpdatePasswordData(hint string) {
	panic("TODO")
}
