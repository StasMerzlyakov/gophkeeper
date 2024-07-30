package app

import (
	"context"
	"crypto/rand"
	"errors"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewAppController(conf *config.ClientConf) *appController {
	helper := NewHelper(rand.Read)
	cntr := &appController{
		conf:         conf,
		loginer:      NewLoginer().LoginHelper(helper),
		registrator:  NewRegistrator().RegHelper(helper),
		dataAccessor: NewDataAccessor().DomainHelper(helper),
		helper:       helper,
	}
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

func (ac *appController) SetAppStorage(storage AppStorage) *appController {
	ac.storage = storage
	ac.dataAccessor.AppStorage(storage)
	return ac
}

func (ac *appController) SetServer(server AppServer) *appController {
	ac.server = server
	ac.loginer.LoginSever(server)
	ac.registrator.RegServer(server)
	ac.dataAccessor.AppSever(server)
	return ac
}

type appController struct {
	conf         *config.ClientConf
	appView      AppView
	helper       DomainHelper // не нужен
	server       AppServer
	loginer      *loginer
	registrator  *registrator
	dataAccessor *dataAccessor
	storage      AppStorage
}

func (ac *appController) invokeFn(fn func(ctx context.Context) error, successFn func()) {
	go func() {
		ctx := context.Background()
		if err := fn(ctx); err != nil {
			ac.appView.ShowMsg(err.Error())
		} else {
			if successFn != nil {
				successFn()
			}
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

func (ac *appController) GetBankCardList() {
	ac.invokeFn(func(ctx context.Context) error {
		if err := ac.dataAccessor.GetBankCardList(ctx); err != nil {
			ac.appView.ShowMsg(err.Error())
			nmbrs := ac.storage.GetBankCardNumberList()
			ac.appView.ShowBankCardListView(nmbrs)
		}
		return nil
	}, func() {
		nmbrs := ac.storage.GetBankCardNumberList()
		ac.appView.ShowBankCardListView(nmbrs)
	})
}

func (ac *appController) AddBankCard(bankCardView *domain.BankCardView) {
	ac.invokeFn(
		func(ctx context.Context) error {
			bankCard, err := bankCardView.ToBankCard()
			if err != nil {
				return err
			}

			if err := ac.dataAccessor.AddBankCard(ctx, bankCard); err != nil {
				if errors.Is(err, domain.ErrClientDataIncorrect) {
					return err
				}
				ac.appView.ShowMsg(err.Error())
			}
			return nil
		}, func() {
			ac.GetBankCardList()
		})
}

func (ac *appController) UpdateBankCard(bankCardView *domain.BankCardView) {
	ac.invokeFn(
		func(ctx context.Context) error {
			bankCard, err := bankCardView.ToBankCard()
			if err != nil {
				return err
			}

			if err := ac.dataAccessor.UpdateBankCard(ctx, bankCard); err != nil {
				if errors.Is(err, domain.ErrClientDataIncorrect) {
					return err
				}
				ac.appView.ShowMsg(err.Error())
			}
			return nil
		},
		func() {
			ac.GetBankCardList()
		})
}

func (ac *appController) DeleteBankCard(number string) {
	ac.invokeFn(
		func(ctx context.Context) error {
			if err := ac.dataAccessor.DeleteBankCard(ctx, number); err != nil {
				ac.appView.ShowMsg(err.Error())
			}
			return nil
		}, func() {
			ac.GetBankCardList()
		})
}

func (ac *appController) GetBankCard(num string) {
	ac.invokeFn(
		func(ctx context.Context) error {
			if num == "" {
				ac.appView.ShowBankCardView(nil)
			} else {
				if data, err := ac.storage.GetBankCard(num); err != nil {
					ac.appView.ShowMsg(err.Error())
					ac.appView.ShowBankCardView(nil)
				} else {
					ac.appView.ShowBankCardView(data)
				}
			}
			return nil
		}, nil)
}

func (ac *appController) GetUserPasswordData(hint string) {
	ac.invokeFn(
		func(ctx context.Context) error {
			if hint == "" {
				ac.appView.ShowUserPasswordDataView(nil)
			} else {
				if data, err := ac.storage.GetUpdatePasswordData(hint); err != nil {
					ac.appView.ShowMsg(err.Error())
					ac.appView.ShowBankCardView(nil)
				} else {
					ac.appView.ShowUserPasswordDataView(data)
				}
			}
			return nil
		}, nil)
}

func (ac *appController) GetUserPasswordDataList() {
	ac.invokeFn(func(ctx context.Context) error {
		if err := ac.dataAccessor.GetUserPasswordDataList(ctx); err != nil {
			ac.appView.ShowMsg(err.Error())
		}
		return nil
	}, func() {
		nmbrs := ac.storage.GetUserPasswordDataList()
		ac.appView.ShowUserPasswordDataListView(nmbrs) // show always
	})
}

func (ac *appController) AddUserPasswordData(data *domain.UserPasswordData) {
	ac.invokeFn(
		func(ctx context.Context) error {
			if err := ac.dataAccessor.AddUserPasswordData(ctx, data); err != nil {
				if errors.Is(err, domain.ErrClientDataIncorrect) {
					return err
				}
				ac.appView.ShowMsg(err.Error())
			}
			return nil
		}, func() {
			ac.GetUserPasswordDataList()
		})
}
func (ac *appController) UpdatePasswordData(data *domain.UserPasswordData) {
	ac.invokeFn(
		func(ctx context.Context) error {
			if err := ac.dataAccessor.UpdateUserPasswordData(ctx, data); err != nil {
				if errors.Is(err, domain.ErrClientDataIncorrect) {
					return err
				}
				ac.appView.ShowMsg(err.Error())
			}
			return nil
		}, func() {
			ac.GetUserPasswordDataList()
		})
}
func (ac *appController) DeleteUpdatePasswordData(hint string) {
	ac.invokeFn(
		func(ctx context.Context) error {
			if err := ac.dataAccessor.DeleteUserPasswordData(ctx, hint); err != nil {
				ac.appView.ShowMsg(err.Error())
			}
			return nil
		}, func() {
			ac.GetUserPasswordDataList()
		})
}
