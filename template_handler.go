// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brot

import (
	"fmt"
	"net/http"

	"github.com/fuxsig/brot/di"
	"github.com/fuxsig/brot/model"
)

type TemplateHandler struct {
	Model     *model.Butter     `brot:"model"`
	Templates ProvidesTemplates `brot:"templates"`
	Data      *DataLayer        `brot:"data"`
}

func (th *TemplateHandler) InitFunc() {
	if th.Templates == nil {
		panic("th.Templates == nil")
	}
	if th.Data == nil {
		panic("th.Data == nil")
	}
}

func (th *TemplateHandler) HandlerFunc() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if Debug {
			Debug.Printf("TemplateHandler begins")
			defer Debug.Printf("TemplateHandler ends")
		}
		values := th.Data.Values(r)
		id, ok := values.Get("id")
		if !ok {
			http.Error(w, "Missing id value", http.StatusBadRequest)
			return
		}
		var template string
		template, ok = values.Get("template")
		if !ok {
			http.Error(w, "Missing template value", http.StatusBadRequest)
			return
		}
		_, obj, err := th.Model.Load(nil, id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Load error: %s", err.Error()), http.StatusBadRequest)
			return
		}
		err = th.Templates.TemplateFunc().ExecuteTemplate(w, template, obj)
		if err != nil {
			http.Error(w, fmt.Sprintf("Template error: %s", err.Error()), http.StatusInternalServerError)
		}
	})
}

var _ di.ProvidesInit = (*TemplateHandler)(nil)
var _ ProvidesHandler = (*TemplateHandler)(nil)
var _ = di.GlobalScope.Declare((*TemplateHandler)(nil))
