package storage

import "github.com/StasMerzlyakov/gophkeeper/internal/client/app"

func NewStorage() *simpleStorage {
	return &simpleStorage{}
}

var _ app.LoginStorage = (*simpleStorage)(nil)

type simpleStorage struct {
	masterKey string
}

func (ss *simpleStorage) SetMasterKey(masterKey string) {
	ss.masterKey = masterKey
}

func (ss *simpleStorage) GetMasterKey() string {
	return ss.masterKey
}
