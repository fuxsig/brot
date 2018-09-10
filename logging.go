// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brot

import (
	"log"
)

type Log bool

func (l Log) Printf(format string, v ...interface{}) Log {
	if l {
		log.Printf(format, v...)
	}
	return l
}

var Info Log = true
var Debug Log = true
var Warning Log = true
