// Code generated by mockery v1.0.0. DO NOT EDIT.
package mocks

import context "context"
import menekel "github.com/golangid/menekel"
import mock "github.com/stretchr/testify/mock"

// AuthorRepository is an autogenerated mock type for the AuthorRepository type
type AuthorRepository struct {
	mock.Mock
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *AuthorRepository) GetByID(ctx context.Context, id int64) (*menekel.Author, error) {
	ret := _m.Called(ctx, id)

	var r0 *menekel.Author
	if rf, ok := ret.Get(0).(func(context.Context, int64) *menekel.Author); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*menekel.Author)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
