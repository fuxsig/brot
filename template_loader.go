// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brot

import (
	"errors"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fuxsig/brot/model"

	"github.com/fuxsig/brot/di"
)

type TemplateLoader struct {
	Paths    []string      `brot:"paths"`
	Model    *model.Butter `brot:"model"`
	template *template.Template
}

type ProvidesTemplates interface {
	TemplateFunc() *template.Template
}

func (tl *TemplateLoader) TemplateFunc() *template.Template {
	return tl.template
}

func (tl *TemplateLoader) InitFunc() {
	for _, basePath := range tl.Paths {

		filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				if Warning {
					Warning.Printf("Skipping dir %s, reason is: %s", path, err.Error())
				}
				return nil
			}
			// is it a regular file?
			if info.Mode().IsRegular() {
				var t *template.Template
				base := filepath.Base(path)
				ext := filepath.Ext(base)
				name := base[0 : len(base)-len(ext)]
				if tl.template == nil {
					// first template
					t = template.New(name).Funcs(template.FuncMap{
						"dict": func(values ...interface{}) (map[string]interface{}, error) {
							if len(values)%2 != 0 {
								return nil, errors.New("invalid dict call")
							}
							dict := make(map[string]interface{}, len(values)/2)
							for i := 0; i < len(values); i += 2 {
								key, ok := values[i].(string)
								if !ok {
									return nil, errors.New("dict keys must be strings")
								}
								dict[key] = values[i+1]
							}
							return dict, nil
						},
						"schema": func(obj map[string]string) *model.Schema {
							return tl.Model.Schema(obj["_schema"])
						}})

				} else {
					t = tl.template.New(name)
				}

				content, err := ioutil.ReadFile(path)
				if err != nil {
					if Warning {
						Warning.Printf("Skipping file %s, reason is %s", path, err.Error())
					}
				}

				if result, err := t.Parse(string(content)); err != nil {
					if Warning {
						Warning.Printf("Skipping file %s, reason is %s", path, err.Error())
					}
				} else {
					if Info {
						Info.Printf("Parsed successfully template %s", path)
						tl.template = result
					}
				}

			}
			return nil
		})
	}
}

var _ = di.GlobalScope.Declare(((*TemplateLoader)(nil)))
