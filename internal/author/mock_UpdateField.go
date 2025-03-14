// Code generated by mockery v2.53.2. DO NOT EDIT.

package author

import mock "github.com/stretchr/testify/mock"

// MockUpdateField is an autogenerated mock type for the UpdateField type
type MockUpdateField struct {
	mock.Mock
}

type MockUpdateField_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUpdateField) EXPECT() *MockUpdateField_Expecter {
	return &MockUpdateField_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: _a0
func (_m *MockUpdateField) Execute(_a0 map[string]interface{}) {
	_m.Called(_a0)
}

// MockUpdateField_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockUpdateField_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - _a0 map[string]interface{}
func (_e *MockUpdateField_Expecter) Execute(_a0 interface{}) *MockUpdateField_Execute_Call {
	return &MockUpdateField_Execute_Call{Call: _e.mock.On("Execute", _a0)}
}

func (_c *MockUpdateField_Execute_Call) Run(run func(_a0 map[string]interface{})) *MockUpdateField_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(map[string]interface{}))
	})
	return _c
}

func (_c *MockUpdateField_Execute_Call) Return() *MockUpdateField_Execute_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockUpdateField_Execute_Call) RunAndReturn(run func(map[string]interface{})) *MockUpdateField_Execute_Call {
	_c.Run(run)
	return _c
}

// NewMockUpdateField creates a new instance of MockUpdateField. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUpdateField(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUpdateField {
	mock := &MockUpdateField{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
