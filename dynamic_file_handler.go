// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brot

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"path/filepath"

	"github.com/fuxsig/brot/di"
)

type DynamicFileHandler struct {
	Dir string `brot:"dir"`
}

func (dh *DynamicFileHandler) HandlerFunc() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Cache-Control", "no-cache, no-store, must-revalidate")
		h.Set("Pragma", "no-cache")
		h.Set("Expires", "0")

		session, err := sessionStore.Get(r, "okta-hosted-login-session-store")
		if err != nil {
			log.Printf("DynamicFileHandler: Session error: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		p := filepath.Join(dh.Dir, r.URL.Path)
		base := path.Base(p)

		t := template.New(base)

		if t, err = t.ParseFiles(p); err != nil {
			io.WriteString(w, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		values := make(map[string]string)
		values["email"], _ = session.Values["email"].(string)
		values["given_name"], _ = session.Values["given_name"].(string)
		if err = t.Execute(w, values); err != nil {
			io.WriteString(w, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	})
}

var _ ProvidesHandler = (*DynamicFileHandler)(nil)
var _ = di.GlobalScope.Declare((*DynamicFileHandler)(nil))
