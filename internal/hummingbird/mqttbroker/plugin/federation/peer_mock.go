// Code generated by MockGen. DO NOT EDIT.
// Source: plugin/federation/peer.go

// Package federation is a generated GoMock package.
package federation

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// Mockqueue is a mock of queue interface
type Mockqueue struct {
	ctrl     *gomock.Controller
	recorder *MockqueueMockRecorder
}

// MockqueueMockRecorder is the mock recorder for Mockqueue
type MockqueueMockRecorder struct {
	mock *Mockqueue
}

// NewMockqueue creates a new mock instance
func NewMockqueue(ctrl *gomock.Controller) *Mockqueue {
	mock := &Mockqueue{ctrl: ctrl}
	mock.recorder = &MockqueueMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *Mockqueue) EXPECT() *MockqueueMockRecorder {
	return m.recorder
}

// clear mocks base method
func (m *Mockqueue) clear() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "clear")
}

// clear indicates an expected call of clear
func (mr *MockqueueMockRecorder) clear() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "clear", reflect.TypeOf((*Mockqueue)(nil).clear))
}

// close mocks base method
func (m *Mockqueue) close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "close")
}

// close indicates an expected call of close
func (mr *MockqueueMockRecorder) close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "close", reflect.TypeOf((*Mockqueue)(nil).close))
}

// open mocks base method
func (m *Mockqueue) open() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "open")
}

// open indicates an expected call of open
func (mr *MockqueueMockRecorder) open() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "open", reflect.TypeOf((*Mockqueue)(nil).open))
}

// setReadPosition mocks base method
func (m *Mockqueue) setReadPosition(id uint64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "setReadPosition", id)
}

// setReadPosition indicates an expected call of setReadPosition
func (mr *MockqueueMockRecorder) setReadPosition(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "setReadPosition", reflect.TypeOf((*Mockqueue)(nil).setReadPosition), id)
}

// add mocks base method
func (m *Mockqueue) add(event *Event) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "add", event)
}

// add indicates an expected call of add
func (mr *MockqueueMockRecorder) add(event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "add", reflect.TypeOf((*Mockqueue)(nil).add), event)
}

// fetchEvents mocks base method
func (m *Mockqueue) fetchEvents() []*Event {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "fetchEvents")
	ret0, _ := ret[0].([]*Event)
	return ret0
}

// fetchEvents indicates an expected call of fetchEvents
func (mr *MockqueueMockRecorder) fetchEvents() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "fetchEvents", reflect.TypeOf((*Mockqueue)(nil).fetchEvents))
}

// ack mocks base method
func (m *Mockqueue) ack(id uint64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ack", id)
}

// ack indicates an expected call of ack
func (mr *MockqueueMockRecorder) ack(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ack", reflect.TypeOf((*Mockqueue)(nil).ack), id)
}