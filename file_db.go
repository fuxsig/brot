// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brot

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/fuxsig/brot/di"
)

type FileDB struct {
	Database ProvidesDatabase `brot:"database"`
	Path     string           `brot:"path"`
}

type FileInfo struct {
	Key         string
	Name        string
	ContentType string
	Path        string
	Size        int64
}

func (fdb *FileDB) Info(key string) (*FileInfo, error) {
	fi := &FileInfo{}
	if err := fdb.Database.Get(key, fi); err != nil {
		return nil, err
	}
	return fi, nil
}

func (fdb *FileDB) Image(fi *FileInfo) (image.Image, error) {
	f, err := os.Open(fi.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	switch fi.ContentType {
	case "image/jpeg":
		return jpeg.Decode(f)
	case "image/png":
		return png.Decode(f)
	case "image/gif":
		return gif.Decode(f)
	default:
		return nil, fmt.Errorf("Unsupported format %s", fi.ContentType)
	}
}

func (fdb *FileDB) Store(r io.Reader, name, mimeType, extension string) (*FileInfo, error) {
	f, path, key, err := nextFile(fdb.Path, extension)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	size, err := io.Copy(f, r)
	if err != nil {
		return nil, err
	}

	fi := &FileInfo{Key: key, Name: name, ContentType: mimeType, Path: path, Size: size}
	if err = fdb.Database.Set(key, fi); err != nil {
		return nil, err
	}
	return fi, nil
}

var rand uint32
var randmu sync.Mutex

func reseed() uint32 {
	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
}

func nextSuffix() string {
	randmu.Lock()
	r := rand
	if r == 0 {
		r = reseed()
	}
	r = r*1664525 + 1013904223 // constants from Numerical Recipes
	rand = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}

func nextFile(dir, postfix string) (f *os.File, path string, key string, err error) {
	nconflict := 0
	for i := 0; i < 10000; i++ {
		id := nextSuffix() + postfix
		name := filepath.Join(dir, id)
		f, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
		if os.IsExist(err) {
			if nconflict++; nconflict > 10 {
				randmu.Lock()
				rand = reseed()
				randmu.Unlock()
			}
			continue
		}
		path = name
		key = id
		break
	}
	return
}

var _ = di.GlobalScope.Declare((*FileDB)(nil))
