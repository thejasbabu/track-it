// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	gocui "github.com/jroimartin/gocui"
	display "github.com/thejasbabu/track-it/pkg/task/display"

	mock "github.com/stretchr/testify/mock"
)

// Pane is an autogenerated mock type for the Pane type
type Pane struct {
	mock.Mock
}

// KeyDown provides a mock function with given fields: g, v
func (_m *Pane) KeyDown(g *gocui.Gui, v *gocui.View) (interface{}, error) {
	ret := _m.Called(g, v)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(*gocui.Gui, *gocui.View) interface{}); ok {
		r0 = rf(g, v)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gocui.Gui, *gocui.View) error); ok {
		r1 = rf(g, v)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// KeyUp provides a mock function with given fields: g, v
func (_m *Pane) KeyUp(g *gocui.Gui, v *gocui.View) (interface{}, error) {
	ret := _m.Called(g, v)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(*gocui.Gui, *gocui.View) interface{}); ok {
		r0 = rf(g, v)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gocui.Gui, *gocui.View) error); ok {
		r1 = rf(g, v)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Layout provides a mock function with given fields: gui, pos
func (_m *Pane) Layout(gui *gocui.Gui, pos display.Position) error {
	ret := _m.Called(gui, pos)

	var r0 error
	if rf, ok := ret.Get(0).(func(*gocui.Gui, display.Position) error); ok {
		r0 = rf(gui, pos)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: g, op, data
func (_m *Pane) Update(g *gocui.Gui, op display.Operation, data interface{}) (interface{}, error) {
	ret := _m.Called(g, op, data)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(*gocui.Gui, display.Operation, interface{}) interface{}); ok {
		r0 = rf(g, op, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*gocui.Gui, display.Operation, interface{}) error); ok {
		r1 = rf(g, op, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
