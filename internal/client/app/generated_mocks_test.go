// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/StasMerzlyakov/gophkeeper/internal/client/app (interfaces: AppServer,AppView,Pinger,DomainHelper,AppStorage)

// Package app_test is a generated GoMock package.
package app_test

import (
	context "context"
	reflect "reflect"

	domain "github.com/StasMerzlyakov/gophkeeper/internal/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockAppServer is a mock of AppServer interface.
type MockAppServer struct {
	ctrl     *gomock.Controller
	recorder *MockAppServerMockRecorder
}

// MockAppServerMockRecorder is the mock recorder for MockAppServer.
type MockAppServerMockRecorder struct {
	mock *MockAppServer
}

// NewMockAppServer creates a new mock instance.
func NewMockAppServer(ctrl *gomock.Controller) *MockAppServer {
	mock := &MockAppServer{ctrl: ctrl}
	mock.recorder = &MockAppServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAppServer) EXPECT() *MockAppServerMockRecorder {
	return m.recorder
}

// CheckEMail mocks base method.
func (m *MockAppServer) CheckEMail(arg0 context.Context, arg1 string) (domain.EMailStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckEMail", arg0, arg1)
	ret0, _ := ret[0].(domain.EMailStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckEMail indicates an expected call of CheckEMail.
func (mr *MockAppServerMockRecorder) CheckEMail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckEMail", reflect.TypeOf((*MockAppServer)(nil).CheckEMail), arg0, arg1)
}

// GetHelloData mocks base method.
func (m *MockAppServer) GetHelloData(arg0 context.Context) (*domain.HelloData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHelloData", arg0)
	ret0, _ := ret[0].(*domain.HelloData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHelloData indicates an expected call of GetHelloData.
func (mr *MockAppServerMockRecorder) GetHelloData(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHelloData", reflect.TypeOf((*MockAppServer)(nil).GetHelloData), arg0)
}

// InitMasterKey mocks base method.
func (m *MockAppServer) InitMasterKey(arg0 context.Context, arg1 *domain.MasterKeyData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InitMasterKey", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InitMasterKey indicates an expected call of InitMasterKey.
func (mr *MockAppServerMockRecorder) InitMasterKey(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitMasterKey", reflect.TypeOf((*MockAppServer)(nil).InitMasterKey), arg0, arg1)
}

// Login mocks base method.
func (m *MockAppServer) Login(arg0 context.Context, arg1 *domain.EMailData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Login indicates an expected call of Login.
func (mr *MockAppServerMockRecorder) Login(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAppServer)(nil).Login), arg0, arg1)
}

// PassLoginOTP mocks base method.
func (m *MockAppServer) PassLoginOTP(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PassLoginOTP", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PassLoginOTP indicates an expected call of PassLoginOTP.
func (mr *MockAppServerMockRecorder) PassLoginOTP(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PassLoginOTP", reflect.TypeOf((*MockAppServer)(nil).PassLoginOTP), arg0, arg1)
}

// PassRegOTP mocks base method.
func (m *MockAppServer) PassRegOTP(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PassRegOTP", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PassRegOTP indicates an expected call of PassRegOTP.
func (mr *MockAppServerMockRecorder) PassRegOTP(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PassRegOTP", reflect.TypeOf((*MockAppServer)(nil).PassRegOTP), arg0, arg1)
}

// Ping mocks base method.
func (m *MockAppServer) Ping(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockAppServerMockRecorder) Ping(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockAppServer)(nil).Ping), arg0)
}

// Registrate mocks base method.
func (m *MockAppServer) Registrate(arg0 context.Context, arg1 *domain.EMailData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Registrate", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Registrate indicates an expected call of Registrate.
func (mr *MockAppServerMockRecorder) Registrate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Registrate", reflect.TypeOf((*MockAppServer)(nil).Registrate), arg0, arg1)
}

// Stop mocks base method.
func (m *MockAppServer) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *MockAppServerMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockAppServer)(nil).Stop))
}

// MockAppView is a mock of AppView interface.
type MockAppView struct {
	ctrl     *gomock.Controller
	recorder *MockAppViewMockRecorder
}

// MockAppViewMockRecorder is the mock recorder for MockAppView.
type MockAppViewMockRecorder struct {
	mock *MockAppView
}

// NewMockAppView creates a new mock instance.
func NewMockAppView(ctrl *gomock.Controller) *MockAppView {
	mock := &MockAppView{ctrl: ctrl}
	mock.recorder = &MockAppViewMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAppView) EXPECT() *MockAppViewMockRecorder {
	return m.recorder
}

// ShowBankCardListView mocks base method.
func (m *MockAppView) ShowBankCardListView(arg0 []string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowBankCardListView", arg0)
}

// ShowBankCardListView indicates an expected call of ShowBankCardListView.
func (mr *MockAppViewMockRecorder) ShowBankCardListView(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowBankCardListView", reflect.TypeOf((*MockAppView)(nil).ShowBankCardListView), arg0)
}

// ShowBankCardView mocks base method.
func (m *MockAppView) ShowBankCardView(arg0 *domain.BankCard) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowBankCardView", arg0)
}

// ShowBankCardView indicates an expected call of ShowBankCardView.
func (mr *MockAppViewMockRecorder) ShowBankCardView(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowBankCardView", reflect.TypeOf((*MockAppView)(nil).ShowBankCardView), arg0)
}

// ShowDataAccessView mocks base method.
func (m *MockAppView) ShowDataAccessView() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowDataAccessView")
}

// ShowDataAccessView indicates an expected call of ShowDataAccessView.
func (mr *MockAppViewMockRecorder) ShowDataAccessView() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowDataAccessView", reflect.TypeOf((*MockAppView)(nil).ShowDataAccessView))
}

// ShowError mocks base method.
func (m *MockAppView) ShowError(arg0 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowError", arg0)
}

// ShowError indicates an expected call of ShowError.
func (mr *MockAppViewMockRecorder) ShowError(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowError", reflect.TypeOf((*MockAppView)(nil).ShowError), arg0)
}

// ShowLogOTPView mocks base method.
func (m *MockAppView) ShowLogOTPView() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowLogOTPView")
}

// ShowLogOTPView indicates an expected call of ShowLogOTPView.
func (mr *MockAppViewMockRecorder) ShowLogOTPView() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowLogOTPView", reflect.TypeOf((*MockAppView)(nil).ShowLogOTPView))
}

// ShowLoginView mocks base method.
func (m *MockAppView) ShowLoginView() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowLoginView")
}

// ShowLoginView indicates an expected call of ShowLoginView.
func (mr *MockAppViewMockRecorder) ShowLoginView() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowLoginView", reflect.TypeOf((*MockAppView)(nil).ShowLoginView))
}

// ShowMasterKeyView mocks base method.
func (m *MockAppView) ShowMasterKeyView(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowMasterKeyView", arg0)
}

// ShowMasterKeyView indicates an expected call of ShowMasterKeyView.
func (mr *MockAppViewMockRecorder) ShowMasterKeyView(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowMasterKeyView", reflect.TypeOf((*MockAppView)(nil).ShowMasterKeyView), arg0)
}

// ShowMsg mocks base method.
func (m *MockAppView) ShowMsg(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowMsg", arg0)
}

// ShowMsg indicates an expected call of ShowMsg.
func (mr *MockAppViewMockRecorder) ShowMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowMsg", reflect.TypeOf((*MockAppView)(nil).ShowMsg), arg0)
}

// ShowRegMasterKeyView mocks base method.
func (m *MockAppView) ShowRegMasterKeyView() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowRegMasterKeyView")
}

// ShowRegMasterKeyView indicates an expected call of ShowRegMasterKeyView.
func (mr *MockAppViewMockRecorder) ShowRegMasterKeyView() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowRegMasterKeyView", reflect.TypeOf((*MockAppView)(nil).ShowRegMasterKeyView))
}

// ShowRegOTPView mocks base method.
func (m *MockAppView) ShowRegOTPView() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowRegOTPView")
}

// ShowRegOTPView indicates an expected call of ShowRegOTPView.
func (mr *MockAppViewMockRecorder) ShowRegOTPView() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowRegOTPView", reflect.TypeOf((*MockAppView)(nil).ShowRegOTPView))
}

// ShowRegView mocks base method.
func (m *MockAppView) ShowRegView() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowRegView")
}

// ShowRegView indicates an expected call of ShowRegView.
func (mr *MockAppViewMockRecorder) ShowRegView() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowRegView", reflect.TypeOf((*MockAppView)(nil).ShowRegView))
}

// MockPinger is a mock of Pinger interface.
type MockPinger struct {
	ctrl     *gomock.Controller
	recorder *MockPingerMockRecorder
}

// MockPingerMockRecorder is the mock recorder for MockPinger.
type MockPingerMockRecorder struct {
	mock *MockPinger
}

// NewMockPinger creates a new mock instance.
func NewMockPinger(ctrl *gomock.Controller) *MockPinger {
	mock := &MockPinger{ctrl: ctrl}
	mock.recorder = &MockPingerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPinger) EXPECT() *MockPingerMockRecorder {
	return m.recorder
}

// Ping mocks base method.
func (m *MockPinger) Ping(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockPingerMockRecorder) Ping(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockPinger)(nil).Ping), arg0)
}

// MockDomainHelper is a mock of DomainHelper interface.
type MockDomainHelper struct {
	ctrl     *gomock.Controller
	recorder *MockDomainHelperMockRecorder
}

// MockDomainHelperMockRecorder is the mock recorder for MockDomainHelper.
type MockDomainHelperMockRecorder struct {
	mock *MockDomainHelper
}

// NewMockDomainHelper creates a new mock instance.
func NewMockDomainHelper(ctrl *gomock.Controller) *MockDomainHelper {
	mock := &MockDomainHelper{ctrl: ctrl}
	mock.recorder = &MockDomainHelperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDomainHelper) EXPECT() *MockDomainHelperMockRecorder {
	return m.recorder
}

// CheckAuthPasswordComplexityLevel mocks base method.
func (m *MockDomainHelper) CheckAuthPasswordComplexityLevel(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAuthPasswordComplexityLevel", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CheckAuthPasswordComplexityLevel indicates an expected call of CheckAuthPasswordComplexityLevel.
func (mr *MockDomainHelperMockRecorder) CheckAuthPasswordComplexityLevel(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAuthPasswordComplexityLevel", reflect.TypeOf((*MockDomainHelper)(nil).CheckAuthPasswordComplexityLevel), arg0)
}

// CheckBankCardData mocks base method.
func (m *MockDomainHelper) CheckBankCardData(arg0 *domain.BankCard) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckBankCardData", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckBankCardData indicates an expected call of CheckBankCardData.
func (mr *MockDomainHelperMockRecorder) CheckBankCardData(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckBankCardData", reflect.TypeOf((*MockDomainHelper)(nil).CheckBankCardData), arg0)
}

// CheckMasterPasswordComplexityLevel mocks base method.
func (m *MockDomainHelper) CheckMasterPasswordComplexityLevel(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckMasterPasswordComplexityLevel", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CheckMasterPasswordComplexityLevel indicates an expected call of CheckMasterPasswordComplexityLevel.
func (mr *MockDomainHelperMockRecorder) CheckMasterPasswordComplexityLevel(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckMasterPasswordComplexityLevel", reflect.TypeOf((*MockDomainHelper)(nil).CheckMasterPasswordComplexityLevel), arg0)
}

// CheckUserPasswordData mocks base method.
func (m *MockDomainHelper) CheckUserPasswordData(arg0 *domain.UserPasswordData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUserPasswordData", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckUserPasswordData indicates an expected call of CheckUserPasswordData.
func (mr *MockDomainHelperMockRecorder) CheckUserPasswordData(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUserPasswordData", reflect.TypeOf((*MockDomainHelper)(nil).CheckUserPasswordData), arg0)
}

// DecryptHello mocks base method.
func (m *MockDomainHelper) DecryptHello(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecryptHello", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DecryptHello indicates an expected call of DecryptHello.
func (mr *MockDomainHelperMockRecorder) DecryptHello(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecryptHello", reflect.TypeOf((*MockDomainHelper)(nil).DecryptHello), arg0, arg1)
}

// EncryptHello mocks base method.
func (m *MockDomainHelper) EncryptHello(arg0, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EncryptHello", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EncryptHello indicates an expected call of EncryptHello.
func (mr *MockDomainHelperMockRecorder) EncryptHello(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EncryptHello", reflect.TypeOf((*MockDomainHelper)(nil).EncryptHello), arg0, arg1)
}

// ParseEMail mocks base method.
func (m *MockDomainHelper) ParseEMail(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseEMail", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// ParseEMail indicates an expected call of ParseEMail.
func (mr *MockDomainHelperMockRecorder) ParseEMail(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseEMail", reflect.TypeOf((*MockDomainHelper)(nil).ParseEMail), arg0)
}

// Random32ByteString mocks base method.
func (m *MockDomainHelper) Random32ByteString() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Random32ByteString")
	ret0, _ := ret[0].(string)
	return ret0
}

// Random32ByteString indicates an expected call of Random32ByteString.
func (mr *MockDomainHelperMockRecorder) Random32ByteString() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Random32ByteString", reflect.TypeOf((*MockDomainHelper)(nil).Random32ByteString))
}

// MockAppStorage is a mock of AppStorage interface.
type MockAppStorage struct {
	ctrl     *gomock.Controller
	recorder *MockAppStorageMockRecorder
}

// MockAppStorageMockRecorder is the mock recorder for MockAppStorage.
type MockAppStorageMockRecorder struct {
	mock *MockAppStorage
}

// NewMockAppStorage creates a new mock instance.
func NewMockAppStorage(ctrl *gomock.Controller) *MockAppStorage {
	mock := &MockAppStorage{ctrl: ctrl}
	mock.recorder = &MockAppStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAppStorage) EXPECT() *MockAppStorageMockRecorder {
	return m.recorder
}

// AddBankCard mocks base method.
func (m *MockAppStorage) AddBankCard(arg0 *domain.BankCard) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddBankCard", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddBankCard indicates an expected call of AddBankCard.
func (mr *MockAppStorageMockRecorder) AddBankCard(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBankCard", reflect.TypeOf((*MockAppStorage)(nil).AddBankCard), arg0)
}

// AddUserPasswordData mocks base method.
func (m *MockAppStorage) AddUserPasswordData(arg0 *domain.UserPasswordData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserPasswordData", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUserPasswordData indicates an expected call of AddUserPasswordData.
func (mr *MockAppStorageMockRecorder) AddUserPasswordData(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserPasswordData", reflect.TypeOf((*MockAppStorage)(nil).AddUserPasswordData), arg0)
}

// DeleteBankCard mocks base method.
func (m *MockAppStorage) DeleteBankCard(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBankCard", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBankCard indicates an expected call of DeleteBankCard.
func (mr *MockAppStorageMockRecorder) DeleteBankCard(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBankCard", reflect.TypeOf((*MockAppStorage)(nil).DeleteBankCard), arg0)
}

// DeleteUpdatePasswordData mocks base method.
func (m *MockAppStorage) DeleteUpdatePasswordData(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUpdatePasswordData", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUpdatePasswordData indicates an expected call of DeleteUpdatePasswordData.
func (mr *MockAppStorageMockRecorder) DeleteUpdatePasswordData(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUpdatePasswordData", reflect.TypeOf((*MockAppStorage)(nil).DeleteUpdatePasswordData), arg0)
}

// GetBankCard mocks base method.
func (m *MockAppStorage) GetBankCard(arg0 string) (*domain.BankCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBankCard", arg0)
	ret0, _ := ret[0].(*domain.BankCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBankCard indicates an expected call of GetBankCard.
func (mr *MockAppStorageMockRecorder) GetBankCard(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBankCard", reflect.TypeOf((*MockAppStorage)(nil).GetBankCard), arg0)
}

// GetBankCardNumberList mocks base method.
func (m *MockAppStorage) GetBankCardNumberList() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBankCardNumberList")
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetBankCardNumberList indicates an expected call of GetBankCardNumberList.
func (mr *MockAppStorageMockRecorder) GetBankCardNumberList() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBankCardNumberList", reflect.TypeOf((*MockAppStorage)(nil).GetBankCardNumberList))
}

// GetMasterPassword mocks base method.
func (m *MockAppStorage) GetMasterPassword() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMasterPassword")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetMasterPassword indicates an expected call of GetMasterPassword.
func (mr *MockAppStorageMockRecorder) GetMasterPassword() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMasterPassword", reflect.TypeOf((*MockAppStorage)(nil).GetMasterPassword))
}

// GetUpdatePasswordData mocks base method.
func (m *MockAppStorage) GetUpdatePasswordData(arg0 string) (*domain.UserPasswordData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUpdatePasswordData", arg0)
	ret0, _ := ret[0].(*domain.UserPasswordData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUpdatePasswordData indicates an expected call of GetUpdatePasswordData.
func (mr *MockAppStorageMockRecorder) GetUpdatePasswordData(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUpdatePasswordData", reflect.TypeOf((*MockAppStorage)(nil).GetUpdatePasswordData), arg0)
}

// GetUserPasswordDataList mocks base method.
func (m *MockAppStorage) GetUserPasswordDataList() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserPasswordDataList")
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetUserPasswordDataList indicates an expected call of GetUserPasswordDataList.
func (mr *MockAppStorageMockRecorder) GetUserPasswordDataList() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserPasswordDataList", reflect.TypeOf((*MockAppStorage)(nil).GetUserPasswordDataList))
}

// SetMasterPassword mocks base method.
func (m *MockAppStorage) SetMasterPassword(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetMasterPassword", arg0)
}

// SetMasterPassword indicates an expected call of SetMasterPassword.
func (mr *MockAppStorageMockRecorder) SetMasterPassword(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMasterPassword", reflect.TypeOf((*MockAppStorage)(nil).SetMasterPassword), arg0)
}

// UpdateBankCard mocks base method.
func (m *MockAppStorage) UpdateBankCard(arg0 *domain.BankCard) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBankCard", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBankCard indicates an expected call of UpdateBankCard.
func (mr *MockAppStorageMockRecorder) UpdateBankCard(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBankCard", reflect.TypeOf((*MockAppStorage)(nil).UpdateBankCard), arg0)
}

// UpdatePasswordData mocks base method.
func (m *MockAppStorage) UpdatePasswordData(arg0 *domain.UserPasswordData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePasswordData", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePasswordData indicates an expected call of UpdatePasswordData.
func (mr *MockAppStorageMockRecorder) UpdatePasswordData(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePasswordData", reflect.TypeOf((*MockAppStorage)(nil).UpdatePasswordData), arg0)
}
