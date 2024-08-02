package app

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"sync"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewViewController(conf *config.ClientConf) *viewController {
	helper := NewHelper(rand.Read)
	cntr := &viewController{
		conf:         conf,
		loginer:      NewLoginer().LoginHelper(helper),
		registrator:  NewRegistrator().RegHelper(helper),
		dataAccessor: NewDataAccessor().DomainHelper(helper),
		fileAccessor: NewFileAccessor().DomainHelper(helper),
		helper:       helper,
	}
	return cntr
}

func (ac *viewController) Start() {
	if ac.server == nil {
		panic("appController is not initialized - server is nil")
	}
	ac.server.Start()
}

func (ac *viewController) Stop(stopCtx context.Context) {
	if ac.server == nil {
		panic("appController is not initialized - server is nil")
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ac.server.Stop()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ac.fileAccessor.Stop(stopCtx)
	}()

	wg.Wait()
}

func (ac *viewController) SetInfoView(view AppView) *viewController {
	ac.appView = view
	return ac
}

func (ac *viewController) SetAppStorage(storage AppStorage) *viewController {
	ac.storage = storage
	ac.dataAccessor.AppStorage(storage)
	ac.fileAccessor.AppStorage(storage)
	ac.loginer.LoginStorage(storage)
	return ac
}

func (ac *viewController) SetServer(server AppServer) *viewController {
	ac.server = server
	ac.loginer.LoginSever(server)
	ac.registrator.RegServer(server)
	ac.dataAccessor.AppSever(server)
	ac.fileAccessor.AppServer(server)
	return ac
}

type viewController struct {
	conf         *config.ClientConf
	appView      AppView
	helper       DomainHelper // не нужен
	server       AppServer
	loginer      *loginer
	registrator  *registrator
	dataAccessor *dataAccessor
	fileAccessor *fileAccessor
	storage      AppStorage
}

func (ac *viewController) invokeFn(fn func(ctx context.Context) error, successFn func()) {
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

func (ac *viewController) LoginEMail(data *domain.EMailData) {
	ac.invokeFn(
		func(ctx context.Context) error {
			return ac.loginer.Login(ctx, data)
		},
		func() {
			ac.appView.ShowLogOTPView()
		})
}
func (ac *viewController) LoginPassOTP(otpPass *domain.OTPPass) {
	ac.invokeFn(
		func(ctx context.Context) error {
			return ac.loginer.PassOTP(ctx, otpPass)
		},
		func() {
			ac.appView.ShowMasterKeyView("")
		})
}
func (ac *viewController) LoginCheckMasterKey(masterKeyPassword string) {
	ac.invokeFn(
		func(ctx context.Context) error {
			err, hint := ac.loginer.CheckMasterKey(ctx, masterKeyPassword)
			if err != nil {
				ac.appView.ShowMasterKeyView(hint)
			}
			return err
		},
		func() {
			ac.appView.ShowDataAccessView()
		})
}
func (ac *viewController) RegEMail(data *domain.EMailData) {
	ac.invokeFn(
		func(ctx context.Context) error {
			return ac.registrator.Registrate(ctx, data)
		},
		func() {
			ac.appView.ShowRegOTPView()
		})
}
func (ac *viewController) RegPassOTP(otpPass *domain.OTPPass) {
	ac.invokeFn(
		func(ctx context.Context) error {
			return ac.registrator.PassOTP(ctx, otpPass)
		},
		func() {
			ac.appView.ShowRegMasterKeyView()
		})
}
func (ac *viewController) RegInitMasterKey(mKey *domain.UnencryptedMasterKeyData) {
	ac.invokeFn(func(ctx context.Context) error {
		return ac.registrator.InitMasterKey(ctx, mKey)
	}, func() {
		ac.appView.ShowLoginView()
	})
}

func (ac *viewController) GetBankCardList() {
	ac.invokeFn(func(ctx context.Context) error {
		return ac.dataAccessor.GetBankCardList(ctx)
	}, func() {
		nmbrs := ac.storage.GetBankCardNumberList()
		ac.appView.ShowBankCardListView(nmbrs)
	})
}

func (ac *viewController) AddBankCard(bankCardView *domain.BankCardView) {
	ac.invokeFn(
		func(ctx context.Context) error {
			bankCard, err := bankCardView.ToBankCard()
			if err != nil {
				return err
			}

			if err := ac.dataAccessor.AddBankCard(ctx, bankCard); err != nil {
				if errors.Is(err, domain.ErrClientDataIncorrect) {
					return err // client error - show the error and leave the current page
				}
				ac.appView.ShowMsg(err.Error()) // show the error and change page to cardList
			}
			return nil
		}, func() {
			ac.GetBankCardList()
		})
}

func (ac *viewController) UpdateBankCard(bankCardView *domain.BankCardView) {
	ac.invokeFn(
		func(ctx context.Context) error {
			bankCard, err := bankCardView.ToBankCard()
			if err != nil {
				return err
			}

			if err := ac.dataAccessor.UpdateBankCard(ctx, bankCard); err != nil {
				if errors.Is(err, domain.ErrClientDataIncorrect) {
					return err // client error - show the error and leave the current page
				}
				ac.appView.ShowMsg(err.Error()) // show the error and change page to cardList
			}
			return nil
		},
		func() {
			ac.GetBankCardList()
		})
}

func (ac *viewController) DeleteBankCard(number string) {
	ac.invokeFn(
		func(ctx context.Context) error {
			if err := ac.dataAccessor.DeleteBankCard(ctx, number); err != nil {
				ac.appView.ShowMsg(err.Error()) // show the error and leave the current page
			}
			return nil
		}, func() {
			ac.GetBankCardList()
		})
}

func (ac *viewController) GetBankCard(num string) {
	ac.invokeFn(
		func(ctx context.Context) error {
			if data, err := ac.storage.GetBankCard(num); err != nil {
				return err //nothig to show
			} else {
				ac.appView.ShowEditBankCardView(data)
			}
			return nil
		}, nil)
}

func (ac *viewController) NewBankCard() {
	ac.invokeFn(
		func(ctx context.Context) error {
			return nil
		}, func() {
			ac.appView.ShowNewBankCardView()
		})
}

func (ac *viewController) GetFileInfo(name string) {
	ac.invokeFn(
		func(ctx context.Context) error {
			if data, err := ac.storage.GetFileInfo(name); err != nil {
				return err
			} else {
				ac.appView.ShowFileInfoView(data)
			}
			return nil
		}, nil)
}

func (ac *viewController) SaveFile(info *domain.FileInfo) {
	ac.invokeFn(
		func(ctx context.Context) error {
			if err := ac.fileAccessor.SaveFile(ctx, info); err != nil {
				return err
			}
			return nil
		}, func() {
			ac.GetFilesInfoList()
		})
}

func (ac *viewController) DeleteFile(name string) {
	ac.invokeFn(
		func(ctx context.Context) error {
			if err := ac.fileAccessor.DeleteFile(ctx, name); err != nil {
				ac.appView.ShowMsg(err.Error()) //
			}
			return nil
		}, func() {
			ac.GetFilesInfoList()
		})
}

// GetUserPasswordData invoked by tui view
func (ac *viewController) GetUserPasswordData(hint string) {
	ac.invokeFn(
		func(ctx context.Context) error {
			if data, err := ac.storage.GetUserPasswordData(hint); err != nil {
				return err //nothig to show
			} else {
				ac.appView.ShowEditUserPasswordDataView(data)
			}
			return nil
		}, nil)
}

func (ac *viewController) NewUserPasswordData() {
	ac.invokeFn(
		func(ctx context.Context) error {
			return nil
		}, func() {
			ac.appView.ShowNewUserPasswordDataView()
		})
}

func (ac *viewController) GetFilesInfoList() {
	ac.invokeFn(func(ctx context.Context) error {
		if err := ac.fileAccessor.GetFileInfoList(ctx); err != nil {
			ac.appView.ShowMsg(err.Error())
		}
		return nil
	}, func() {
		infos := ac.storage.GetFileInfoList()
		ac.appView.ShowFileInfoListView(infos) // show always
	})
}

func (ac *viewController) GetUserPasswordDataList() {
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

func (ac *viewController) AddUserPasswordData(data *domain.UserPasswordData) {
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
func (ac *viewController) UpdatePasswordData(data *domain.UserPasswordData) {
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
func (ac *viewController) DeleteUpdatePasswordData(hint string) {
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

func (ac *viewController) UploadFile(info *domain.FileInfo) {

	var cancelFnHandler func()

	cancelFn := func() {
		if cancelFnHandler != nil {
			cancelFnHandler()
		}
	}

	resultHandler := func(err error) {
		if err != nil {
			ac.appView.ShowMsg(err.Error())
		} else {
			ac.GetFilesInfoList()
		}
	}

	progerssFn := func(done int, common int) {
		go func() {
			percentage := float64(done*100) / float64(common)
			progressText := fmt.Sprintf("uploading %d of %d", done, common)
			ac.appView.ShowProgressBar(fmt.Sprintf("Uploading %s", info.Name), progressText, percentage, cancelFn)
		}()
	}

	ctx := context.Background()
	if hndl, err := ac.fileAccessor.UploadFile(ctx, info, resultHandler, progerssFn); err != nil {
		ac.appView.ShowMsg(err.Error())
	} else {
		cancelFnHandler = hndl
	}
}
