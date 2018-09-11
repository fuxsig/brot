// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brot

import (
	"log"

	"github.com/fuxsig/brot/di"
	"github.com/gorilla/mux"
)

// GorillaRouter defines the rules for the Gorilla multiplexer
type GorillaRouter struct {
	Name      string   `brot:"name"`
	Subrouter string   `brot:"subrouter"`
	Use       []string `brot:"use"`
	Routes    []struct {
		Name    string            `brot:"name"`
		Path    string            `brot:"path"`
		Handler string            `brot:"handler"`
		Host    string            `brot:"host"`
		Prefix  string            `brot:"prefix"`
		Methods []string          `brot:"methods"`
		Schemes []string          `brot:"schemes"`
		Headers map[string]string `brot:"headers"`
		Queries map[string]string `brot:"queries"`
	} `brot:"routes"`
}

// InitFunc initialize opens a new connection to the configured Redis database
func (gr *GorillaRouter) InitFunc() (err error) {

	var router *mux.Router
	if gr.Subrouter != "" {
		if subrouter := di.GlobalScope.Get(gr.Subrouter); subrouter != nil {
			if subrouter, ok := subrouter.(*mux.Router); ok {
				router = subrouter
			} else {
				log.Printf("Warning: skipping route %s. Subrouter %s is not a mux.Router struct", gr.Name, gr.Subrouter)
				return
			}
		} else {
			log.Printf("Warning: skipping route %s. Could not find subrouter %s", gr.Name, gr.Subrouter)
			return
		}
	} else {
		router = mux.NewRouter()
		if gr.Name != "" {
			di.GlobalScope.Set(gr.Name, router)
		}
	}

	for _, current := range gr.Routes {
		// do we have a handler object with the given name?
		if handler := di.GlobalScope.Get(current.Handler); handler != nil {
			if handler, ok := handler.(ProvidesHandler); ok {

				route := router.NewRoute()

				if current.Name != "" {
					route.Name(current.Name)
				}
				if current.Path != "" {
					route.Path(current.Path)
				} else {
					if handler, ok := handler.(ProvidesPath); ok {
						route.Path(handler.PathFunc())
					}
				}
				if current.Prefix != "" {
					route.PathPrefix(current.Prefix)
				} else {
					if handler, ok := handler.(ProvidesPrefix); ok {
						route.PathPrefix(handler.PrefixFunc())
					}
				}

				if current.Host != "" {
					route.Host(current.Host)
				}
				if len(current.Methods) > 0 {
					route.Methods(current.Methods...)
				}
				if len(current.Schemes) > 0 {
					route.Schemes(current.Schemes...)
				}
				if len(current.Headers) > 0 {
					route.Headers(map2array(&current.Headers)...)
				}
				if len(current.Queries) > 0 {
					route.Queries(map2array(&current.Queries)...)
				}
				route.Handler(handler.HandlerFunc())
			} else {
				log.Printf("Warning: skipping route for %s. Handler %s does not support interface ProvidesHandler", current.Path, current.Handler)
			}
		} else {
			log.Printf("Warning: skipping route for %s. Could not find handler %s", current.Path, current.Handler)
		}
	}
	// http wrapper or middleware in mux language
	for _, use := range gr.Use {
		if wrapper := di.GlobalScope.Get(use); wrapper != nil {
			if wrapper, ok := wrapper.(ProvidesWrapper); ok {
				router.Use(wrapper.WrapperFunc())
			} else {
				log.Printf("Warning: skipping wrapper for %s. Object %s does not support interface ProvidesWrapper", gr.Name, use)
			}
		} else {
			log.Printf("Warning: skipping wrapper for %s. Object %s does not exist", gr.Name, use)
		}
	}
	return
}

func (gr *GorillaRouter) Retry() bool {
	return false
}

var _ di.ProvidesInit = (*GorillaRouter)(nil)
var _ = di.GlobalScope.Declare((*GorillaRouter)(nil))
