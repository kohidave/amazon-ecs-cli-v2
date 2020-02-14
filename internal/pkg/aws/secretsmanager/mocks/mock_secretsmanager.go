// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/pkg/aws/secretsmanager/secretsmanager.go

// Package mocks is a generated GoMock package.
package mocks

import (
	secretsmanager "github.com/aws/aws-sdk-go/service/secretsmanager"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockSecretsManagerAPI is a mock of SecretsManagerAPI interface
type MockSecretsManagerAPI struct {
	ctrl     *gomock.Controller
	recorder *MockSecretsManagerAPIMockRecorder
}

// MockSecretsManagerAPIMockRecorder is the mock recorder for MockSecretsManagerAPI
type MockSecretsManagerAPIMockRecorder struct {
	mock *MockSecretsManagerAPI
}

// NewMockSecretsManagerAPI creates a new mock instance
func NewMockSecretsManagerAPI(ctrl *gomock.Controller) *MockSecretsManagerAPI {
	mock := &MockSecretsManagerAPI{ctrl: ctrl}
	mock.recorder = &MockSecretsManagerAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSecretsManagerAPI) EXPECT() *MockSecretsManagerAPIMockRecorder {
	return m.recorder
}

// CreateSecret mocks base method
func (m *MockSecretsManagerAPI) CreateSecret(arg0 *secretsmanager.CreateSecretInput) (*secretsmanager.CreateSecretOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSecret", arg0)
	ret0, _ := ret[0].(*secretsmanager.CreateSecretOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSecret indicates an expected call of CreateSecret
func (mr *MockSecretsManagerAPIMockRecorder) CreateSecret(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSecret", reflect.TypeOf((*MockSecretsManagerAPI)(nil).CreateSecret), arg0)
}
