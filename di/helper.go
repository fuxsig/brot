// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package di

import (
	"errors"
	"math"
	"strconv"
)

var errUnexpectedType = errors.New("Non-numeric type could not be converted to float")

func GetBool(value interface{}) (result bool, err error) {
	switch b := value.(type) {
	case bool:
		result = b
	case string:
		result, err = strconv.ParseBool(b)
	default:
		err = errUnexpectedType
	}
	return
}

func GetString(value interface{}) (result string, err error) {
	switch i := value.(type) {
	case string:
		result = i
	case float64:
		result = strconv.FormatFloat(i, 'f', -1, 64)
	case float32:
		result = strconv.FormatFloat(float64(i), 'f', -1, 32)
	case int64:
		result = strconv.FormatInt(i, 10)
	case int32:
		result = strconv.FormatInt(int64(i), 10)
	case int:
		result = strconv.FormatInt(int64(i), 10)
	case int16:
		result = strconv.FormatInt(int64(i), 10)
	case int8:
		result = strconv.FormatInt(int64(i), 10)
	case uint64:
		result = strconv.FormatUint(i, 10)
	case uint32:
		result = strconv.FormatUint(uint64(i), 10)
	case uint:
		result = strconv.FormatUint(uint64(i), 10)
	case uint16:
		result = strconv.FormatUint(uint64(i), 10)
	case uint8:
		result = strconv.FormatUint(uint64(i), 10)
	case bool:
		result = strconv.FormatBool(i)
	default:
		err = errUnexpectedType
	}
	return
}

func GetFloat64(value interface{}) (result float64, err error) {
	switch i := value.(type) {
	case float64:
		result = i
	case float32:
		result = float64(i)
	case int64:
		result = float64(i)
	case int32:
		result = float64(i)
	case int:
		result = float64(i)
	case uint64:
		result = float64(i)
	case uint32:
		result = float64(i)
	case uint:
		result = float64(i)
	case string:
		result, err = strconv.ParseFloat(i, 64)
	default:
		result, err = math.NaN(), errUnexpectedType
	}
	return
}

func GetInt64(value interface{}) (result int64, err error) {
	switch i := value.(type) {
	case int64:
		result = i
	case float64:
		result = int64(i)
	case float32:
		result = int64(i)
	case int32:
		result = int64(i)
	case int:
		result = int64(i)
	case int16:
		result = int64(i)
	case int8:
		result = int64(i)
	case uint64:
		result = int64(i)
	case uint32:
		result = int64(i)
	case uint:
		result = int64(i)
	case uint16:
		result = int64(i)
	case uint8:
		result = int64(i)
	case string:
		result, err = strconv.ParseInt(i, 10, 64)
	default:
		err = errUnexpectedType
	}
	return
}

func GetInt(value interface{}) (result int, err error) {
	switch i := value.(type) {
	case int:
		result = i
	case float64:
		result = int(i)
	case float32:
		result = int(i)
	case int64:
		result = int(i)
	case int32:
		result = int(i)
	case int16:
		result = int(i)
	case int8:
		result = int(i)
	case uint64:
		result = int(i)
	case uint32:
		result = int(i)
	case uint:
		result = int(i)
	case uint16:
		result = int(i)
	case uint8:
		result = int(i)
	case string:
		var i64 int64
		if i64, err = strconv.ParseInt(i, 10, 32); err == nil {
			result = int(i64)
		}
	default:
		err = errUnexpectedType
	}
	return
}

func GetUint64(value interface{}) (result uint64, err error) {
	switch i := value.(type) {
	case uint64:
		result = i
	case float64:
		result = uint64(i)
	case float32:
		result = uint64(i)
	case int64:
		result = uint64(i)
	case int32:
		result = uint64(i)
	case int:
		result = uint64(i)
	case int16:
		result = uint64(i)
	case int8:
		result = uint64(i)
	case uint32:
		result = uint64(i)
	case uint:
		result = uint64(i)
	case uint16:
		result = uint64(i)
	case uint8:
		result = uint64(i)
	case string:
		result, err = strconv.ParseUint(i, 10, 64)
	default:
		err = errUnexpectedType
	}
	return
}
