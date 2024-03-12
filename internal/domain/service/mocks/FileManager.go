// Code generated by mockery v2.36.1. DO NOT EDIT.

package mocks

import (
	model "go-clickhouse-migrator/internal/domain/model"

	mock "github.com/stretchr/testify/mock"
)

// FileManager is an autogenerated mock type for the FileManager type
type FileManager struct {
	mock.Mock
}

// Create provides a mock function with given fields: name
func (_m *FileManager) Create(name string) (string, error) {
	ret := _m.Called(name)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Read provides a mock function with given fields: data
func (_m *FileManager) Read(data model.MigrationFileInfo) (string, error) {
	ret := _m.Called(data)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(model.MigrationFileInfo) (string, error)); ok {
		return rf(data)
	}
	if rf, ok := ret.Get(0).(func(model.MigrationFileInfo) string); ok {
		r0 = rf(data)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(model.MigrationFileInfo) error); ok {
		r1 = rf(data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadHelp provides a mock function with given fields: file
func (_m *FileManager) ReadHelp(file string) (string, error) {
	ret := _m.Called(file)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(file)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(file)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(file)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SortedMigrationFilesData provides a mock function with given fields:
func (_m *FileManager) SortedMigrationFilesData() ([]model.MigrationFileInfo, error) {
	ret := _m.Called()

	var r0 []model.MigrationFileInfo
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]model.MigrationFileInfo, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []model.MigrationFileInfo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.MigrationFileInfo)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewFileManager creates a new instance of FileManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFileManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *FileManager {
	mock := &FileManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}