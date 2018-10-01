// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package di

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/fuxsig/brot/aux"
)

type ProvidesInit interface {
	InitFunc() error
	Retry() bool
}

type Constructor interface {
	Allocated()
}

type Scope struct {
	parent  *Scope
	types   map[string]reflect.Type
	objects map[string]interface{}
	funcs   map[string]interface{}
}

var GlobalScope = NewScope()

func NewScope() (result *Scope) {
	result = new(Scope)
	result.types = make(map[string]reflect.Type, 0)
	result.objects = make(map[string]interface{}, 0)
	result.funcs = make(map[string]interface{}, 0)
	return
}

var brotTag, err = regexp.Compile(`brot:()`)

var (
	UnaddressableError = errors.New("value is unaddressable")
	UnexportedError    = errors.New("value is unexported struct field")
)

func (scope *Scope) assignValue(dest reflect.Value, src interface{}) error {
	var me aux.MultiError
	if !dest.CanSet() {
		if !dest.CanAddr() {
			return UnaddressableError
		}
		return UnexportedError
	}
	switch dest.Kind() {
	case reflect.Bool:
		if b, err := GetBool(src); err == nil {
			dest.SetBool(b)
		} else {
			me.Append(err)
		}
	case reflect.Float32, reflect.Float64:
		if f, err := GetFloat64(src); err == nil {
			dest.SetFloat(f)
		} else {
			me.Append(err)
		}
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		if i, err := GetInt64(src); err == nil {
			dest.SetInt(i)
		} else {
			me.Append(err)
		}
	case reflect.String:
		if str, err := GetString(src); err == nil {
			dest.SetString(str)
		} else {
			me.Append(err)
		}
	case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uint:
		if ui, err := GetUint64(src); err == nil {
			dest.SetUint(ui)
		} else {
			me.Append(err)
		}
	case reflect.Slice:
		// returns string for []string
		et := dest.Type().Elem()
		// value of val
		vv := reflect.ValueOf(src)
		st := vv.Type()
		cap := 1

		switch st.Kind() {
		case reflect.Slice:
			cap = vv.Len()
			slice := reflect.MakeSlice(reflect.SliceOf(et), 0, cap)
			for i := 0; i < cap; i++ {
				act := reflect.New(et).Elem()
				if err := scope.assignValue(act, vv.Index(i).Interface()); err == nil {
					slice = reflect.Append(slice, act)
				} else {
					me.Merge(err)
				}
			}
			dest.Set(slice)
		default:
			panic("Not yet implemented")
		}

	case reflect.Map:
		dt := dest.Type()
		vt := dt.Elem()
		kt := dt.Key()

		vv := reflect.ValueOf(src)
		st := vv.Type()
		var dict reflect.Value
		switch st.Kind() {
		case reflect.Map:
			dict = reflect.MakeMapWithSize(dt, vv.Len())
			for _, key := range vv.MapKeys() {
				val := vv.MapIndex(key)
				nKey := reflect.New(kt).Elem()
				if err := scope.assignValue(nKey, key.Interface()); err == nil {
					nVal := reflect.New(vt).Elem()
					if err := scope.assignValue(nVal, val.Interface()); err == nil {
						dict.SetMapIndex(nKey, nVal)
					} else {
						me.Merge(err)
					}
				} else {
					me.Merge(err)
				}

			}
		default:
			panic("Not yet implemented")
		}
		dest.Set(dict)
	case reflect.Struct:
		t := dest.Type()
		if m, ok := src.(map[string]interface{}); ok {
			for i := 0; i < dest.NumField(); i++ {
				// struct field
				sf := dest.Field(i)
				// struct field type
				sft := t.Field(i)
				name := sft.Name
				if !sf.CanSet() {
					continue
				}
				// processing the tag
				mandatory := false
				if val, ok := sft.Tag.Lookup("brot"); ok {
					tags := strings.Split(val, ",")
					if len(tags) > 0 {
						name = strings.Trim(tags[0], " ")
					}
					for i := 1; i < len(tags); i++ {
						if tags[i] == "mandatory" {
							mandatory = true
						}
					}
				}

				if val, ok := m[name]; ok {
					me.Merge(scope.assignValue(sf, val))
				} else {
					if mandatory {
						log.Printf("Mandatory value for %s not defined in configuration", name)
					}
				}
			}

			ptr := dest.Addr().Interface()
			if aux, ok := ptr.(Constructor); ok {
				aux.Allocated()
			}
			if aux, ok := ptr.(ProvidesInit); ok {
				r := true
				for r {
					if err := aux.InitFunc(); err != nil {
						log.Printf("Error in init: %s\n", err.Error())
						r = aux.Retry()
						if r {
							log.Printf("Retrying init\n")
						}
					} else {
						break
					}
				}
			}
		} else {

			log.Printf("Expected a source value of type map[string]interface{}, found %s", reflect.TypeOf(src).String())

		}
	case reflect.Ptr:

		// get the type of the referenced object, e.g. *string -> string
		et := dest.Type().Elem()
		// get the address of the pointer, e.g. *string -> **string
		pptr := dest.Addr().Elem()
		// does the pointer references a struct?
		if et.Kind() == reflect.Struct {
			// a map means an anonymous strcut
			if reflect.TypeOf(src).Kind() == reflect.Map {
				// create a new struct
				elem := reflect.New(et)
				if err := scope.assignValue(elem.Elem(), src); err == nil {
					pptr.Set(elem)
				} else {
					me.Merge(err)
				}
			} else {
				if str, err := GetString(src); err == nil {
					if val := scope.Get(str); val != nil {
						elem := reflect.ValueOf(val)
						if elem.Type() == dest.Type() {
							pptr.Set(elem)
						} else {
							log.Printf("Cannot assign value. Expecting type %s, %s is %s", dest.Type().Name(), str, elem.Type().Name())
						}

					} else {
						log.Printf("Cannot find object %s", str)
					}
				} else {
					me.Append(err)
				}

			}
		} else {
			ptr := reflect.New(et)
			if err := scope.assignValue(ptr.Elem(), src); err == nil {
				pptr.Set(ptr)
			} else {
				me.Merge(err)
			}
		}

	case reflect.Interface:
		if dest.NumMethod() == 0 {
			vv := reflect.ValueOf(src)
			dest.Set(vv)
			break
		}

		if m, ok := src.(map[string]interface{}); ok {
			// get struct name
			if sn, ok := m["struct"]; ok {
				if sn, ok := sn.(string); ok {
					if args, ok := m["args"]; ok {
						if args, ok := args.(map[string]interface{}); ok {
							if val, err := scope.Allocate(sn, args); err == nil {
								elem := reflect.ValueOf(val)
								if elem.Type().Implements(dest.Type()) {
									dest.Set(elem)
								} else {

									log.Printf("Object does not implement interface %s", dest.Type().Name())

								}
							} else {
								me.Append(err)
							}
						} else {
							me.AppendString("args is not of type map[string]interface{}")
						}
					} else {
						me.AppendString("missing args in object definition")
					}
				} else {
					me.AppendString("args is not of type string")
				}
			} else {
				me.AppendString("missing struct in object definition")
			}
		} else {
			if str, err := GetString(src); err == nil {

				if val := scope.Get(str); val != nil {
					elem := reflect.ValueOf(val)
					if elem.Type().Implements(dest.Type()) {
						dest.Set(elem)
					} else {

						log.Printf("Object does not implement interface %s", dest.Type().Name())

					}
				} else {
					log.Printf("Cannot find object %s", str)
				}
			} else {
				me.Append(err)
			}
		}
	default:
		sft := dest.Type()
		log.Printf("Unsupported type for field %s", sft.String())
	}
	return me.ErrorOrNil()
}

var errorNotAStruct = errors.New("passed value is not a pointer to struct")

func (scope *Scope) Assign(v interface{}, m map[string]interface{}) (err error) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		err = scope.assignValue(reflect.ValueOf(v).Elem(), m)
	} else {
		err = errorNotAStruct
	}
	return
}

func (scope *Scope) Get(name string) (result interface{}) {
	current := scope
	for result == nil && current != nil {
		result = current.objects[name]
		current = current.parent
	}
	return
}

func (scope *Scope) Set(name string, object interface{}) {
	scope.objects[name] = object
}

func (scope *Scope) Declare(object interface{}) interface{} {
	t := reflect.TypeOf(object)
	if t.Kind() == reflect.Ptr {
		for t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() != reflect.Struct {
			log.Panicf("Expected a struct but received %s", t.Name())
		}
		scope.types[t.String()] = t
	} else if t.Kind() == reflect.Func {
		scope.funcs[t.String()] = object
	}
	return object
}

func (scope *Scope) TypeOf(name string) reflect.Type {
	return scope.types[name]
}

func (scope *Scope) New(name string, typeName string, m map[string]interface{}) (result interface{}, err error) {
	result, err = scope.Allocate(typeName, m)
	if err == nil && name != "" {
		scope.Set(name, result)
	}
	return
}

func (scope *Scope) Allocate(typeName string, m map[string]interface{}) (result interface{}, err error) {
	if t := scope.types[typeName]; t != nil {
		ptr := reflect.New(t)
		if err = scope.assignValue(ptr.Elem(), m); err == nil {
			result = ptr.Interface()
		}
	} else {
		err = fmt.Errorf("Could not find type %s", typeName)
	}
	return
}

func (scope *Scope) Create(typeName string) (result reflect.Value, err error) {
	if t := scope.types[typeName]; t != nil {
		result = reflect.New(t)
	} else {
		err = fmt.Errorf("Could not find type %s", typeName)
	}
	return
}

func (scope *Scope) Call(objName string, funcName string, args map[string]interface{}) (result interface{}, err error) {
	// search for the object
	if f := scope.funcs[funcName]; f != nil {
		// is it a DynamicFunc?
		if df, ok := f.(*DynamicFunc); ok {
			// call DynamicFunc with arguments
			if obj := df.Call(args); obj != nil {
				// store the result
				scope.objects[objName] = obj
				result = obj
			}
		}
	} else {
		err = fmt.Errorf("Could not find constructor %s", funcName)
	}
	return
}

func (scope *Scope) RegisterDefaults() *Scope {
	scope.objects["os.Stdout"] = os.Stdout
	scope.funcs["log.New"] = NewDynamicFunc(log.New, []string{"out", "prefix", "flag"})
	return scope
}
