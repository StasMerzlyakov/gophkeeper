package storage

import (
	"fmt"
	"sync"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewStorage() *simpleStorage {
	return &simpleStorage{
		status:           domain.ClientStatusOnline,
		userPasswordData: make(map[string]domain.UserPasswordData),
		bankCards:        make(map[string]domain.BankCard),
	}
}

var _ app.AppStorage = (*simpleStorage)(nil)

type simpleStorage struct {
	masterPassword     string
	status             domain.ClientStatus
	userPasswordData   map[string]domain.UserPasswordData
	userPasswordDataMx sync.Mutex
	bankCards          map[string]domain.BankCard
	bankCardsMx        sync.Mutex
}

func (ss *simpleStorage) SetMasterPassword(masterPassword string) {
	ss.masterPassword = masterPassword
}

func (ss *simpleStorage) GetMasterPassword() string {
	return ss.masterPassword
}

func (ss *simpleStorage) AddBankCard(bankCard *domain.BankCard) error {
	ss.bankCardsMx.Lock()
	defer ss.bankCardsMx.Unlock()
	if _, ok := ss.bankCards[bankCard.Number]; ok {
		// Method on client invoked after success server method invokaction, so it's client error.
		return fmt.Errorf("%w bankCard with number %v exists, reopen client", domain.ErrClientInternal, bankCard.Number)
	}
	ss.bankCards[bankCard.Number] = *bankCard
	return nil
}

func (ss *simpleStorage) AddUserPasswordData(data *domain.UserPasswordData) error {
	ss.userPasswordDataMx.Lock()
	defer ss.userPasswordDataMx.Unlock()
	if _, ok := ss.userPasswordData[data.Hint]; ok {
		// Method on client invoked after success server method invokaction, so it's client error.
		return fmt.Errorf("%w userPasswordData with hint %v exists, reopen client", domain.ErrClientInternal, data.Hint)
	}
	ss.userPasswordData[data.Hint] = *data
	return nil
}

func (ss *simpleStorage) UpdateBankCard(bankCard *domain.BankCard) error {
	ss.bankCardsMx.Lock()
	defer ss.bankCardsMx.Unlock()
	if _, ok := ss.bankCards[bankCard.Number]; !ok {
		// Method on client invoked after success server method invokaction, so it's client error.
		return fmt.Errorf("%w bankCard with number %v is not exists, reopen client", domain.ErrClientInternal, bankCard.Number)
	}
	ss.bankCards[bankCard.Number] = *bankCard
	return nil
}

func (ss *simpleStorage) UpdatePasswordData(data *domain.UserPasswordData) error {
	ss.userPasswordDataMx.Lock()
	defer ss.userPasswordDataMx.Unlock()
	if _, ok := ss.userPasswordData[data.Hint]; !ok {
		// Method on client invoked after success server method invokaction, so it's client error.
		return fmt.Errorf("%w userPasswordData with hint %v is not exists, reopen client", domain.ErrClientInternal, data.Hint)
	}
	ss.userPasswordData[data.Hint] = *data
	return nil
}

func (ss *simpleStorage) DeleteBankCard(number string) error {
	ss.bankCardsMx.Lock()
	defer ss.bankCardsMx.Unlock()
	if _, ok := ss.bankCards[number]; !ok {
		// Method on client invoked after success server method invokaction, so it's client error.
		return fmt.Errorf("%w bankCard with number %v is not exists, reopen client", domain.ErrClientInternal, number)
	}
	delete(ss.bankCards, number)
	return nil
}

func (ss *simpleStorage) DeleteUpdatePasswordData(hint string) error {
	ss.userPasswordDataMx.Lock()
	defer ss.userPasswordDataMx.Unlock()
	if _, ok := ss.userPasswordData[hint]; !ok {
		// Method on client invoked after success server method invokaction, so it's client error.
		return fmt.Errorf("%w userPasswordData with hint %v is not exists, reopen client", domain.ErrClientInternal, hint)
	}
	delete(ss.userPasswordData, hint)
	return nil
}

func (ss *simpleStorage) GetBankCard(number string) (*domain.BankCard, error) {
	ss.bankCardsMx.Lock()
	defer ss.bankCardsMx.Unlock()
	if card, ok := ss.bankCards[number]; !ok {
		// Method on client invoked after success server method invokaction, so it's client error.
		return nil, fmt.Errorf("%w bankCard with number %v is not exists, reopen client", domain.ErrClientInternal, number)
	} else {
		return &card, nil
	}
}

func (ss *simpleStorage) GetUpdatePasswordData(hint string) (*domain.UserPasswordData, error) {
	ss.userPasswordDataMx.Lock()
	defer ss.userPasswordDataMx.Unlock()
	if data, ok := ss.userPasswordData[hint]; !ok {
		// Method on client invoked after success server method invokaction, so it's client error.
		return nil, fmt.Errorf("%w userPasswordData with hint %v is not exists, reopen client", domain.ErrClientInternal, hint)
	} else {
		return &data, nil
	}
}

func (ss *simpleStorage) GetBankCardNumberList() []string {
	ss.bankCardsMx.Lock()
	defer ss.bankCardsMx.Unlock()
	keys := make([]string, len(ss.bankCards))
	i := 0
	for k := range ss.bankCards {
		keys[i] = k
		i++
	}
	return keys
}

func (ss *simpleStorage) GetUserPasswordDataList() []string {
	ss.userPasswordDataMx.Lock()
	defer ss.userPasswordDataMx.Unlock()
	keys := make([]string, len(ss.userPasswordData))
	i := 0
	for k := range ss.userPasswordData {
		keys[i] = k
		i++
	}
	return keys
}

func (ss *simpleStorage) SetBankCards(cards []domain.BankCard) {
	ss.bankCardsMx.Lock()
	defer ss.bankCardsMx.Unlock()
	ss.bankCards = make(map[string]domain.BankCard)
	for _, card := range cards {
		ss.bankCards[card.Number] = card
	}
}

func (ss *simpleStorage) SetUserPasswordDatas(datas []domain.UserPasswordData) {
	ss.userPasswordDataMx.Lock()
	defer ss.userPasswordDataMx.Unlock()
	ss.userPasswordData = make(map[string]domain.UserPasswordData)
	for _, data := range datas {
		ss.userPasswordData[data.Hint] = data
	}
}
