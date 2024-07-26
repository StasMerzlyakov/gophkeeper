// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/StasMerzlyakov/gophkeeper/internal/client/app (interfaces: RegServer,RegView,RegHelper,LoginServer,LoginView,LoginHelper)

// Package app_test is a generated GoMock package.
package app_test

import (
	context "context"
	reflect "reflect"

	domain "github.com/StasMerzlyakov/gophkeeper/internal/domain"
	gomock "github.com/golang/mock/gomock"
)

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

// PassOTP mocks base method.
func (m *MockRegServer) PassOTP(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PassOTP", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PassOTP indicates an expected call of PassOTP.
func (mr *MockRegServerMockRecorder) PassOTP(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PassOTP", reflect.TypeOf((*MockRegServer)(nil).PassOTP), arg0, arg1)
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

// ShowInitMasterKeyView mocks base method.
func (m *MockRegView) ShowInitMasterKeyView() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowInitMasterKeyView")
}

// ShowInitMasterKeyView indicates an expected call of ShowInitMasterKeyView.
func (mr *MockRegViewMockRecorder) ShowInitMasterKeyView() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowInitMasterKeyView", reflect.TypeOf((*MockRegView)(nil).ShowInitMasterKeyView))
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

// ShowRegForm mocks base method.
func (m *MockRegView) ShowRegForm() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShowRegForm")
}

// ShowRegForm indicates an expected call of ShowRegForm.
func (mr *MockRegViewMockRecorder) ShowRegForm() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowRegForm", reflect.TypeOf((*MockRegView)(nil).ShowRegForm))
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

// CheckMasterKeyPasswordComplexityLevel mocks base method.
func (m *MockRegHelper) CheckMasterKeyPasswordComplexityLevel(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckMasterKeyPasswordComplexityLevel", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CheckMasterKeyPasswordComplexityLevel indicates an expected call of CheckMasterKeyPasswordComplexityLevel.
func (mr *MockRegHelperMockRecorder) CheckMasterKeyPasswordComplexityLevel(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckMasterKeyPasswordComplexityLevel", reflect.TypeOf((*MockRegHelper)(nil).CheckMasterKeyPasswordComplexityLevel), arg0)
}

// EncryptMasterKey mocks base method.
func (m *MockRegHelper) EncryptMasterKey(arg0, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EncryptMasterKey", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EncryptMasterKey indicates an expected call of EncryptMasterKey.
func (mr *MockRegHelperMockRecorder) EncryptMasterKey(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EncryptMasterKey", reflect.TypeOf((*MockRegHelper)(nil).EncryptMasterKey), arg0, arg1)
}

// EncryptShortData mocks base method.
func (m *MockRegHelper) EncryptShortData(arg0 []byte, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EncryptShortData", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EncryptShortData indicates an expected call of EncryptShortData.
func (mr *MockRegHelperMockRecorder) EncryptShortData(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EncryptShortData", reflect.TypeOf((*MockRegHelper)(nil).EncryptShortData), arg0, arg1)
}

// GenerateHello mocks base method.
func (m *MockRegHelper) GenerateHello() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateHello")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateHello indicates an expected call of GenerateHello.
func (mr *MockRegHelperMockRecorder) GenerateHello() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateHello", reflect.TypeOf((*MockRegHelper)(nil).GenerateHello))
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

// GetMasterKey mocks base method.
func (m *MockLoginServer) GetMasterKey(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetMasterKey", arg0)
}

// GetMasterKey indicates an expected call of GetMasterKey.
func (mr *MockLoginServerMockRecorder) GetMasterKey(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMasterKey", reflect.TypeOf((*MockLoginServer)(nil).GetMasterKey), arg0)
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

// PassOTP mocks base method.
func (m *MockLoginServer) PassOTP(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PassOTP", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PassOTP indicates an expected call of PassOTP.
func (mr *MockLoginServerMockRecorder) PassOTP(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PassOTP", reflect.TypeOf((*MockLoginServer)(nil).PassOTP), arg0, arg1)
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

// CheckHello mocks base method.
func (m *MockLoginHelper) CheckHello(arg0 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckHello", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckHello indicates an expected call of CheckHello.
func (mr *MockLoginHelperMockRecorder) CheckHello(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckHello", reflect.TypeOf((*MockLoginHelper)(nil).CheckHello), arg0)
}

// DecryptMasterKey mocks base method.
func (m *MockLoginHelper) DecryptMasterKey(arg0, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecryptMasterKey", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DecryptMasterKey indicates an expected call of DecryptMasterKey.
func (mr *MockLoginHelperMockRecorder) DecryptMasterKey(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecryptMasterKey", reflect.TypeOf((*MockLoginHelper)(nil).DecryptMasterKey), arg0, arg1)
}

// DecryptShortData mocks base method.
func (m *MockLoginHelper) DecryptShortData(arg0, arg1 string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecryptShortData", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DecryptShortData indicates an expected call of DecryptShortData.
func (mr *MockLoginHelperMockRecorder) DecryptShortData(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecryptShortData", reflect.TypeOf((*MockLoginHelper)(nil).DecryptShortData), arg0, arg1)
}
