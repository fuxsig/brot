// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package di

import (
	"testing"
)

func TestGetBool(t *testing.T) {
	tables := []struct {
		v interface{}
		r bool
	}{
		{true, true},
		{"true", true},
		{"True", true},
		{"TRUE", true},
		{"t", true},
		{"T", true},
		{(uint8)(1), true},
		{(uint16)(1), true},
		{(uint32)(1), true},
		{(uint)(1), true},
		{(uint64)(1), true},
		{(int8)(1), true},
		{(int16)(1), true},
		{(int32)(1), true},
		{(int)(1), true},
		{(int64)(1), true},

		{false, false},
		{"false", false},
		{"False", false},
		{"FALSE", false},
		{"f", false},
		{"F", false},
		{(uint8)(0), false},
		{(uint16)(0), false},
		{(uint32)(0), false},
		{(uint)(0), false},
		{(uint64)(0), false},
		{(int8)(0), false},
		{(int16)(0), false},
		{(int32)(0), false},
		{(int)(0), false},
		{(int64)(0), false},
	}

	for _, table := range tables {
		r, err := GetBool(table.v)
		if err == nil {
			if table.r != r {
				t.Errorf("expected %t, got %t", table.r, r)
			}
		} else {
			t.Errorf("unexpected error: %s", err.Error())
		}

	}

}
