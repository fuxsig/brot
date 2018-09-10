// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brot

import (
	"net/http"

	"github.com/fuxsig/brot/di"
	"github.com/gorilla/mux"
)

type MultiValueMap map[string][]string

func (de MultiValueMap) Get(name string) (string, bool) {
	aux := de[name]
	if len(aux) > 0 {
		return aux[0], true
	}
	return "", false
}

type ValueProvider interface {
	Initialize(request *http.Request, m map[string][]string)
	CapacityReco() int
}

type DataLayer struct {
	Providers []ValueProvider `brot:"providers"`
	capacity  int
}

func (dl *DataLayer) InitFunc() {
	dl.capacity = 0
	for _, p := range dl.Providers {
		dl.capacity += p.CapacityReco()
	}
}

func (dl *DataLayer) Values(request *http.Request) (result MultiValueMap) {
	result = make(MultiValueMap, dl.capacity)
	for _, p := range dl.Providers {
		p.Initialize(request, result)
	}
	return
}

var _ di.ProvidesInit = (*DataLayer)(nil)
var _ = di.GlobalScope.Declare((*DataLayer)(nil))

type MuxVarsProvider struct {
	Mapping map[string]string `brot:"mapping"`
}

func (me *MuxVarsProvider) Initialize(request *http.Request, m map[string][]string) {
	vars := mux.Vars(request)
	for key, value := range me.Mapping {
		parameter := vars[key]
		if parameter != "" {
			current := m[value]
			if current != nil {
				m[value] = append(current, parameter)
			} else {
				m[value] = []string{parameter}
			}
		}
	}
}

func (me *MuxVarsProvider) CapacityReco() int {
	return len(me.Mapping)
}

var _ ValueProvider = (*MuxVarsProvider)(nil)
var _ = di.GlobalScope.Declare((*MuxVarsProvider)(nil))

type URLParameterProvider struct {
	Mapping map[string]string `brot:"mapping"`
}

func (ue *URLParameterProvider) Initialize(request *http.Request, m map[string][]string) {
	query := request.URL.Query()
	for key, value := range ue.Mapping {
		parameters := query[key]
		if parameters != nil {
			current := m[value]
			if current != nil {
				m[value] = append(current, parameters...)
			} else {
				m[value] = parameters
			}
		}
	}
}

func (ue *URLParameterProvider) CapacityReco() int {
	return len(ue.Mapping)
}

var _ ValueProvider = (*URLParameterProvider)(nil)
var _ = di.GlobalScope.Declare((*URLParameterProvider)(nil))

type ConstValuesProvider struct {
	Mapping map[string]string `brot:"mapping"`
}
