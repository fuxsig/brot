// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package di

import (
	"strconv"
)

// GetInt64 accepts a value and tries to return an int64 representation.
// If a conversion is not possible, then this function returns an error.
func GetInt64(value interface{}) (result int64, err error) {
	switch i := value.(type) {
	case int64:
		result = i
	case float64:
		result = int64(i)
	case float32:
		result = int64(i)
	case int:
		result = int64(i)
	case int32:
		result = int64(i)
	case int16:
		result = int64(i)
	case int8:
		result = int64(i)
	case uint64:
		result = int64(i)
	case uint:
		result = int64(i)
	case uint32:
		result = int64(i)
	case uint16:
		result = int64(i)
	case uint8:
		result = int64(i)
	case string:
		result, err = strconv.ParseInt(i, 10, 64)
	default:
		err = ErrUnexpectedType
	}
	return
}

// GetInt accepts a value and tries to return an int representation.
// If a conversion is not possible, then this function returns an error.
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
	case uint:
		result = int(i)
	case uint32:
		result = int(i)
	case uint16:
		result = int(i)
	case uint8:
		result = int(i)
	case string:
		var aux int64
		if aux, err = strconv.ParseInt(i, 10, 64); err == nil {
			result = int(aux)
		}
	default:
		err = ErrUnexpectedType
	}
	return
}

// GetInt32 accepts a value and tries to return an int32 representation.
// If a conversion is not possible, then this function returns an error.
func GetInt32(value interface{}) (result int32, err error) {
	switch i := value.(type) {
	case int32:
		result = i
	case float64:
		result = int32(i)
	case float32:
		result = int32(i)
	case int64:
		result = int32(i)
	case int:
		result = int32(i)
	case int16:
		result = int32(i)
	case int8:
		result = int32(i)
	case uint64:
		result = int32(i)
	case uint:
		result = int32(i)
	case uint32:
		result = int32(i)
	case uint16:
		result = int32(i)
	case uint8:
		result = int32(i)
	case string:
		var aux int64
		if aux, err = strconv.ParseInt(i, 10, 64); err == nil {
			result = int32(aux)
		}
	default:
		err = ErrUnexpectedType
	}
	return
}

// GetInt16 accepts a value and tries to return an int16 representation.
// If a conversion is not possible, then this function returns an error.
func GetInt16(value interface{}) (result int16, err error) {
	switch i := value.(type) {
	case int16:
		result = i
	case float64:
		result = int16(i)
	case float32:
		result = int16(i)
	case int64:
		result = int16(i)
	case int:
		result = int16(i)
	case int32:
		result = int16(i)
	case int8:
		result = int16(i)
	case uint64:
		result = int16(i)
	case uint:
		result = int16(i)
	case uint32:
		result = int16(i)
	case uint16:
		result = int16(i)
	case uint8:
		result = int16(i)
	case string:
		var aux int64
		if aux, err = strconv.ParseInt(i, 10, 64); err == nil {
			result = int16(aux)
		}
	default:
		err = ErrUnexpectedType
	}
	return
}

// GetInt8 accepts a value and tries to return an int8 representation.
// If a conversion is not possible, then this function returns an error.
func GetInt8(value interface{}) (result int8, err error) {
	switch i := value.(type) {
	case int8:
		result = i
	case float64:
		result = int8(i)
	case float32:
		result = int8(i)
	case int64:
		result = int8(i)
	case int:
		result = int8(i)
	case int32:
		result = int8(i)
	case int16:
		result = int8(i)
	case uint64:
		result = int8(i)
	case uint:
		result = int8(i)
	case uint32:
		result = int8(i)
	case uint16:
		result = int8(i)
	case uint8:
		result = int8(i)
	case string:
		var aux int64
		if aux, err = strconv.ParseInt(i, 10, 64); err == nil {
			result = int8(aux)
		}
	default:
		err = ErrUnexpectedType
	}
	return
}
