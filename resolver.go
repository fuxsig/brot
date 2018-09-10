// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brot

import (
	"net/http"
	"path"
	"regexp"

	"github.com/fuxsig/brot/di"
	"github.com/gorilla/mux"
)

type Resolver interface {
	Resolve(*http.Request, map[string]string)
}

type MuxResolver struct {
	Mapping map[string]string
}

func (mr *MuxResolver) Resolve(r *http.Request, result map[string]string) {
	vars := mux.Vars(r)
	for k, v := range mr.Mapping {
		result[k] = vars[v]
	}
}

var _ = di.GlobalScope.Declare((*MuxResolver)(nil))

type QueryResolver struct {
	Mapping map[string]string
}

func (mr *QueryResolver) Resolve(r *http.Request, result map[string]string) {
	values := r.URL.Query()
	for k, v := range mr.Mapping {
		result[k] = values.Get(v)
	}
}

var _ = di.GlobalScope.Declare((*QueryResolver)(nil))

type URLBaseResolver struct {
	RegexpStr string `brot:"regexpStr"`
	re        *regexp.Regexp
	names     []string
}

func (mr *URLBaseResolver) InitFunc() {
	mr.re = regexp.MustCompile(mr.RegexpStr)
	mr.names = mr.re.SubexpNames()
}

func (mr *URLBaseResolver) Resolve(r *http.Request, result map[string]string) {
	base := path.Base(r.URL.Path)
	matches := mr.re.FindStringSubmatch(base)
	for i, str := range matches {
		if i != 0 {
			result[mr.names[i]] = str
		}
	}
}

var _ = di.GlobalScope.Declare((*URLBaseResolver)(nil))
