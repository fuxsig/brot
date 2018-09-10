// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brot

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/fuxsig/brot/di"
)

type Object struct {
	Name       string                 `json:"name"`
	StructName string                 `json:"struct"`
	FuncName   string                 `json:"func"`
	Args       map[string]interface{} `json:"args"`
}

// Configuration contains all configuration settings
type Configuration struct {
	Handlers []Object `json:"handlers"`

	Views struct {
		Paths []string `json:"paths"`
	} `json:"views"`
}

// ViewRegistry creates a new ViewRegistry object from the configuration
func (conf *Configuration) ViewRegistry() (result *ViewRegistry) {
	result = new(ViewRegistry)
	result.Init(conf.Views.Paths)
	return
}

// LoadConfiguration loads a configuration
func LoadConfiguration(path string) (result *Configuration, err error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	result, err = ParseConfiguration(raw)
	return
}

// ParseConfiguration parses the configuration data
func ParseConfiguration(raw []byte) (result *Configuration, err error) {
	result = new(Configuration)
	err = json.Unmarshal(raw, result)
	if err != nil {
		result = nil
	}
	return
}

func map2array(m *map[string]string) (result []string) {
	result = make([]string, len(*m)*2)
	i := 0
	for k, v := range *m {
		result[i] = k
		i++
		result[i] = v
		i++
	}
	return
}

func (c Configuration) Process() {
	// initialize all handler objects
	for _, handler := range c.Handlers {
		if handler.StructName != "" && handler.FuncName != "" {
			log.Panicf("Invalid configuration, please use either struct or func. Current values are struct '%s' and func '%s'", handler.StructName, handler.FuncName)
		}
		if handler.StructName == "" && handler.FuncName == "" {
			log.Panic("Invalid configuration, please set at least struct or func.")
		}
		if handler.StructName != "" {
			if _, err := di.GlobalScope.New(handler.Name, handler.StructName, handler.Args); err == nil {
				log.Printf("Created successfully struct %s of type %s", handler.Name, handler.StructName)
			} else {
				log.Printf("Could not create object %s. Reason: %s", handler.Name, err.Error())
			}

		} else {
			if _, err := di.GlobalScope.Call(handler.Name, handler.FuncName, handler.Args); err == nil {
			}
		}
	}

}
