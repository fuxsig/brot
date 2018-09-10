// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brot

import (
	"net/http"
)

type ProvidesPath interface {
	PathFunc() string
}

type ProvidesPrefix interface {
	PrefixFunc() string
}

type ProvidesHandler interface {
	HandlerFunc() http.Handler
}

type ProvidesWrapper interface {
	WrapperFunc() func(http.Handler) http.Handler
}

type ProvidesDatabase interface {
	Get(string, interface{}) error
	GetPtr(string, string) interface{}
	Set(string, interface{}) error
}
