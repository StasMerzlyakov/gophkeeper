// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/StasMerzlyakov/gophkeeper/internal/client/app (interfaces: Server,InfoView,Pinger,RegServer,RegView,RegHelper,LoginServer,LoginView,LoginHelper,AppStorage)

// Package app_test is a generated GoMock package.
package app_test

import (
	context "context"
	reflect "reflect"

	domain "github.com/StasMerzlyakov/gophkeeper/internal/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockServer is a mock of Server interface.
type MockServer struct {
	ctrl     *gomock.Controller
	recorder *MockServerMockRecorder
}

// MockServerMockRecorder is the mock recorder for MockServer.
type MockServerMockRecorder struct {
	mock *MockServer
}

// NewMockServer creates a new mock instance.
func NewMockServer(ctrl *gomock.Controller) *MockServer {
	mock := &MockServer{ctrl: ctrl}
	mock.recorder = &MockServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServer) EXPECT() *MockServerMockRecorder {
	return m.recorder
}

// CheckEMail mocks base method.
func (m *MockServer) CheckEMail(arg0 context.Context, arg1 string) (domain.EMailStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckEMail", arg0, arg1)
	ret0, _ := ret[0].(domain.EMailStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckEMail indicates an expected call of CheckEMail.
func (mr *MockServerMockRecorder) CheckEMail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckEMail", reflect.TypeOf((*MockServer)(nil).CheckEMail), arg0, arg1)
}

// GetHelloData mocks base method.
func (m *MockServer) GetHelloData(arg0 context.Context) (*domain.HelloData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHelloData", arg0)
	ret0, _ := ret[0].(*domain.HelloData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHelloData indicates an expected call of GetHelloData.
func (mr *MockServerMockRecorder) GetHelloData(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHelloData", reflect.TypeOf((*MockServer)(nil).GetHelloData), arg0)
}

// InitMasterKey mocks base method.
func (m *MockServer) InitMasterKey(arg0 context.Context, arg1 *domain.MasterKeyData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InitMasterKey", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InitMasterKey indicates an expected call of InitMasterKey.
func (mr *MockServerMockRecorder) InitMasterKey(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitMasterKey", reflect.TypeOf((*MockServer)(nil).InitMasterKey), arg0, arg1)
}

// Login mocks base method.
func (m *MockServer) Login(arg0 context.Context, arg1 *domain.EMailData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Login indicates an expected call of Login.
func (mr *MockServerMockRecorder) Login(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockServer)(nil).Login), arg0, arg1)
}

// PassLoginOTP mocks base method.
func (m *MockServer) PassLoginOTP(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PassLoginOTP", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PassLoginOTP indicates an expected call of PassLoginOTP.
func (mr *MockServerMockRecorder) PassLoginOTP(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PassLoginOTP", reflect.TypeOf((*MockServer)(nil).PassLoginOTP), arg0, arg1)
}

// PassRegOTP mocks base method.
func (m *MockServer) PassRegOTP(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PassRegOTP", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PassRegOTP indicates an expected call of PassRegOTP.
func (mr *MockServerMockRecorder) PassRegOTP(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PassRegOTP", reflect.TypeOf((*MockServer)(nil).PassRegOTP), arg0, arg1)
}

// Ping mocks base method.
func (m *MockServer) Ping(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockServerMockRecorder) Ping(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockServer)(nil).Ping), arg0)
}

// Registrate mocks base method.
func (m *MockServer) Registrate(arg0 context.Context, arg1 *domain.EMailData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Registrate", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Registrate indicates an expected call of Registrate.
func (mr *MockServerMockRecorder) Registrate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Registrate", reflect.TypeOf((*MockServer)(nil).Registrate), arg0, arg1)
}

// Stop mocks base method.
func (m *MockServer) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *MockServerMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockServer)(nil).Stop))
}

// MockInfoView is a mock of InfoView interface.
type MockInfoView struct {
	ctrl     *gomock.Controller
	recorder *MockInfoViewMockRecorder
}

// MockInfoViewMockRecorder is the mock recorder for MockInfoView.
type MockInfoViewMockRecorder struct {
	mock *MockInfoView
}

// NewMockInfoView creates a new mock instance.
func NewMockInfoView(ctrl *gomock.Controller) *MockInfoView {
	mock := &MockInfoView{ctrl: ctrl}
	mock.recorder = &MockInfoViewMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInfoView) EXPECT() *MockInfoViewMockRecorder {
	return m.recorder
}

// ShowDataAccessView mocks base method.
func (m *MockInfoView) ShowDataAccessView() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowDataAccessView")
}

// ShowDataAccessView indicates an expected call of ShowDataAccessView.
func (mr *MockInfoViewMockRecorder) ShowDataAccessView() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowDataAccessView", reflect.TypeOf((*MockInfoView)(nil).ShowDataAccessView))
}

// ShowError mocks base method.
func (m *MockInfoView) ShowError(arg0 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowError", arg0)
}

// ShowError indicates an expected call of ShowError.
func (mr *MockInfoViewMockRecorder) ShowError(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowError", reflect.TypeOf((*MockInfoView)(nil).ShowError), arg0)
}

// ShowLogOTPView mocks base method.
func (m *MockInfoView) ShowLogOTPView() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowLogOTPView")
}

// ShowLogOTPView indicates an expected call of ShowLogOTPView.
func (mr *MockInfoViewMockRecorder) ShowLogOTPView() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowLogOTPView", reflect.TypeOf((*MockInfoView)(nil).ShowLogOTPView))
}

// ShowLoginView mocks base method.
func (m *MockInfoView) ShowLoginView() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowLoginView")
}

// ShowLoginView indicates an expected call of ShowLoginView.
func (mr *MockInfoViewMockRecorder) ShowLoginView() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowLoginView", reflect.TypeOf((*MockInfoView)(nil).ShowLoginView))
}

// ShowMasterKeyView mocks base method.
func (m *MockInfoView) ShowMasterKeyView(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowMasterKeyView", arg0)
}

// ShowMasterKeyView indicates an expected call of ShowMasterKeyView.
func (mr *MockInfoViewMockRecorder) ShowMasterKeyView(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowMasterKeyView", reflect.TypeOf((*MockInfoView)(nil).ShowMasterKeyView), arg0)
}

// ShowMsg mocks base method.
func (m *MockInfoView) ShowMsg(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowMsg", arg0)
}

// ShowMsg indicates an expected call of ShowMsg.
func (mr *MockInfoViewMockRecorder) ShowMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowMsg", reflect.TypeOf((*MockInfoView)(nil).ShowMsg), arg0)
}

// ShowRegMasterKeyView mocks base method.
func (m *MockInfoView) ShowRegMasterKeyView() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowRegMasterKeyView")
}

// ShowRegMasterKeyView indicates an expected call of ShowRegMasterKeyView.
func (mr *MockInfoViewMockRecorder) ShowRegMasterKeyView() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowRegMasterKeyView", reflect.TypeOf((*MockInfoView)(nil).ShowRegMasterKeyView))
}

// ShowRegOTPView mocks base method.
func (m *MockInfoView) ShowRegOTPView() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowRegOTPView")
}

// ShowRegOTPView indicates an expected call of ShowRegOTPView.
func (mr *MockInfoViewMockRecorder) ShowRegOTPView() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowRegOTPView", reflect.TypeOf((*MockInfoView)(nil).ShowRegOTPView))
}

// ShowRegView mocks base method.
func (m *MockInfoView) ShowRegView() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowRegView")
}

// ShowRegView indicates an expected call of ShowRegView.
func (mr *MockInfoViewMockRecorder) ShowRegView() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowRegView", reflect.TypeOf((*MockInfoView)(nil).ShowRegView))
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

// MockRegServer is a mock of RegServer interface.
type MockRegServer struct {
	ctrl     *gomock.Controller
	recorder *MockRegServerMockRecorder
}

// MockRegServerMockRecorder is the mock recorder for MockRegServer.
type MockRegServerMockRecorder struct {
	mock *MockRegServer
}

// NewMockRegServer creates a new mock instance.
func NewMockRegServer(ctrl *gomock.Controller) *MockRegServer {
	mock := &MockRegServer{ctrl: ctrl}
	mock.recorder = &MockRegServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRegServer) EXPECT() *MockRegServerMockRecorder {
	return m.recorder
}

// CheckEMail mocks base method.
func (m *MockRegServer) CheckEMail(arg0 context.Context, arg1 string) (domain.EMailStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckEMail", arg0, arg1)
	ret0, _ := ret[0].(domain.EMailStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckEMail indicates an expected call of CheckEMail.
func (mr *MockRegServerMockRecorder) CheckEMail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckEMail", reflect.TypeOf((*MockRegServer)(nil).CheckEMail), arg0, arg1)
}

// InitMasterKey mocks base method.
func (m *MockRegServer) InitMasterKey(arg0 context.Context, arg1 *domain.MasterKeyData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InitMasterKey", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InitMasterKey indicates an expected call of InitMasterKey.
func (mr *MockRegServerMockRecorder) InitMasterKey(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitMasterKey", reflect.TypeOf((*MockRegServer)(nil).InitMasterKey), arg0, arg1)
}

// PassRegOTP mocks base method.
func (m *MockRegServer) PassRegOTP(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PassRegOTP", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PassRegOTP indicates an expected call of PassRegOTP.
func (mr *MockRegServerMockRecorder) PassRegOTP(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PassRegOTP", reflect.TypeOf((*MockRegServer)(nil).PassRegOTP), arg0, arg1)
}

// Registrate mocks base method.
func (m *MockRegServer) Registrate(arg0 context.Context, arg1 *domain.EMailData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Registrate", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Registrate indicates an expected call of Registrate.
func (mr *MockRegServerMockRecorder) Registrate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Registrate", reflect.TypeOf((*MockRegServer)(nil).Registrate), arg0, arg1)
}

// MockRegView is a mock of RegView interface.
type MockRegView struct {
	ctrl     *gomock.Controller
	recorder *MockRegViewMockRecorder
}

// MockRegViewMockRecorder is the mock recorder for MockRegView.
type MockRegViewMockRecorder struct {
	mock *MockRegView
}

// NewMockRegView creates a new mock instance.
func NewMockRegView(ctrl *gomock.Controller) *MockRegView {
	mock := &MockRegView{ctrl: ctrl}
	mock.recorder = &MockRegViewMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRegView) EXPECT() *MockRegViewMockRecorder {
	return m.recorder
}

// ShowError mocks base method.
func (m *MockRegView) ShowError(arg0 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowError", arg0)
}

// ShowError indicates an expected call of ShowError.
func (mr *MockRegViewMockRecorder) ShowError(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowError", reflect.TypeOf((*MockRegView)(nil).ShowError), arg0)
}

// ShowLoginView mocks base method.
func (m *MockRegView) ShowLoginView() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowLoginView")
}

// ShowLoginView indicates an expected call of ShowLoginView.
func (mr *MockRegViewMockRecorder) ShowLoginView() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowLoginView", reflect.TypeOf((*MockRegView)(nil).ShowLoginView))
}

// ShowMsg mocks base method.
func (m *MockRegView) ShowMsg(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowMsg", arg0)
}

// ShowMsg indicates an expected call of ShowMsg.
func (mr *MockRegViewMockRecorder) ShowMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowMsg", reflect.TypeOf((*MockRegView)(nil).ShowMsg), arg0)
}

// ShowRegMasterKeyView mocks base method.
func (m *MockRegView) ShowRegMasterKeyView() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowRegMasterKeyView")
}

// ShowRegMasterKeyView indicates an expected call of ShowRegMasterKeyView.
func (mr *MockRegViewMockRecorder) ShowRegMasterKeyView() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowRegMasterKeyView", reflect.TypeOf((*MockRegView)(nil).ShowRegMasterKeyView))
}

// ShowRegOTPView mocks base method.
func (m *MockRegView) ShowRegOTPView() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowRegOTPView")
}

// ShowRegOTPView indicates an expected call of ShowRegOTPView.
func (mr *MockRegViewMockRecorder) ShowRegOTPView() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowRegOTPView", reflect.TypeOf((*MockRegView)(nil).ShowRegOTPView))
}

// ShowRegView mocks base method.
func (m *MockRegView) ShowRegView() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowRegView")
}

// ShowRegView indicates an expected call of ShowRegView.
func (mr *MockRegViewMockRecorder) ShowRegView() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowRegView", reflect.TypeOf((*MockRegView)(nil).ShowRegView))
}

// MockRegHelper is a mock of RegHelper interface.
type MockRegHelper struct {
	ctrl     *gomock.Controller
	recorder *MockRegHelperMockRecorder
}

// MockRegHelperMockRecorder is the mock recorder for MockRegHelper.
type MockRegHelperMockRecorder struct {
	mock *MockRegHelper
}

// NewMockRegHelper creates a new mock instance.
func NewMockRegHelper(ctrl *gomock.Controller) *MockRegHelper {
	mock := &MockRegHelper{ctrl: ctrl}
	mock.recorder = &MockRegHelperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRegHelper) EXPECT() *MockRegHelperMockRecorder {
	return m.recorder
}

// CheckAuthPasswordComplexityLevel mocks base method.
func (m *MockRegHelper) CheckAuthPasswordComplexityLevel(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAuthPasswordComplexityLevel", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CheckAuthPasswordComplexityLevel indicates an expected call of CheckAuthPasswordComplexityLevel.
func (mr *MockRegHelperMockRecorder) CheckAuthPasswordComplexityLevel(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAuthPasswordComplexityLevel", reflect.TypeOf((*MockRegHelper)(nil).CheckAuthPasswordComplexityLevel), arg0)
}

// CheckMasterPasswordComplexityLevel mocks base method.
func (m *MockRegHelper) CheckMasterPasswordComplexityLevel(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckMasterPasswordComplexityLevel", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CheckMasterPasswordComplexityLevel indicates an expected call of CheckMasterPasswordComplexityLevel.
func (mr *MockRegHelperMockRecorder) CheckMasterPasswordComplexityLevel(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckMasterPasswordComplexityLevel", reflect.TypeOf((*MockRegHelper)(nil).CheckMasterPasswordComplexityLevel), arg0)
}

// EncryptHello mocks base method.
func (m *MockRegHelper) EncryptHello(arg0, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EncryptHello", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EncryptHello indicates an expected call of EncryptHello.
func (mr *MockRegHelperMockRecorder) EncryptHello(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EncryptHello", reflect.TypeOf((*MockRegHelper)(nil).EncryptHello), arg0, arg1)
}

// ParseEMail mocks base method.
func (m *MockRegHelper) ParseEMail(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseEMail", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// ParseEMail indicates an expected call of ParseEMail.
func (mr *MockRegHelperMockRecorder) ParseEMail(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseEMail", reflect.TypeOf((*MockRegHelper)(nil).ParseEMail), arg0)
}

// Random32ByteString mocks base method.
func (m *MockRegHelper) Random32ByteString() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Random32ByteString")
	ret0, _ := ret[0].(string)
	return ret0
}

// Random32ByteString indicates an expected call of Random32ByteString.
func (mr *MockRegHelperMockRecorder) Random32ByteString() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Random32ByteString", reflect.TypeOf((*MockRegHelper)(nil).Random32ByteString))
}

// MockLoginServer is a mock of LoginServer interface.
type MockLoginServer struct {
	ctrl     *gomock.Controller
	recorder *MockLoginServerMockRecorder
}

// MockLoginServerMockRecorder is the mock recorder for MockLoginServer.
type MockLoginServerMockRecorder struct {
	mock *MockLoginServer
}

// NewMockLoginServer creates a new mock instance.
func NewMockLoginServer(ctrl *gomock.Controller) *MockLoginServer {
	mock := &MockLoginServer{ctrl: ctrl}
	mock.recorder = &MockLoginServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLoginServer) EXPECT() *MockLoginServerMockRecorder {
	return m.recorder
}

// GetHelloData mocks base method.
func (m *MockLoginServer) GetHelloData(arg0 context.Context) (*domain.HelloData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHelloData", arg0)
	ret0, _ := ret[0].(*domain.HelloData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHelloData indicates an expected call of GetHelloData.
func (mr *MockLoginServerMockRecorder) GetHelloData(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHelloData", reflect.TypeOf((*MockLoginServer)(nil).GetHelloData), arg0)
}

// Login mocks base method.
func (m *MockLoginServer) Login(arg0 context.Context, arg1 *domain.EMailData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Login indicates an expected call of Login.
func (mr *MockLoginServerMockRecorder) Login(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockLoginServer)(nil).Login), arg0, arg1)
}

// PassLoginOTP mocks base method.
func (m *MockLoginServer) PassLoginOTP(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PassLoginOTP", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PassLoginOTP indicates an expected call of PassLoginOTP.
func (mr *MockLoginServerMockRecorder) PassLoginOTP(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PassLoginOTP", reflect.TypeOf((*MockLoginServer)(nil).PassLoginOTP), arg0, arg1)
}

// MockLoginView is a mock of LoginView interface.
type MockLoginView struct {
	ctrl     *gomock.Controller
	recorder *MockLoginViewMockRecorder
}

// MockLoginViewMockRecorder is the mock recorder for MockLoginView.
type MockLoginViewMockRecorder struct {
	mock *MockLoginView
}

// NewMockLoginView creates a new mock instance.
func NewMockLoginView(ctrl *gomock.Controller) *MockLoginView {
	mock := &MockLoginView{ctrl: ctrl}
	mock.recorder = &MockLoginViewMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLoginView) EXPECT() *MockLoginViewMockRecorder {
	return m.recorder
}

// ShowDataAccessView mocks base method.
func (m *MockLoginView) ShowDataAccessView() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowDataAccessView")
}

// ShowDataAccessView indicates an expected call of ShowDataAccessView.
func (mr *MockLoginViewMockRecorder) ShowDataAccessView() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowDataAccessView", reflect.TypeOf((*MockLoginView)(nil).ShowDataAccessView))
}

// ShowError mocks base method.
func (m *MockLoginView) ShowError(arg0 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowError", arg0)
}

// ShowError indicates an expected call of ShowError.
func (mr *MockLoginViewMockRecorder) ShowError(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowError", reflect.TypeOf((*MockLoginView)(nil).ShowError), arg0)
}

// ShowLogOTPView mocks base method.
func (m *MockLoginView) ShowLogOTPView() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowLogOTPView")
}

// ShowLogOTPView indicates an expected call of ShowLogOTPView.
func (mr *MockLoginViewMockRecorder) ShowLogOTPView() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowLogOTPView", reflect.TypeOf((*MockLoginView)(nil).ShowLogOTPView))
}

// ShowMasterKeyView mocks base method.
func (m *MockLoginView) ShowMasterKeyView(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowMasterKeyView", arg0)
}

// ShowMasterKeyView indicates an expected call of ShowMasterKeyView.
func (mr *MockLoginViewMockRecorder) ShowMasterKeyView(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowMasterKeyView", reflect.TypeOf((*MockLoginView)(nil).ShowMasterKeyView), arg0)
}

// ShowMsg mocks base method.
func (m *MockLoginView) ShowMsg(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowMsg", arg0)
}

// ShowMsg indicates an expected call of ShowMsg.
func (mr *MockLoginViewMockRecorder) ShowMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowMsg", reflect.TypeOf((*MockLoginView)(nil).ShowMsg), arg0)
}

// MockLoginHelper is a mock of LoginHelper interface.
type MockLoginHelper struct {
	ctrl     *gomock.Controller
	recorder *MockLoginHelperMockRecorder
}

// MockLoginHelperMockRecorder is the mock recorder for MockLoginHelper.
type MockLoginHelperMockRecorder struct {
	mock *MockLoginHelper
}

// NewMockLoginHelper creates a new mock instance.
func NewMockLoginHelper(ctrl *gomock.Controller) *MockLoginHelper {
	mock := &MockLoginHelper{ctrl: ctrl}
	mock.recorder = &MockLoginHelperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLoginHelper) EXPECT() *MockLoginHelperMockRecorder {
	return m.recorder
}

// DecryptHello mocks base method.
func (m *MockLoginHelper) DecryptHello(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecryptHello", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DecryptHello indicates an expected call of DecryptHello.
func (mr *MockLoginHelperMockRecorder) DecryptHello(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecryptHello", reflect.TypeOf((*MockLoginHelper)(nil).DecryptHello), arg0, arg1)
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
