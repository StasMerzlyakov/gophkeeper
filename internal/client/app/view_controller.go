package app

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

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
		if err := ac.dataAccessor.GetBankCardList(ctx); err != nil {
			ac.appView.ShowMsg(err.Error())
			// do not return error - show cache
		}
		return nil
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
				ac.appView.ShowMsg(err.Error()) // show the error and change page to cardList
				return err
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
				ac.appView.ShowMsg(err.Error()) // show the error and change page to cardList
				return err
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
			// do not return error - show cache
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
			// do not return error - show cache
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
				ac.appView.ShowMsg(err.Error())
				return err
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
				ac.appView.ShowMsg(err.Error())
				return err
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

func (ac *viewController) SaveFile(info *domain.FileInfo) {
	log.Debug("SaveFile start")
	cancelChan := make(chan struct{}, 1)
	errorChan := make(chan error, 1)
	log := GetMainLogger()
	var canceled atomic.Bool

	var progressCount atomic.Int32
	var progressCommon atomic.Int32

	closeProgress := make(chan struct{}, 1)

	go func() {
		defer func() {
			closeProgress <- struct{}{}
		}()

		cancelFnHandler := func() {
			if canceled.CompareAndSwap(false, true) {
				log.Debug("cancel invoked !!!")
				cancelChan <- struct{}{}
			}
		}

		go func() { // separate goroutine for proggress view
			for {
				select {
				case <-closeProgress:
					return
				case <-time.After(1 * time.Second):
					procesed := progressCount.Load()
					common := progressCommon.Load()
					progressText := fmt.Sprintf("loading %d of %d", procesed, common)
					percentage := float64(procesed) * 100 / float64(common)
					if percentage > 100 {
						percentage = 100
					}
					ac.appView.CreateProgressBar(fmt.Sprintf("Loading %s", info.Name), percentage, progressText, cancelFnHandler)
				}
			}
		}()

		progerssFn := func(procesed int, common int) {
			if common > 0 && procesed > 0 {
				progressCount.Store(int32(procesed))
				progressCommon.Store(int32(common))

			}
		}

		ctx := context.Background()

		ac.fileAccessor.LoadFile(ctx, info, progerssFn, cancelChan, errorChan)
		log.Debug("Load complete")
		select {
		case err := <-errorChan:
			ac.appView.ShowMsg(err.Error())
		default:
			ac.GetFilesInfoList()
		}

	}()
}

func (ac *viewController) UploadFile(info *domain.FileInfo) {
	log.Debug("Upload start")
	cancelChan := make(chan struct{}, 1)
	errorChan := make(chan error, 1)
	log := GetMainLogger()
	var canceled atomic.Bool

	var progressCount atomic.Int32
	var progressCommon atomic.Int32

	closeProgress := make(chan struct{}, 1)

	go func() {
		cancelFnHandler := func() {
			if canceled.CompareAndSwap(false, true) {
				log.Debug("cancel invoked !!!")
				cancelChan <- struct{}{}
			}
		}

		go func() { // separate goroutine for proggress view
			for {
				select {
				case <-closeProgress:
					return
				case <-time.After(1 * time.Second):
					procesed := progressCount.Load()
					common := progressCommon.Load()
					progressText := fmt.Sprintf("loading %d of %d", procesed, common)
					percentage := float64(procesed) * 100 / float64(common)
					if percentage > 100 {
						percentage = 100
					}
					ac.appView.CreateProgressBar(fmt.Sprintf("Loading %s", info.Name), percentage, progressText, cancelFnHandler)
				}
			}
		}()
		progerssFn := func(procesed int, common int) {
			if common > 0 && procesed > 0 {
				progressCount.Store(int32(procesed))
				progressCommon.Store(int32(common))

			}
		}

		ctx := context.Background()

		ac.fileAccessor.UploadFile(ctx, info, progerssFn, cancelChan, errorChan)
		log.Debug("Upload complete")
		select {
		case err := <-errorChan:
			ac.appView.ShowMsg(err.Error())
		default:
			ac.GetFilesInfoList()
		}
	}()

}
