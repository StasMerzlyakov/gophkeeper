package storage

import (
	"github.com/StasMerzlyakov/gophkeeper/internal/client/app"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewStorage() *simpleStorage {
	return &simpleStorage{
		status: domain.ClientStatusOnline,
	}
}

var _ app.LoginStorage = (*simpleStorage)(nil)

type simpleStorage struct {
	masterKey string
	status    domain.ClientStatus
}

func (ss *simpleStorage) SetMasterKey(masterKey string) {
	ss.masterKey = masterKey
}

func (ss *simpleStorage) GetMasterKey() string {
	return ss.masterKey
}

func (ss *simpleStorage) SetStatus(status domain.ClientStatus) {
	ss.status = status
}

func (ss *simpleStorage) GetStatus() domain.ClientStatus {
	return ss.status
}
