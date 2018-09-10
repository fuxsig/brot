// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brot

import "html/template"

// View renders data to a specific format
type View struct {
	template *template.Template
}

// Create creates a new template
func Create(name, dt, code string) *View {
	qname := name + "@" + dt
	t, _ := template.New(qname).Parse(code)
	return &View{t}
}
