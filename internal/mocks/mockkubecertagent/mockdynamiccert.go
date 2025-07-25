// Copyright 2020-2025 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

// Code generated by MockGen. DO NOT EDIT.
// Source: go.pinniped.dev/internal/dynamiccert (interfaces: Private)
//
// Generated by this command:
//
//	mockgen -destination=mockdynamiccert.go -package=mocks -copyright_file=../../../hack/header.txt -mock_names Private=MockDynamicCertPrivate go.pinniped.dev/internal/dynamiccert Private
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	dynamiccertificates "k8s.io/apiserver/pkg/server/dynamiccertificates"
)

// MockDynamicCertPrivate is a mock of Private interface.
type MockDynamicCertPrivate struct {
	ctrl     *gomock.Controller
	recorder *MockDynamicCertPrivateMockRecorder
	isgomock struct{}
}

// MockDynamicCertPrivateMockRecorder is the mock recorder for MockDynamicCertPrivate.
type MockDynamicCertPrivateMockRecorder struct {
	mock *MockDynamicCertPrivate
}

// NewMockDynamicCertPrivate creates a new mock instance.
func NewMockDynamicCertPrivate(ctrl *gomock.Controller) *MockDynamicCertPrivate {
	mock := &MockDynamicCertPrivate{ctrl: ctrl}
	mock.recorder = &MockDynamicCertPrivateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDynamicCertPrivate) EXPECT() *MockDynamicCertPrivateMockRecorder {
	return m.recorder
}

// AddListener mocks base method.
func (m *MockDynamicCertPrivate) AddListener(listener dynamiccertificates.Listener) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddListener", listener)
}

// AddListener indicates an expected call of AddListener.
func (mr *MockDynamicCertPrivateMockRecorder) AddListener(listener any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddListener", reflect.TypeOf((*MockDynamicCertPrivate)(nil).AddListener), listener)
}

// CurrentCertKeyContent mocks base method.
func (m *MockDynamicCertPrivate) CurrentCertKeyContent() ([]byte, []byte) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentCertKeyContent")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].([]byte)
	return ret0, ret1
}

// CurrentCertKeyContent indicates an expected call of CurrentCertKeyContent.
func (mr *MockDynamicCertPrivateMockRecorder) CurrentCertKeyContent() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentCertKeyContent", reflect.TypeOf((*MockDynamicCertPrivate)(nil).CurrentCertKeyContent))
}

// Name mocks base method.
func (m *MockDynamicCertPrivate) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockDynamicCertPrivateMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockDynamicCertPrivate)(nil).Name))
}

// Run mocks base method.
func (m *MockDynamicCertPrivate) Run(ctx context.Context, workers int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run", ctx, workers)
}

// Run indicates an expected call of Run.
func (mr *MockDynamicCertPrivateMockRecorder) Run(ctx, workers any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockDynamicCertPrivate)(nil).Run), ctx, workers)
}

// RunOnce mocks base method.
func (m *MockDynamicCertPrivate) RunOnce(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunOnce", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunOnce indicates an expected call of RunOnce.
func (mr *MockDynamicCertPrivateMockRecorder) RunOnce(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunOnce", reflect.TypeOf((*MockDynamicCertPrivate)(nil).RunOnce), ctx)
}

// SetCertKeyContent mocks base method.
func (m *MockDynamicCertPrivate) SetCertKeyContent(certPEM, keyPEM []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetCertKeyContent", certPEM, keyPEM)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetCertKeyContent indicates an expected call of SetCertKeyContent.
func (mr *MockDynamicCertPrivateMockRecorder) SetCertKeyContent(certPEM, keyPEM any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCertKeyContent", reflect.TypeOf((*MockDynamicCertPrivate)(nil).SetCertKeyContent), certPEM, keyPEM)
}

// UnsetCertKeyContent mocks base method.
func (m *MockDynamicCertPrivate) UnsetCertKeyContent() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UnsetCertKeyContent")
}

// UnsetCertKeyContent indicates an expected call of UnsetCertKeyContent.
func (mr *MockDynamicCertPrivateMockRecorder) UnsetCertKeyContent() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnsetCertKeyContent", reflect.TypeOf((*MockDynamicCertPrivate)(nil).UnsetCertKeyContent))
}
