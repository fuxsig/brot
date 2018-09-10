// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"bytes"
	"strconv"
)

type Vanilla struct {
	Schema *Schema
	Data   map[string]string
}

func (v *Vanilla) Label() string {
	return v.Data["_label"]
}

func (v *Vanilla) ID() string {
	return v.Data["_id"]
}

func (v *Vanilla) Int(name string) int {
	result, _ := strconv.Atoi(v.Data[name])
	return result
}

func (v *Vanilla) String(name string) string {
	return v.Data[name]
}

func (v *Vanilla) Json() []byte {
	var b bytes.Buffer
	b.WriteString("{")
	i := len(v.Data)
	for key, value := range v.Data {
		b.WriteString(`"`)
		b.WriteString(key)
		b.WriteString(`":"`)
		b.WriteString(value)
		i--
		if i > 0 {
			b.WriteString(`",`)
		} else {
			b.WriteString(`"`)
		}
	}
	b.WriteString("}")
	return b.Bytes()
}

func (v *Vanilla) JsonPartial(fields ...string) []byte {
	var b bytes.Buffer
	b.WriteString("{")
	i := len(fields)
	for _, field := range fields {
		b.WriteString(`"`)
		b.WriteString(field)
		b.WriteString(`":"`)
		b.WriteString(v.Data[field])
		i--
		if i > 0 {
			b.WriteString(`",`)
		} else {
			b.WriteString(`"`)
		}
	}
	b.WriteString("}")
	return b.Bytes()
}
