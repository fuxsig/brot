// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brot

import (
	"net/http"
	"path/filepath"

	"github.com/fuxsig/brot/di"
)

type TemplateFileHandler struct {
	Dir       string     `brot:"dir"`
	Path      string     `brot:"path"`
	Resolvers []Resolver `brot:"resolvers"`
}

func (sh *TemplateFileHandler) HandlerFunc() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := make(map[string]string, 0)
		for _, resolver := range sh.Resolvers {
			resolver.Resolve(r, vars)
		}

		path := filepath.Clean(r.URL.Path)
		_ = path
	})
}

func (sh *TemplateFileHandler) PrefixFunc() string {
	return sh.Path
}

var _ = di.GlobalScope.Declare((*TemplateFileHandler)(nil))
