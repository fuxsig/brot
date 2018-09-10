// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brot

import (
	"net/http"

	"github.com/fuxsig/brot/di"
)

type StaticFileHandler struct {
	Dir  string `brot:"dir"`
	Path string `brot:"path"`
}

func (sh StaticFileHandler) HandlerFunc() http.Handler {
	return http.StripPrefix(sh.Path, http.FileServer(http.Dir(sh.Dir)))
}

func (sh StaticFileHandler) PrefixFunc() string {
	return sh.Path
}

var _ = di.GlobalScope.Declare((*StaticFileHandler)(nil))
