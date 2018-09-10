// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brot

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/RedisLabs/redisearch-go/redisearch"
	"github.com/fuxsig/brot/di"
	"github.com/fuxsig/brot/model"
	"github.com/garyburd/redigo/redis"
)

// RedisearchHandler is used for the network configuration. Address accepts a single host:port or a comma
// separated host:port,host:port,... value
type RedisearchHandler struct {
	Address string `brot:"address"`
	client  *redisearch.Client
}

// InitFunc initialize opens a new connection to the configured Redis database
func (r *RedisearchHandler) InitFunc() {
	r.client = redisearch.NewClient(r.Address, "vanilla")
	if _, err := r.client.Info(); err != nil {
		schema := redisearch.NewSchema(redisearch.DefaultOptions).
			AddField(redisearch.NewSortableTextField("_schema", 1.0)).
			AddField(redisearch.NewSortableTextField("_owner", 1.0)).
			AddField(redisearch.NewSortableTextField("_label", 2.0)).
			AddField(redisearch.NewSortableNumericField("_created")).
			AddField(redisearch.NewSortableNumericField("_modified")).
			AddField(redisearch.NewTagField("_read")).
			AddField(redisearch.NewTagField("_write")).
			AddField(redisearch.NewTagField("_delete"))
		if err := r.client.CreateIndex(schema); err != nil {
			log.Panicf("Could not create index for %s", r.Address)
		}
	}
}

func (r *RedisearchHandler) BuildIndex(schema *model.Schema) {
	r.client = redisearch.NewClient(r.Address, schema.Plural)
	if _, err := r.client.Info(); err != nil {
		s := redisearch.NewSchema(redisearch.DefaultOptions).
			AddField(redisearch.NewSortableTextField("_schema", 1.0)).
			AddField(redisearch.NewSortableTextField("_owner", 1.0)).
			AddField(redisearch.NewSortableTextField("_label", 2.0)).
			AddField(redisearch.NewSortableNumericField("_created")).
			AddField(redisearch.NewSortableNumericField("_modified")).
			AddField(redisearch.NewTextField("_content")).
			AddField(redisearch.NewTagField("_read")).
			AddField(redisearch.NewTagField("_write")).
			AddField(redisearch.NewTagField("_delete"))
		for _, element := range schema.Elements {
			if !element.Index {
				continue
			}
			name := strings.ToLower(element.Name)
			switch element.Type {
			case "string":
				if element.Sort {
					s.AddField(redisearch.NewSortableTextField(name, 1.0))
				} else {
					s.AddField(redisearch.NewTextField(name))
				}
			case "int":
				if element.Sort {
					s.AddField(redisearch.NewSortableNumericField(name))
				} else {
					s.AddField(redisearch.NewNumericField(name))
				}
			default:
				log.Panicf("ups, not implemented yet")
			}
		}
		if err := r.client.CreateIndex(s); err != nil {
			log.Panicf("Could not create index for %s", r.Address)
		}
	}
}

func (r *RedisearchHandler) Save(ctx *model.Context, data map[string]string, schema *model.Schema) (err error) {
	if ctx == nil {
		ctx = model.Anonymous
	}
	var conn redis.Conn
	conn, err = redis.Dial("tcp", r.Address)
	if err != nil {
		return
	}
	defer conn.Close()
	// check the id
	id := data["_id"]
	create := id == ""
	if create {
		var (
			next   int64
			exists int
		)
		searching := true
		for searching {
			if next, err = redis.Int64(conn.Do("INCR", "brotid")); err != nil {
				return
			}
			id = "brot:" + strconv.FormatInt(next, 36)
			if exists, err = redis.Int(conn.Do("EXISTS", id)); err != nil {
				return
			}
			searching = exists > 0
		}
		data["_id"] = id
	} else {
		// get meta data
		var values []string
		if values, err = redis.Strings(conn.Do("HMGET", id, "_schema", "_write")); err != nil {
			return
		}
		if values[0] == "" && values[1] == "" {
			create = true
		} else {
			// check schema
			if values[0] != data["_schema"] {
				return fmt.Errorf("wrong schema %s, expected %s", data["_schema"], values[0])
			}
			// check write groups
			groups := strings.Split(values[1], ",")
			if !ctx.MemberOf(groups...) {
				return fmt.Errorf("access denied")
			}
		}
	}
	doc := redisearch.NewDocument(id, 1.0)
	read := data["_read"]
	write := data["_write"]
	delete := data["_delete"]
	if create {
		// fix rights if necessary
		if read == "" {
			read = "all"
			data["_read"] = read
		}

		if write == "" {
			write = "all"
			data["_write"] = write
		}

		if delete == "" {
			delete = "all"
			data["_delete"] = delete
		}
	}

	if create {
		data["_created"] = strconv.FormatInt(time.Now().Unix(), 10)
	}
	data["_modified"] = strconv.FormatInt(time.Now().Unix(), 10)

	for key, value := range data {
		doc.Set(key, value)
	}
	c := redisearch.NewClient(r.Address, schema.Plural)
	if create {
		if err = c.Index(doc); err != nil {
			return
		}
		_, err = conn.Do("FT.ADDHASH", "vanilla", id, 1.0)
	} else {

		// update
		if err = c.IndexOptions(redisearch.IndexingOptions{Replace: true, Partial: true}, doc); err != nil {
			return
		}
		_, err = conn.Do("FT.ADDHASH", "vanilla", id, 1.0, "REPLACE")
	}
	return
}

func (r *RedisearchHandler) Load(ctx *model.Context, id string) (int, map[string]string, error) {
	if ctx == nil {
		ctx = model.Anonymous
	}
	conn, err := redis.Dial("tcp", r.Address)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	defer conn.Close()
	var result map[string]string
	result, err = redis.StringMap(conn.Do("HGETALL", id))
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	read := result["_read"]
	groups := strings.Split(read, ",")
	if !ctx.MemberOf(groups...) {
		return http.StatusForbidden, nil, fmt.Errorf("forbidden")
	}
	return http.StatusOK, result, nil
}

func (r *RedisearchHandler) Delete(ctx *model.Context, id string) (bool, error) {
	return r.client.Delete(id, true)
}

func (r *RedisearchHandler) Search(ctx *model.Context, index, query, sort string, offset, num int) ([]map[string]string, int, error) {
	if ctx == nil {
		ctx = model.Anonymous
	}
	conn, err := redis.Dial("tcp", r.Address)
	if err != nil {
		return nil, 0, err
	}
	defer conn.Close()

	var buf bytes.Buffer
	buf.WriteString(query)
	buf.WriteString(" @_read:{")
	buf.WriteString(ctx.GroupsString())
	buf.WriteString("}")

	var values []interface{}
	args := redis.Args{strings.ToLower(index), buf.String(), "LIMIT", offset, num}
	if sort != "" {
		asc := !strings.HasPrefix(sort, "-")
		var order string
		if asc {
			if strings.HasPrefix(sort, "+") {
				sort = sort[1:]
			}
			order = "ASC"
		} else {
			sort = sort[1:]
			order = "DESC"
		}
		args = append(args, "SORTBY", sort, order)

	}
	values, err = redis.Values(conn.Do("FT.SEARCH", args...))
	if err != nil {
		return nil, 0, err
	}
	var total int
	if total, err = redis.Int(values[0], nil); err != nil {
		return nil, 0, err
	}

	result := make([]map[string]string, (len(values)-1)/2)

	for i, j := 0, 2; i < len(result); i++ {
		data, ok := values[j].([]interface{})
		j += 2
		if !ok {
			return nil, total, fmt.Errorf("internal error, expected []interface{}")
		}
		result[i], err = redis.StringMap(data, nil)
		if err != nil {
			return nil, 0, err
		}
	}
	return result, total, nil
}

var _ di.ProvidesInit = (*RedisearchHandler)(nil)
var _ model.Connection = (*RedisearchHandler)(nil)
var _ = di.GlobalScope.Declare((*RedisearchHandler)(nil))
