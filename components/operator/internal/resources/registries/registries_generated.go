// Code generated by MockGen. DO NOT EDIT.
// Source: ./registries.go
//
// Generated by this command:
//
//	mockgen -source ./registries.go -destination ./registries_generated.go -package registries . ImageSettingsOverrider
//

// Package registries is a generated GoMock package.
package registries

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockImageSettingsOverrider is a mock of ImageSettingsOverrider interface.
type MockImageSettingsOverrider struct {
	ctrl     *gomock.Controller
	recorder *MockImageSettingsOverriderMockRecorder
}

// MockImageSettingsOverriderMockRecorder is the mock recorder for MockImageSettingsOverrider.
type MockImageSettingsOverriderMockRecorder struct {
	mock *MockImageSettingsOverrider
}

// NewMockImageSettingsOverrider creates a new mock instance.
func NewMockImageSettingsOverrider(ctrl *gomock.Controller) *MockImageSettingsOverrider {
	mock := &MockImageSettingsOverrider{ctrl: ctrl}
	mock.recorder = &MockImageSettingsOverriderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockImageSettingsOverrider) EXPECT() *MockImageSettingsOverriderMockRecorder {
	return m.recorder
}

// OverrideWithSetting mocks base method.
func (m *MockImageSettingsOverrider) OverrideWithSetting(arg0 *imageOrigin, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OverrideWithSetting", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// OverrideWithSetting indicates an expected call of OverrideWithSetting.
func (mr *MockImageSettingsOverriderMockRecorder) OverrideWithSetting(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OverrideWithSetting", reflect.TypeOf((*MockImageSettingsOverrider)(nil).OverrideWithSetting), arg0, arg1)
}
