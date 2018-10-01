// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package di

import (
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
	s := NewScope()
	v := true
	ptr := reflect.ValueOf(&v)
	err := s.assignValue(ptr.Elem(), false)
	if err == nil {
		if v != false {
			t.Error("expected false")
		}
	} else {
		t.Errorf("unexpected error: %s", err.Error())
	}

	err = s.assignValue(ptr.Elem(), true)
	if err == nil {
		if v != true {
			t.Error("expected true")
		}
	} else {
		t.Errorf("unexpected error: %s", err.Error())
	}

}
