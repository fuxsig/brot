// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brot

import (
	"log"
	"net/http"
	"time"

	"github.com/fuxsig/brot/di"
)

// Dispatcher dispatches
type Dispatcher interface {
	Dispatch(http.ResponseWriter, *http.Request)
}

type Server struct {
	Addr         string `brot:"addr"`
	WriteTimeout int    `brot:"writeTimeout"`
	ReadTimeout  int    `brot:"readTimeout"`
	Router       string `brot:"router"`
	CertPath     string `brot:cert`
	KeyPath      string `brot:key`
}

// InitFunc initialize opens a new connection to the configured Redis database
func (s *Server) InitFunc() {
	server := new(http.Server)
	if s.Addr != "" {
		server.Addr = s.Addr
	}
	if s.WriteTimeout >= 0 {
		server.WriteTimeout = time.Duration(s.WriteTimeout) * time.Second
	}
	if s.ReadTimeout >= 0 {
		server.ReadTimeout = time.Duration(s.ReadTimeout) * time.Second
	}
	if s.Router != "" {
		obj := di.GlobalScope.Get(s.Router)
		if obj == nil {
			log.Printf("Warning: cannot find router %s", s.Router)
			return
		}
		handler, ok := obj.(http.Handler)
		if !ok {
			log.Printf("Warning: skipping router %s. Router is not http.Handler", s.Router)
			return
		}
		server.Handler = handler

	}

	if s.CertPath != "" || s.KeyPath != "" {
		if s.CertPath == "" || s.KeyPath == "" {
			log.Fatalf("Certificate path and key path must be set")
		}
		go func() {
			if err := server.ListenAndServeTLS("/Users/mbungens/server.crt", "/Users/mbungens/server.key"); err != nil {
				log.Fatalln(err.Error())
			}
		}()
	} else {
		go func() {
			if err := server.ListenAndServe(); err != nil {
				log.Fatalln(err.Error())
			}
		}()
	}
}

var _ di.ProvidesInit = (*Server)(nil)
var _ = di.GlobalScope.Declare((*Server)(nil))
