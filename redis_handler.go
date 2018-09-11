// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brot

import (
	"bytes"
	"encoding/gob"
	"log"
	"reflect"
	"time"

	"github.com/fuxsig/brot/di"

	"github.com/mediocregopher/radix.v2/pool"
)

// RedisHandler provides a connection to a Redis database
type RedisHandler struct {
	Network string `brot:"network"`
	Address string `brot:"address"`
	Size    int    `brot:"size"`
	pool    *pool.Pool
}

// InitFunc initialize opens a new connection to the configured Redis database
func (r *RedisHandler) InitFunc() (err error) {
	r.pool, err = pool.New(r.Network, r.Address, r.Size)
	if err != nil {
		// handle err
		log.Printf("Error in RedisHandler: %s", err.Error())
		panic("No Database!!!")
	}
	return
}

func (r *RedisHandler) Retry() bool {
	time.Sleep(2 * time.Second)
	return true
}

func (r *RedisHandler) Get(name string, object interface{}) error {
	value := reflect.ValueOf(object)
	conn, err := r.pool.Get()
	if err != nil {
		return err
	}
	defer r.pool.Put(conn)
	data, err := conn.Cmd("GET", name).Bytes()
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(bytes.NewReader(data))
	return dec.DecodeValue(value)
}

func (r *RedisHandler) GetPtr(name, typeName string) interface{} {
	val, err := di.GlobalScope.Create(typeName)
	if err != nil {
		return nil
	}

	conn, err := r.pool.Get()
	if err != nil {
		return nil
	}
	defer r.pool.Put(conn)
	data, err := conn.Cmd("GET", name).Bytes()
	if err != nil {
		return nil
	}
	dec := gob.NewDecoder(bytes.NewReader(data))
	dec.DecodeValue(val)
	return val.Interface()
}

func (r *RedisHandler) Set(name string, object interface{}) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(object); err != nil {
		return err
	}

	conn, err := r.pool.Get()
	if err != nil {
		return err
	}
	defer r.pool.Put(conn)

	resp := conn.Cmd("SET", name, b.Bytes())
	if resp.Err != nil {
		return resp.Err
	}
	return nil
}

var _ di.ProvidesInit = (*RedisHandler)(nil)
var _ ProvidesDatabase = (*RedisHandler)(nil)
var _ = di.GlobalScope.Declare((*RedisHandler)(nil))

//HandlerDefaultRegistry.Register()
