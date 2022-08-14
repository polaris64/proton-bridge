// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ProtonMail/proton-bridge/v2/internal/transfer (interfaces: PanicHandler,IMAPClientProvider)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	imap "github.com/emersion/go-imap"
	sasl "github.com/emersion/go-sasl"
	gomock "github.com/golang/mock/gomock"
)

// MockPanicHandler is a mock of PanicHandler interface.
type MockPanicHandler struct {
	ctrl     *gomock.Controller
	recorder *MockPanicHandlerMockRecorder
}

// MockPanicHandlerMockRecorder is the mock recorder for MockPanicHandler.
type MockPanicHandlerMockRecorder struct {
	mock *MockPanicHandler
}

// NewMockPanicHandler creates a new mock instance.
func NewMockPanicHandler(ctrl *gomock.Controller) *MockPanicHandler {
	mock := &MockPanicHandler{ctrl: ctrl}
	mock.recorder = &MockPanicHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPanicHandler) EXPECT() *MockPanicHandlerMockRecorder {
	return m.recorder
}

// HandlePanic mocks base method.
func (m *MockPanicHandler) HandlePanic() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandlePanic")
}

// HandlePanic indicates an expected call of HandlePanic.
func (mr *MockPanicHandlerMockRecorder) HandlePanic() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandlePanic", reflect.TypeOf((*MockPanicHandler)(nil).HandlePanic))
}

// MockIMAPClientProvider is a mock of IMAPClientProvider interface.
type MockIMAPClientProvider struct {
	ctrl     *gomock.Controller
	recorder *MockIMAPClientProviderMockRecorder
}

// MockIMAPClientProviderMockRecorder is the mock recorder for MockIMAPClientProvider.
type MockIMAPClientProviderMockRecorder struct {
	mock *MockIMAPClientProvider
}

// NewMockIMAPClientProvider creates a new mock instance.
func NewMockIMAPClientProvider(ctrl *gomock.Controller) *MockIMAPClientProvider {
	mock := &MockIMAPClientProvider{ctrl: ctrl}
	mock.recorder = &MockIMAPClientProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIMAPClientProvider) EXPECT() *MockIMAPClientProviderMockRecorder {
	return m.recorder
}

// Authenticate mocks base method.
func (m *MockIMAPClientProvider) Authenticate(arg0 sasl.Client) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authenticate", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Authenticate indicates an expected call of Authenticate.
func (mr *MockIMAPClientProviderMockRecorder) Authenticate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authenticate", reflect.TypeOf((*MockIMAPClientProvider)(nil).Authenticate), arg0)
}

// Capability mocks base method.
func (m *MockIMAPClientProvider) Capability() (map[string]bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Capability")
	ret0, _ := ret[0].(map[string]bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Capability indicates an expected call of Capability.
func (mr *MockIMAPClientProviderMockRecorder) Capability() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Capability", reflect.TypeOf((*MockIMAPClientProvider)(nil).Capability))
}

// Fetch mocks base method.
func (m *MockIMAPClientProvider) Fetch(arg0 *imap.SeqSet, arg1 []imap.FetchItem, arg2 chan *imap.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fetch", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Fetch indicates an expected call of Fetch.
func (mr *MockIMAPClientProviderMockRecorder) Fetch(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockIMAPClientProvider)(nil).Fetch), arg0, arg1, arg2)
}

// List mocks base method.
func (m *MockIMAPClientProvider) List(arg0, arg1 string, arg2 chan *imap.MailboxInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// List indicates an expected call of List.
func (mr *MockIMAPClientProviderMockRecorder) List(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockIMAPClientProvider)(nil).List), arg0, arg1, arg2)
}

// Login mocks base method.
func (m *MockIMAPClientProvider) Login(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Login indicates an expected call of Login.
func (mr *MockIMAPClientProviderMockRecorder) Login(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockIMAPClientProvider)(nil).Login), arg0, arg1)
}

// Select mocks base method.
func (m *MockIMAPClientProvider) Select(arg0 string, arg1 bool) (*imap.MailboxStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Select", arg0, arg1)
	ret0, _ := ret[0].(*imap.MailboxStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Select indicates an expected call of Select.
func (mr *MockIMAPClientProviderMockRecorder) Select(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Select", reflect.TypeOf((*MockIMAPClientProvider)(nil).Select), arg0, arg1)
}

// State mocks base method.
func (m *MockIMAPClientProvider) State() imap.ConnState {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "State")
	ret0, _ := ret[0].(imap.ConnState)
	return ret0
}

// State indicates an expected call of State.
func (mr *MockIMAPClientProviderMockRecorder) State() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "State", reflect.TypeOf((*MockIMAPClientProvider)(nil).State))
}

// Support mocks base method.
func (m *MockIMAPClientProvider) Support(arg0 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Support", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Support indicates an expected call of Support.
func (mr *MockIMAPClientProviderMockRecorder) Support(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Support", reflect.TypeOf((*MockIMAPClientProvider)(nil).Support), arg0)
}

// SupportAuth mocks base method.
func (m *MockIMAPClientProvider) SupportAuth(arg0 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SupportAuth", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SupportAuth indicates an expected call of SupportAuth.
func (mr *MockIMAPClientProviderMockRecorder) SupportAuth(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SupportAuth", reflect.TypeOf((*MockIMAPClientProvider)(nil).SupportAuth), arg0)
}

// UidFetch mocks base method.
func (m *MockIMAPClientProvider) UidFetch(arg0 *imap.SeqSet, arg1 []imap.FetchItem, arg2 chan *imap.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UidFetch", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UidFetch indicates an expected call of UidFetch.
func (mr *MockIMAPClientProviderMockRecorder) UidFetch(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UidFetch", reflect.TypeOf((*MockIMAPClientProvider)(nil).UidFetch), arg0, arg1, arg2)
}
