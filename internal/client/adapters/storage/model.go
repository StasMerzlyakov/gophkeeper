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
	masterKey string
	status    domain.ClientStatus
}

func (ss *simpleStorage) SetMasterKey(masterKey string) {
	ss.masterKey = masterKey
}

func (ss *simpleStorage) GetMasterKey() string {
	return ss.masterKey
}
