package storage

import (
	"fmt"
	"sync"

	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

// NewStorage create simple client storage.
func NewStorage() *simpleStorage {
	return &simpleStorage{
		status:           domain.ClientStatusOnline,
		userPasswordData: make(map[string]domain.UserPasswordData),
		bankCards:        make(map[string]domain.BankCard),
		filesInfo:        make(map[string]domain.FileInfo),
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
	filesInfo          map[string]domain.FileInfo
	filesInfoMx        sync.Mutex
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

func (ss *simpleStorage) AddFileInfo(fileInfo *domain.FileInfo) error {
	ss.filesInfoMx.Lock()
	defer ss.filesInfoMx.Unlock()
	if _, ok := ss.filesInfo[fileInfo.Name]; ok {
		// Method on client invoked after success server method invokaction, so it's client error.
		return fmt.Errorf("%w fileInfo with name %v exists, reopen client", domain.ErrClientInternal, fileInfo.Name)
	}
	ss.filesInfo[fileInfo.Name] = *fileInfo
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

func (ss *simpleStorage) UpdateUserPasswordData(data *domain.UserPasswordData) error {
	ss.userPasswordDataMx.Lock()
	defer ss.userPasswordDataMx.Unlock()
	if _, ok := ss.userPasswordData[data.Hint]; !ok {
		// Method on client invoked after success server method invokaction, so it's client error.
		return fmt.Errorf("%w userPasswordData with hint %v is not exists, reopen client", domain.ErrClientInternal, data.Hint)
	}
	ss.userPasswordData[data.Hint] = *data
	return nil
}

func (ss *simpleStorage) UpdateFileInfo(data *domain.FileInfo) error {
	ss.filesInfoMx.Lock()
	defer ss.filesInfoMx.Unlock()
	if _, ok := ss.filesInfo[data.Name]; !ok {
		// Method on client invoked after success server method invokaction, so it's client error.
		return fmt.Errorf("%w fileInfo with name %v is not exists, reopen client", domain.ErrClientInternal, data.Name)
	}
	ss.filesInfo[data.Name] = *data
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

func (ss *simpleStorage) DeleteUserPasswordData(hint string) error {
	ss.userPasswordDataMx.Lock()
	defer ss.userPasswordDataMx.Unlock()
	if _, ok := ss.userPasswordData[hint]; !ok {
		// Method on client invoked after success server method invokaction, so it's client error.
		return fmt.Errorf("%w userPasswordData with hint %v is not exists, reopen client", domain.ErrClientInternal, hint)
	}
	delete(ss.userPasswordData, hint)
	return nil
}

func (ss *simpleStorage) DeleteFileInfo(name string) error {
	ss.filesInfoMx.Lock()
	defer ss.filesInfoMx.Unlock()
	if _, ok := ss.filesInfo[name]; !ok {
		// Method on client invoked after success server method invokaction, so it's client error.
		return fmt.Errorf("%w fileName with name %v is not exists, reopen client", domain.ErrClientInternal, name)
	}
	delete(ss.filesInfo, name)
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

func (ss *simpleStorage) GetUserPasswordData(hint string) (*domain.UserPasswordData, error) {
	ss.userPasswordDataMx.Lock()
	defer ss.userPasswordDataMx.Unlock()
	if data, ok := ss.userPasswordData[hint]; !ok {
		// Method on client invoked after success server method invokaction, so it's client error.
		return nil, fmt.Errorf("%w userPasswordData with hint %v is not exists, reopen client", domain.ErrClientInternal, hint)
	} else {
		return &data, nil
	}
}

func (ss *simpleStorage) GetFileInfo(name string) (*domain.FileInfo, error) {
	ss.filesInfoMx.Lock()
	defer ss.filesInfoMx.Unlock()
	if info, ok := ss.filesInfo[name]; !ok {
		// Method on client invoked after success server method invokaction, so it's client error.
		return nil, fmt.Errorf("%w fileInfo with name %v is not exists, reopen client", domain.ErrClientInternal, name)
	} else {
		return &info, nil
	}
}

func (ss *simpleStorage) GetFileInfoList() []string {
	ss.filesInfoMx.Lock()
	defer ss.filesInfoMx.Unlock()
	keys := make([]string, len(ss.filesInfo))
	i := 0
	for k := range ss.filesInfo {
		keys[i] = k
		i++
	}
	return keys
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

func (ss *simpleStorage) SetFilesInfo(infs []domain.FileInfo) {
	ss.filesInfoMx.Lock()
	defer ss.filesInfoMx.Unlock()
	ss.filesInfo = make(map[string]domain.FileInfo)
	for _, inf := range infs {
		ss.filesInfo[inf.Name] = inf
	}
}
