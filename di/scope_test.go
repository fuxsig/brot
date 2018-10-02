// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package di

import (
	"math"
	"reflect"
	"testing"
)

func TestAssignError(t *testing.T) {
	s := NewScope()
	v := true
	err := s.assignValue(reflect.ValueOf(&v), false)
	if err != UnaddressableError {
		t.Error("expected UnaddressableError")
	}
}

func TestAssignBool(t *testing.T) {
	tables := []struct {
		i bool
		v interface{}
		e error
		r bool
	}{
		{false, true, nil, true},
		{false, "True", nil, true},
		{false, 1, nil, true},
		{true, false, nil, false},
		{true, "False", nil, false},
		{true, 0, nil, false},
	}
	s := NewScope()

	for _, table := range tables {
		v := table.i
		ptr := reflect.ValueOf(&v)
		err := s.assignValue(ptr.Elem(), table.v)
		if err == table.e {
			if table.r != v {
				t.Errorf("expected %t, got %t", table.r, v)
			}
		} else {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestAssignFloat32(t *testing.T) {
	tables := []struct {
		i float32
		v interface{}
		e error
		r float32
	}{
		{0, 1.2345, nil, 1.2345},
		{0, "1.2345", nil, 1.2345},
		{0, 255, nil, 255},
		{0, -1.2345, nil, -1.2345},
		{0, "-1.2345", nil, -1.2345},
		{0, -255, nil, -255},
	}
	s := NewScope()

	for _, table := range tables {
		v := table.i
		ptr := reflect.ValueOf(&v)
		err := s.assignValue(ptr.Elem(), table.v)
		if err == table.e {
			if table.r != v {
				t.Errorf("expected %f, got %f", table.r, v)
			}
		} else {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestAssignFloat64(t *testing.T) {
	tables := []struct {
		i float64
		v interface{}
		e error
		r float64
	}{
		{0, 1.2345, nil, 1.2345},
		{0, "1.2345", nil, 1.2345},
		{0, 255, nil, 255},
		{0, -1.2345, nil, -1.2345},
		{0, "-1.2345", nil, -1.2345},
		{0, -255, nil, -255},
	}
	s := NewScope()

	for _, table := range tables {
		v := table.i
		ptr := reflect.ValueOf(&v)
		err := s.assignValue(ptr.Elem(), table.v)
		if err == table.e {
			if table.r != v {
				t.Errorf("expected %f, got %f", table.r, v)
			}
		} else {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestAssignInt64(t *testing.T) {
	tables := []struct {
		i int64
		v interface{}
		e bool
		r int64
	}{
		{0, 255, false, 255},
		{0, "155", false, 155},
		{0, "+123", false, 123},
		{0, "1+55", true, 155},
		{0, 55.1, false, 55},
		{0, -250, false, -250},
		{0, "-215", false, -215},
		{0, -25.2, false, -25},
		{0, -25.9, false, -25},
	}
	s := NewScope()

	for _, table := range tables {
		v := table.i
		ptr := reflect.ValueOf(&v)
		err := s.assignValue(ptr.Elem(), table.v)
		if table.e {
			if err == nil {
				t.Error("expected error")
			}
		} else {
			if err == nil {
				if table.r != v {
					t.Errorf("expected %d, got %d", table.r, v)
				}
			} else {
				t.Errorf("unexpected error: %s", err.Error())
			}
		}
	}
}

func TestAssignInt(t *testing.T) {
	tables := []struct {
		i int
		v interface{}
		e bool
		r int
	}{
		{0, 255, false, 255},
		{0, "155", false, 155},
		{0, "+123", false, 123},
		{0, "1+55", true, 155},
		{0, 55.1, false, 55},
		{0, -250, false, -250},
		{0, "-215", false, -215},
		{0, -25.2, false, -25},
		{0, -25.9, false, -25},
	}
	s := NewScope()

	for _, table := range tables {
		v := table.i
		ptr := reflect.ValueOf(&v)
		err := s.assignValue(ptr.Elem(), table.v)
		if table.e {
			if err == nil {
				t.Error("expected error")
			}
		} else {
			if err == nil {
				if table.r != v {
					t.Errorf("expected %d, got %d", table.r, v)
				}
			} else {
				t.Errorf("unexpected error: %s", err.Error())
			}
		}
	}
}

func TestAssignInt32(t *testing.T) {
	tables := []struct {
		i int32
		v interface{}
		e bool
		r int32
	}{
		{0, 255, false, 255},
		{0, "155", false, 155},
		{0, "+123", false, 123},
		{0, "1+55", true, 155},
		{0, 55.1, false, 55},
		{0, -250, false, -250},
		{0, "-215", false, -215},
		{0, -25.2, false, -25},
		{0, -25.9, false, -25},
	}
	s := NewScope()

	for _, table := range tables {
		v := table.i
		ptr := reflect.ValueOf(&v)
		err := s.assignValue(ptr.Elem(), table.v)
		if table.e {
			if err == nil {
				t.Error("expected error")
			}
		} else {
			if err == nil {
				if table.r != v {
					t.Errorf("expected %d, got %d", table.r, v)
				}
			} else {
				t.Errorf("unexpected error: %s", err.Error())
			}
		}
	}
}

func TestAssignInt16(t *testing.T) {
	tables := []struct {
		i int16
		v interface{}
		e bool
		r int16
	}{
		{0, 255, false, 255},
		{0, "155", false, 155},
		{0, "+123", false, 123},
		{0, "1+55", true, 155},
		{0, 55.1, false, 55},
		{0, -250, false, -250},
		{0, "-215", false, -215},
		{0, -25.2, false, -25},
		{0, -25.9, false, -25},
	}
	s := NewScope()

	for _, table := range tables {
		v := table.i
		ptr := reflect.ValueOf(&v)
		err := s.assignValue(ptr.Elem(), table.v)
		if table.e {
			if err == nil {
				t.Error("expected error")
			}
		} else {
			if err == nil {
				if table.r != v {
					t.Errorf("expected %d, got %d", table.r, v)
				}
			} else {
				t.Errorf("unexpected error: %s", err.Error())
			}
		}
	}
}

func TestAssignInt8(t *testing.T) {
	tables := []struct {
		i int8
		v interface{}
		e bool
		r int8
	}{
		{0, 127, false, 127},
		{0, 128, false, -128},
		{0, 129, false, -127},
		{0, 256, false, 0},
		{0, "55", false, 55},
		{0, "+123", false, 123},
		{0, "1+22", true, 0},
		{0, 55.1, false, 55},
		{0, -50, false, -50},
		{0, "-15", false, -15},
		{0, -25.2, false, -25},
		{0, -25.9, false, -25},
		{0, -128, false, -128},
		{0, -129, false, 127},
		{0, -256, false, 0},
	}
	s := NewScope()

	for _, table := range tables {
		v := table.i
		ptr := reflect.ValueOf(&v)
		err := s.assignValue(ptr.Elem(), table.v)
		if table.e {
			if err == nil {
				t.Error("expected error")
			}
		} else {
			if err == nil {
				if table.r != v {
					t.Errorf("expected %d, got %d", table.r, v)
				}
			} else {
				t.Errorf("unexpected error: %s", err.Error())
			}
		}
	}
}

func TestAssignUint64(t *testing.T) {
	tables := []struct {
		i uint64
		v interface{}
		e bool
		r uint64
	}{
		{1, nil, false, 0},
		{1, "", false, 0},
		{1, false, false, 0},
		{0, true, false, 1},
		{0, 255, false, 255},
		{0, "155", false, 155},
		{0, "155.0", true, 0},
		{0, "18446744073709551615", false, math.MaxUint64},
		{0, "18446744073709551616", true, 0},
		{0, "+123", true, 0},
		{0, "1+55", true, 0},
		{0, 55.1, false, 55},
		{0, -1, false, math.MaxUint64},
		{0, math.MinInt64, false, math.MaxInt64 + 1},
		{0, "-215", true, 0},
	}
	s := NewScope()

	for _, table := range tables {
		v := table.i
		ptr := reflect.ValueOf(&v)
		err := s.assignValue(ptr.Elem(), table.v)
		if table.e {
			if err == nil {
				t.Error("expected error")
			}
		} else {
			if err == nil {
				if table.r != v {
					t.Errorf("expected %d, got %d", table.r, v)
				}
			} else {
				t.Errorf("unexpected error: %s", err.Error())
			}
		}
	}
}

func TestAssignUint(t *testing.T) {
	tables := []struct {
		i uint
		v interface{}
		e bool
		r uint
	}{
		{0, 255, false, 255},
		{0, "155", false, 155},
		{0, "+123", true, 0},
		{0, "1+55", true, 0},
		{0, 55.1, false, 55},
		{0, -1, false, math.MaxUint64},
		{0, math.MinInt64, false, math.MaxInt64 + 1},
		{0, "-215", true, 0},
	}
	s := NewScope()

	for _, table := range tables {
		v := table.i
		ptr := reflect.ValueOf(&v)
		err := s.assignValue(ptr.Elem(), table.v)
		if table.e {
			if err == nil {
				t.Error("expected error")
			}
		} else {
			if err == nil {
				if table.r != v {
					t.Errorf("expected %d, got %d", table.r, v)
				}
			} else {
				t.Errorf("unexpected error: %s", err.Error())
			}
		}
	}
}

func TestAssignUint32(t *testing.T) {
	tables := []struct {
		i uint32
		v interface{}
		e bool
		r uint32
	}{
		{0, 255, false, 255},
		{0, "155", false, 155},
		{0, "+123", true, 0},
		{0, "1+55", true, 0},
		{0, 55.1, false, 55},
		{0, -1, false, math.MaxUint32},
		{0, math.MinInt32, false, math.MaxInt32 + 1},
		{0, "-215", true, 0},
	}
	s := NewScope()

	for _, table := range tables {
		v := table.i
		ptr := reflect.ValueOf(&v)
		err := s.assignValue(ptr.Elem(), table.v)
		if table.e {
			if err == nil {
				t.Error("expected error")
			}
		} else {
			if err == nil {
				if table.r != v {
					t.Errorf("expected %d, got %d", table.r, v)
				}
			} else {
				t.Errorf("unexpected error: %s", err.Error())
			}
		}
	}
}

func TestAssignUint16(t *testing.T) {
	tables := []struct {
		i uint16
		v interface{}
		e bool
		r uint16
	}{
		{0, 255, false, 255},
		{0, "155", false, 155},
		{0, "+123", true, 0},
		{0, "1+55", true, 0},
		{0, 55.1, false, 55},
		{0, -1, false, math.MaxUint16},
		{0, math.MinInt16, false, math.MaxInt16 + 1},
		{0, "-215", true, 0},
	}
	s := NewScope()

	for _, table := range tables {
		v := table.i
		ptr := reflect.ValueOf(&v)
		err := s.assignValue(ptr.Elem(), table.v)
		if table.e {
			if err == nil {
				t.Error("expected error")
			}
		} else {
			if err == nil {
				if table.r != v {
					t.Errorf("expected %d, got %d", table.r, v)
				}
			} else {
				t.Errorf("unexpected error: %s", err.Error())
			}
		}
	}
}

func TestAssignUint8(t *testing.T) {
	tables := []struct {
		i uint8
		v interface{}
		e bool
		r uint8
	}{
		{0, 255, false, 255},
		{0, "155", false, 155},
		{0, "257", false, 1},
		{0, "+123", true, 0},
		{0, "1+55", true, 0},
		{0, 55.1, false, 55},
		{0, -1, false, math.MaxUint8},
		{0, math.MinInt8, false, math.MaxInt8 + 1},
		{0, "-215", true, 0},
	}
	s := NewScope()

	for _, table := range tables {
		v := table.i
		ptr := reflect.ValueOf(&v)
		err := s.assignValue(ptr.Elem(), table.v)
		if table.e {
			if err == nil {
				t.Error("expected error")
			}
		} else {
			if err == nil {
				if table.r != v {
					t.Errorf("expected %d, got %d", table.r, v)
				}
			} else {
				t.Errorf("unexpected error: %s", err.Error())
			}
		}
	}
}

func TestAssignString(t *testing.T) {
	tables := []struct {
		i string
		v interface{}
		e bool
		r string
	}{
		{"hello", "", false, ""},
		{"hello", nil, false, ""},
		{"", "world", false, "world"},
		{"", 0, false, "0"},
		{"", 123, false, "123"},
		{"", -456, false, "-456"},
		{"", true, false, "true"},
		{"", false, false, "false"},
		{"", 1.2345, false, "1.2345"},
		{"", -1.2345, false, "-1.2345"},
	}
	s := NewScope()

	for _, table := range tables {
		v := table.i
		ptr := reflect.ValueOf(&v)
		err := s.assignValue(ptr.Elem(), table.v)
		if table.e {
			if err == nil {
				t.Error("expected error")
			}
		} else {
			if err == nil {
				if table.r != v {
					t.Errorf("expected %s, got %s", table.r, v)
				}
			} else {
				t.Errorf("unexpected error: %s", err.Error())
			}
		}
	}
}

func TestAssignSlice(t *testing.T) {
	tables := []struct {
		i []string
		v interface{}
		e bool
		r []string
	}{
		{[]string{"hello", "world"}, nil, false, []string{}},
		{[]string{"hello", "world"}, []string{}, false, []string{}},
		{[]string{}, []string{"hello", "world"}, false, []string{"hello", "world"}},
		{[]string{}, "hello,world", false, []string{"hello", "world"}},
		{[]string{}, "kangaroo", false, []string{"kangaroo"}},
		{[]string{}, []bool{true, false}, false, []string{"true", "false"}},
		{[]string{}, []interface{}{true, 1, -0.12345, "house"}, false, []string{"true", "1", "-0.12345", "house"}},
	}

	s := NewScope()

	for _, table := range tables {
		v := table.i
		ptr := reflect.ValueOf(&v)
		err := s.assignValue(ptr.Elem(), table.v)
		if table.e {
			if err == nil {
				t.Error("expected error")
			}
		} else {
			if err == nil {
				if !reflect.DeepEqual(table.r, v) {
					t.Errorf("expected %#v, got %#v", table.r, v)
				}
			} else {
				t.Errorf("unexpected error: %s", err.Error())
			}
		}
	}
}
