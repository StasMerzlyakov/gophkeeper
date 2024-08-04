package app

import (
	"context"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

//go:generate mockgen -destination "./generated_mocks_test.go" -package ${GOPACKAGE}_test . AppServer,AppView,Pinger,DomainHelper,AppStorage

type Pinger interface {
	Ping(ctx context.Context) error
}

type AppStorage interface {
	SetMasterPassword(masterPassword string)
	GetMasterPassword() string

	SetBankCards(cards []domain.BankCard)
	AddBankCard(bankCard *domain.BankCard) error
	UpdateBankCard(bankCard *domain.BankCard) error
	DeleteBankCard(number string) error
	GetBankCard(number string) (*domain.BankCard, error)
	GetBankCardNumberList() []string

	GetUserPasswordData(hint string) (*domain.UserPasswordData, error)
	UpdateUserPasswordData(data *domain.UserPasswordData) error
	DeleteUserPasswordData(hint string) error
	SetUserPasswordDatas(datas []domain.UserPasswordData)
	AddUserPasswordData(data *domain.UserPasswordData) error
	GetUserPasswordDataList() []string

	AddFileInfo(fileInfo *domain.FileInfo) error
	UpdateFileInfo(data *domain.FileInfo) error
	DeleteFileInfo(name string) error
	GetFileInfo(name string) (*domain.FileInfo, error)
	GetFileInfoList() []domain.FileInfo
	SetFilesInfo(infs []domain.FileInfo)
	IsFileInfoExists(name string) bool
}

type DomainHelper interface {
	ParseEMail(address string) bool
	CheckAuthPasswordComplexityLevel(pass string) bool
	CheckMasterPasswordComplexityLevel(pass string) bool
	EncryptHello(masterPass string, hello string) (string, error)
	Random32ByteString() string
	DecryptHello(masterPassword string, helloEncrypted string) error
	CheckBankCardData(data *domain.BankCard) error
	CheckUserPasswordData(data *domain.UserPasswordData) error
	EncryptShortData(masterKey string, data string) (string, error)
	DecryptShortData(masterKey string, ciphertext string) (string, error)

	CheckFileForRead(info *domain.FileInfo) error
	CheckFileForWrite(inf *domain.FileInfo) error

	CreateStreamFileReader(info *domain.FileInfo) (domain.StreamFileReader, error)
	CreateStreamFileWriter(dir string) (domain.StreamFileWriter, error)

	CreateChunkEncrypter(password string) domain.ChunkEncrypter
	CreateChunkDecrypter(password string) domain.ChunkDecrypter
}

type AppView interface {
	ShowError(err error)
	ShowMsg(msg string)
	ShowLogOTPView()
	ShowMasterKeyView(hint string)
	ShowDataAccessView()
	ShowLoginView()
	ShowRegView()
	ShowRegOTPView()
	ShowRegMasterKeyView()
	ShowBankCardListView(cardsNumber []string)
	ShowEditBankCardView(bankCard *domain.BankCard)
	ShowNewBankCardView()
	ShowUserPasswordDataListView(hints []string)
	ShowEditUserPasswordDataView(data *domain.UserPasswordData)
	ShowNewUserPasswordDataView()
	ShowFileInfoView(info *domain.FileInfo)
	ShowFileInfoListView(filesInfoList []domain.FileInfo)

	CreateProgressBar(title string, percentage float64, progressText string, cancelFn func())
}

type AppServer interface {
	Stop()
	Start()
	Ping(ctx context.Context) error
	CheckEMail(ctx context.Context, email string) (domain.EMailStatus, error)
	Registrate(ctx context.Context, data *domain.EMailData) error
	PassRegOTP(ctx context.Context, otpPass string) error
	InitMasterKey(ctx context.Context, mKey *domain.MasterKeyData) error
	Login(ctx context.Context, data *domain.EMailData) error
	PassLoginOTP(ctx context.Context, otpPass string) error
	GetHelloData(ctx context.Context) (*domain.HelloData, error)

	GetBankCardList(ctx context.Context) ([]domain.EncryptedBankCard, error)
	CreateBankCard(ctx context.Context, bnkCard *domain.EncryptedBankCard) error
	UpdateBankCard(ctx context.Context, bnkCard *domain.EncryptedBankCard) error
	DeleteBankCard(ctx context.Context, number string) error

	GetUserPasswordDataList(ctx context.Context) ([]domain.EncryptedUserPasswordData, error)
	CreateUserPasswordData(ctx context.Context, data *domain.EncryptedUserPasswordData) error
	UpdateUserPasswordData(ctx context.Context, data *domain.EncryptedUserPasswordData) error
	DeleteUserPasswordData(ctx context.Context, hint string) error

	GetFileInfoList(ctx context.Context) ([]domain.FileInfo, error)
	DeleteFileInfo(ctx context.Context, name string) error
	CreateFileSender(ctx context.Context) (domain.StreamFileWriter, error)
	CreateFileReceiver(ctx context.Context, name string) (domain.StreamFileReader, error)
}
