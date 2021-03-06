// Code generated by mockery v2.2.1. DO NOT EDIT.

package mocks

import (
	jwt "github.com/dgrijalva/jwt-go"
	mock "github.com/stretchr/testify/mock"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// Issue provides a mock function with given fields: claims
func (_m *Service) Issue(claims jwt.Claims) *jwt.Token {
	ret := _m.Called(claims)

	var r0 *jwt.Token
	if rf, ok := ret.Get(0).(func(jwt.Claims) *jwt.Token); ok {
		r0 = rf(claims)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jwt.Token)
		}
	}

	return r0
}

// Parse provides a mock function with given fields: s, secret
func (_m *Service) Parse(s string, secret ...[]byte) (*jwt.Token, error) {
	_va := make([]interface{}, len(secret))
	for _i := range secret {
		_va[_i] = secret[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, s)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *jwt.Token
	if rf, ok := ret.Get(0).(func(string, ...[]byte) *jwt.Token); ok {
		r0 = rf(s, secret...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*jwt.Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, ...[]byte) error); ok {
		r1 = rf(s, secret...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Sign provides a mock function with given fields: t, secret
func (_m *Service) Sign(t *jwt.Token, secret ...[]byte) (string, error) {
	_va := make([]interface{}, len(secret))
	for _i := range secret {
		_va[_i] = secret[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, t)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 string
	if rf, ok := ret.Get(0).(func(*jwt.Token, ...[]byte) string); ok {
		r0 = rf(t, secret...)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*jwt.Token, ...[]byte) error); ok {
		r1 = rf(t, secret...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
