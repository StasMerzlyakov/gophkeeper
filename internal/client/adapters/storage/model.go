package storage

import (
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewStorage() *simpleStorage {
	return &simpleStorage{
		status: domain.ClientStatusOnline,
	}
}

type simpleStorage struct {
	masterPassword string
	status         domain.ClientStatus
}

func (ss *simpleStorage) SetMasterPassword(masterPassword string) {
	ss.masterPassword = masterPassword
}

func (ss *simpleStorage) GetMasterPassword() string {
	return ss.masterPassword
}
