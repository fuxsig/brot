// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package di

import (
	"log"
	"reflect"
)

type filler func(args map[string]interface{}) (reflect.Value, error)

type DynamicFunc struct {
	fills []filler
	funk  reflect.Value
	index int
}

func NewDynamicFunc(f interface{}, argNames []string) (result *DynamicFunc) {
	v := reflect.ValueOf(f)
	t := v.Type()
	num := t.NumIn()
	if num != len(argNames) {
		log.Panicf("Number of func arguments and passed names does not match. Expected %d, received %d", num, len(argNames))
	}
	result = &DynamicFunc{}
	result.index = 0
	result.fills = make([]filler, num)
	if t.Kind() == reflect.Func {
		result.funk = v
		for i := 0; i < num; i++ {
			pt := t.In(i)
			name := argNames[i]
			kind := pt.Kind()
			switch kind {
			case reflect.Bool:
				result.fills[i] = func(args map[string]interface{}) (result reflect.Value, err error) {
					if val, ok := args[name]; ok {
						var b bool
						if b, err = GetBool(val); err == nil {
							result = reflect.ValueOf(b)
						}
					}
					return
				}
			case reflect.String:
				result.fills[i] = func(args map[string]interface{}) (result reflect.Value, err error) {
					if val, ok := args[name]; ok {
						var str string
						if str, err = GetString(val); err == nil {
							result = reflect.ValueOf(str)
						}
					}
					return
				}
			case reflect.Int:
				result.fills[i] = func(args map[string]interface{}) (result reflect.Value, err error) {
					if val, ok := args[name]; ok {
						var v int
						if v, err = GetInt(val); err == nil {
							result = reflect.ValueOf(v)
						}
					}
					return
				}
			case reflect.Interface:
				result.fills[i] = func(args map[string]interface{}) (result reflect.Value, err error) {
					if val, ok := args[name]; ok {
						if name, ok := val.(string); ok {
							result = reflect.ValueOf(GlobalScope.Get(name))
						}
					}
					return
				}
			}
		}
	}
	return
}

func (df *DynamicFunc) Call(args map[string]interface{}) (result interface{}) {
	in := make([]reflect.Value, len(df.fills))
	for i, f := range df.fills {
		if in[i], err = f(args); err != nil {
			log.Panicf("Could not set value: %s", err.Error())
		}
	}
	out := df.funk.Call(in)
	if len(out) > 0 {
		result = out[df.index].Interface()
	}
	return
}
