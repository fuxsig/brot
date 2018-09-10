// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brot

import (
	"log"
	"net/http"

	"github.com/fuxsig/brot/di"
)

type LogWrapper struct {
}

func (lw *LogWrapper) WrapperFunc() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Request for ", r.URL.String())
			defer log.Println("Done with ", r.URL.String())
			h.ServeHTTP(w, r)

		})
	}
}

var _ = di.GlobalScope.Declare((*LogWrapper)(nil))
