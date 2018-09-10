// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"sort"
	"strings"
)

type Context struct {
	User   string
	groups []string
	grpstr string
}

var Anonymous = new(Context).SetGroups("all")

func (c *Context) SetGroups(groups ...string) *Context {
	sort.Slice(groups, func(i, j int) bool { return groups[i] < groups[j] })
	c.groups = groups
	return c
}

func (c *Context) GroupsString() string {
	if c.grpstr == "" {
		c.grpstr = strings.Join(c.groups, " | ")
	}
	return c.grpstr
}

func (c *Context) MemberOf(groups ...string) bool {
	l := len(c.groups)
	for i := range groups {
		if sort.SearchStrings(c.groups, groups[i]) < l {
			return true
		}
	}
	return false
}

type Connection interface {
	Load(ctx *Context, id string) (int, map[string]string, error)
	Save(*Context, map[string]string, *Schema) error
	Delete(ctx *Context, id string) (bool, error)

	Search(ctx *Context, schema, query, sort string, offset, num int) ([]map[string]string, int, error)

	BuildIndex(schema *Schema)
}

type Database interface {
	Open() (Connection, error)
}

type IDGenerator interface {
	Next() string
}
