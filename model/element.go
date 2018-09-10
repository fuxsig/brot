// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Element struct {
	Name  string `brot:"name"`
	Type  string `brot:"type"`
	Index bool   `brot:"index"`
	Sort  bool   `brot:"index"`

	slice   bool
	pointer bool
	kind    reflect.Kind
	check   func(string) error
}

func (e *Element) initialize(schemas map[string]*Schema) error {
	t := strings.TrimSpace(e.Type)
	e.Type = t
	// is it a slice?
	s := strings.HasPrefix(t, "[]")
	e.slice = s
	if s {
		t = t[2:]
	}
	// is it a pointer?
	p := strings.HasPrefix(t, "*")
	e.pointer = p
	if p {
		t = t[1:]
		// is it a defined structure?
		_, ok := schemas[t]
		if !ok {
			return fmt.Errorf("Element %s of type %s references unknown schems %s", e.Name, e.Type, t)
		}
		e.kind = reflect.Struct
		e.check = func(value string) error {
			if value == "" {
				return nil
			}
			return nil
		}
		return nil

	}
	// is it a predefined type?
	switch t {
	case "string":
		e.kind = reflect.String
		e.check = func(string) error { return nil }
	case "int":
		e.kind = reflect.Int
		e.check = func(value string) error {
			_, err := strconv.ParseInt(value, 10, 64)
			return err
		}
	default:
		return fmt.Errorf("Invalid declaration: %s", e.Type)
	}
	return nil
}
