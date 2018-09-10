// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brot

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ViewRegistry is the registry for views
type ViewRegistry map[string]View

// Init initilaizes the view registry
func (vr ViewRegistry) Init(paths []string) {
	stack := make([]os.FileInfo, 0)

	for _, path := range paths {

		// open the file
		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		// read directory
		files, err := f.Readdir(-1)
		f.Close()
		if err != nil {
			log.Fatal(err)
		}
		// iterate over all files in directory
		for _, file := range files {
			if file.IsDir() {
				stack = append(stack, file)
			} else {
				fullPath := filepath.Join(path, file.Name())
				content, err := ioutil.ReadFile(fullPath)
				if err != nil {
					log.Fatal(err)
				}
				name := file.Name()
				ext := filepath.Ext(name)
				viewName := name[:len(name)-len(ext)]
				viewType := "unknown"
				if index := strings.Index(viewName, "@"); index >= 0 {
					viewType = viewName[index+1:]
					viewName = viewName[:index]

				}
				qname := viewName + "@" + viewType

				switch strings.ToLower(ext) {
				case ".html":
					t, _ := template.New(qname).Parse(string(content))
					vr[qname] = View{t}
				}

			}
		}

	}
}
