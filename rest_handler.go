// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brot

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/fuxsig/brot/di"
	"github.com/fuxsig/brot/model"
)

type RestHandler struct {
	Model *model.Butter `brot:"model"`
	Data  *DataLayer    `brot:"data"`
}

func (h *RestHandler) HandlerFunc() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := h.Data.Values(r)
		query := r.URL.Query()
		var (
			err  error
			str  []string
			sort string
		)
		id, ok := values.Get("id")
		if !ok {
			var schema string
			if schema, ok = values.Get("schema"); !ok {
				http.Error(w, "schema missing", http.StatusBadRequest)
				return
			}

			offset := 0
			limit := 10
			if str, ok = query["offset"]; ok {
				if offset, err = strconv.Atoi(str[0]); err != nil {
					http.Error(w, "invalid offset", http.StatusBadRequest)
					return
				}
			}
			if str, ok = query["limit"]; ok {
				if limit, err = strconv.Atoi(str[0]); err != nil {
					http.Error(w, "invalid limit", http.StatusBadRequest)
					return
				}
			}
			if str, ok = query["sort"]; ok {
				sort = str[0]
			}

			result, total, err := h.Model.Search(nil, schema, "", sort, offset, limit)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			i := len(result)
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")

			io.WriteString(w, `{"offset":`)
			io.WriteString(w, strconv.Itoa(offset))
			io.WriteString(w, `,"number":`)
			io.WriteString(w, strconv.Itoa(i))
			io.WriteString(w, `,"total":`)
			io.WriteString(w, strconv.Itoa(total))
			io.WriteString(w, `,"objects":[`)
			var j []byte
			for _, current := range result {
				if j, err = json.Marshal(current); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Write(j)
				i--
				if i > 0 {
					io.WriteString(w, ",")
				}
			}
			io.WriteString(w, "]}")

			return
		}

		status, data, err := h.Model.Load(nil, id)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}
		var j []byte
		if j, err = json.Marshal(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(status)
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)

	})
}

var _ ProvidesHandler = (*RestHandler)(nil)
var _ = di.GlobalScope.Declare((*RestHandler)(nil))
