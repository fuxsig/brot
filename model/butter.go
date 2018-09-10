// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"

	"github.com/fuxsig/brot/aux"

	"github.com/fuxsig/brot/di"
)

type Butter struct {
	Conn              Connection  `brot:"connection"`
	ConfiguredSchemas []*Schema   `brot:"schemas"`
	Generator         IDGenerator `brot:"generator"`
	schemas           map[string]*Schema
}

func (b *Butter) InitFunc() {
	b.schemas = make(map[string]*Schema, len(b.ConfiguredSchemas))
	for _, schema := range b.ConfiguredSchemas {
		b.schemas[schema.Name] = schema
	}
	for _, schema := range b.ConfiguredSchemas {
		schema.initialize(b)
		b.Conn.BuildIndex(schema)
	}
	//b.conn = b.Db.Open()
}

func (b *Butter) Schema(name string) *Schema {
	result := b.schemas[name]
	return result
}

func (b *Butter) Save(ctx *Context, data map[string]string) (err error) {
	schemaName := data["_schema"]
	// get the schema
	schema, ok := b.schemas[schemaName]
	if !ok {
		return fmt.Errorf("Schema %s is unknown", schemaName)
	}
	// check values
	for key, value := range data {
		if err = schema.CheckValue(key, value); err != nil {
			return
		}
	}
	// save the payload
	err = b.Conn.Save(ctx, data, schema)
	return
}

func (b *Butter) LoadSlice(ctx *Context, ids ...string) ([]map[string]string, error) {
	var (
		me  aux.MultiError
		err error
	)
	result := make([]map[string]string, 0, len(ids))
	var current map[string]string
	for _, id := range ids {
		_, current, err = b.Load(ctx, id)
		if err != nil {
			me.Append(err)
			continue
		}
		result = append(result, current)
	}
	return result, me.ErrorOrNil()
}

func (b *Butter) Load(ctx *Context, id string) (status int, result map[string]string, err error) {
	return b.Conn.Load(ctx, id)
}

func (b *Butter) Search(ctx *Context, schema, query, sort string, offset, num int) ([]map[string]string, int, error) {
	return b.Conn.Search(ctx, schema, query, sort, offset, num)
}

var _ di.ProvidesInit = (*Butter)(nil)
var _ = di.GlobalScope.Declare((*Butter)(nil))
