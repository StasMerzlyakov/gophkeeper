package app

import (
	"context"
	"crypto/rand"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewAppController(conf *config.ClientConf) *appController {
	helper := NewHelper(rand.Read)
	cntr := &appController{
		conf:        conf,
		status:      domain.ClientStatusOnline,
		loginer:     NewLoginer(),
		pinger:      NewPinger(),
		registrator: NewRegistrator(),
	}

	cntr.SetDomainHelper(helper)
	return cntr
}

func (ac *appController) Start() {
	ac.wg.Add(1)
	ac.exitChan = make(chan struct{}, 1)
	go func() {
		defer ac.wg.Done()
		for {
			select {
			case <-time.After(2 * ac.conf.InterationTimeout):
				ac.CheckServerStatus()
			case <-ac.exitChan:
				return
			}
		}
	}()

}

func (ac *appController) CheckServerStatus() {

	ac.invokeFnHlp(func(ctx context.Context) {
		log := GetMainLogger()
		log.Debug("ping start")
		err := ac.pinger.Ping(ctx)
		if err == nil {
			select {
			case <-ctx.Done():
				// timeout
				log.Error("ping timeout")
				ac.status = domain.ClientStatusOffline
				return
			default:
				log.Info("ping success")
				if ac.status == domain.ClientStatusOffline {
					ac.status = domain.ClientStatusOnline
					ac.appView.ShowMsg("Server is online")
				}
			}
		} else {
			log.Error("ping error")
			ac.status = domain.ClientStatusOffline
			return
		}
	}, false)
}

func (ac *appController) Stop() {
	close(ac.exitChan)
	ac.server.Stop()
	ac.wg.Wait()
}

func (ac *appController) SetInfoView(view AppView) *appController {
	ac.appView = view
	ac.loginer.LoginView(view)
	ac.registrator.RegView(view)
	return ac
}

func (ac *appController) SetDomainHelper(helper DomainHelper) *appController {
	ac.helper = helper
	ac.loginer.LoginHelper(helper)
	ac.registrator.RegHelper(helper)
	return ac
}

func (ac *appController) SetAppStorage(storage AppStorage) *appController {
	ac.loginer.LoginStorage(storage)
	ac.storage = storage
	return ac
}

func (ac *appController) SetServer(server AppServer) *appController {
	ac.pinger.SetPinger(server)
	ac.loginer.LoginSever(server)
	ac.registrator.RegServer(server)
	ac.server = server
	return ac
}

type appController struct {
	conf        *config.ClientConf
	status      domain.ClientStatus
	appView     AppView
	helper      DomainHelper
	loginer     *loginer
	pinger      *pinger
	registrator *registrator
	server      AppServer
	exitChan    chan struct{}
	wg          sync.WaitGroup
	storage     AppStorage
}

func (ac *appController) GetStatus() domain.ClientStatus {
	return ac.status
}

func (ac *appController) invokeFnHlp(fn func(ctx context.Context), showErr bool) {
	timedCtx, cancelCtxFn := context.WithTimeout(context.Background(), ac.conf.InterationTimeout)
	defer cancelCtxFn()
	ac.wg.Add(1)
	doneCh := make(chan struct{}, 1)
	go func() {
		defer ac.wg.Done()
		fn(timedCtx)
		select {
		case <-timedCtx.Done():
			return // timeout
		default:
			doneCh <- struct{}{}
		}
	}()

	select {
	case <-timedCtx.Done():
		ac.status = domain.ClientStatusOffline
		if showErr {
			ac.appView.ShowError(fmt.Errorf("%w server timeout", domain.ErrClientServerTimeout))
		}
	case <-ac.exitChan:
		return
	case <-doneCh:
		return
	}
}

func (ac *appController) invokeFn(fn func(ctx context.Context), runStatus domain.ClientStatus, msg string, showErr bool) {
	go func() {
		if ac.status != runStatus {
			if msg != "" {
				ac.appView.ShowMsg(msg)
			}
		} else {
			ac.invokeFnHlp(fn, showErr)

		}
	}()
}

func (ac *appController) LoginEMail(data *domain.EMailData) {
	ac.invokeFn(func(ctx context.Context) {
		ac.loginer.Login(ctx, data)
	}, domain.ClientStatusOnline, "server is offline", true)
}
func (ac *appController) LoginPassOTP(otpPass *domain.OTPPass) {
	ac.invokeFn(func(ctx context.Context) {
		ac.loginer.PassOTP(ctx, otpPass)
	}, domain.ClientStatusOnline, "server is offline", true)
}
func (ac *appController) LoginCheckMasterKey(masterKeyPassword string) {
	ac.invokeFn(func(ctx context.Context) {
		ac.loginer.CheckMasterKey(ctx, masterKeyPassword)
	}, domain.ClientStatusOnline, "server is offline", true)
}
func (ac *appController) RegEMail(data *domain.EMailData) {
	ac.invokeFn(func(ctx context.Context) {
		ac.registrator.Registrate(ctx, data)
	}, domain.ClientStatusOnline, "server is offline", true)
}
func (ac *appController) RegPassOTP(otpPass *domain.OTPPass) {
	ac.invokeFn(func(ctx context.Context) {
		ac.registrator.PassOTP(ctx, otpPass)
	}, domain.ClientStatusOnline, "server is offline", true)
}
func (ac *appController) RegInitMasterKey(mKey *domain.UnencryptedMasterKeyData) {
	ac.invokeFn(func(ctx context.Context) {
		ac.registrator.InitMasterKey(ctx, mKey)
	}, domain.ClientStatusOnline, "server is offline", true)
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
