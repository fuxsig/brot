// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package di

import (
	"strconv"
)

// GetUint64 accepts a value and tries to return an uint64 representation.
// If a conversion is not possible, then this function returns an error.
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
	case int:
		result = uint64(i)
	case int32:
		result = uint64(i)
	case int16:
		result = uint64(i)
	case int8:
		result = uint64(i)
	case uint:
		result = uint64(i)
	case uint32:
		result = uint64(i)
	case uint16:
		result = uint64(i)
	case uint8:
		result = uint64(i)
	case string:
		result, err = strconv.ParseUint(i, 10, 64)
	default:
		err = ErrUnexpectedType
	}
	return
}

// GetUint accepts a value and tries to return an uint representation.
// If a conversion is not possible, then this function returns an error.
func GetUint(value interface{}) (result uint, err error) {
	switch i := value.(type) {
	case uint:
		result = i
	case float64:
		result = uint(i)
	case float32:
		result = uint(i)
	case int64:
		result = uint(i)
	case int:
		result = uint(i)
	case int32:
		result = uint(i)
	case int16:
		result = uint(i)
	case int8:
		result = uint(i)
	case uint64:
		result = uint(i)
	case uint32:
		result = uint(i)
	case uint16:
		result = uint(i)
	case uint8:
		result = uint(i)
	case string:
		var aux uint64
		if aux, err = strconv.ParseUint(i, 10, 64); err != nil {
			result = uint(aux)
		}
	default:
		err = ErrUnexpectedType
	}
	return
}

// GetUint32 accepts a value and tries to return an uint32 representation.
// If a conversion is not possible, then this function returns an error.
func GetUint32(value interface{}) (result uint32, err error) {
	switch i := value.(type) {
	case uint32:
		result = i
	case float64:
		result = uint32(i)
	case float32:
		result = uint32(i)
	case int64:
		result = uint32(i)
	case int:
		result = uint32(i)
	case int32:
		result = uint32(i)
	case int16:
		result = uint32(i)
	case int8:
		result = uint32(i)
	case uint64:
		result = uint32(i)
	case uint:
		result = uint32(i)
	case uint16:
		result = uint32(i)
	case uint8:
		result = uint32(i)
	case string:
		var aux uint64
		if aux, err = strconv.ParseUint(i, 10, 64); err != nil {
			result = uint32(aux)
		}
	default:
		err = ErrUnexpectedType
	}
	return
}

// GetUint16 accepts a value and tries to return an uint16 representation.
// If a conversion is not possible, then this function returns an error.
func GetUint16(value interface{}) (result uint16, err error) {
	switch i := value.(type) {
	case uint16:
		result = i
	case float64:
		result = uint16(i)
	case float32:
		result = uint16(i)
	case int64:
		result = uint16(i)
	case int:
		result = uint16(i)
	case int32:
		result = uint16(i)
	case int16:
		result = uint16(i)
	case int8:
		result = uint16(i)
	case uint64:
		result = uint16(i)
	case uint:
		result = uint16(i)
	case uint32:
		result = uint16(i)
	case uint8:
		result = uint16(i)
	case string:
		var aux uint64
		if aux, err = strconv.ParseUint(i, 10, 64); err != nil {
			result = uint16(aux)
		}
	default:
		err = ErrUnexpectedType
	}
	return
}

// GetUint8 accepts a value and tries to return an uint8 representation.
// If a conversion is not possible, then this function returns an error.
func GetUint8(value interface{}) (result uint8, err error) {
	switch i := value.(type) {
	case uint8:
		result = i
	case float64:
		result = uint8(i)
	case float32:
		result = uint8(i)
	case int64:
		result = uint8(i)
	case int:
		result = uint8(i)
	case int32:
		result = uint8(i)
	case int16:
		result = uint8(i)
	case int8:
		result = uint8(i)
	case uint64:
		result = uint8(i)
	case uint:
		result = uint8(i)
	case uint32:
		result = uint8(i)
	case uint16:
		result = uint8(i)
	case string:
		var aux uint64
		if aux, err = strconv.ParseUint(i, 10, 64); err != nil {
			result = uint8(aux)
		}
	default:
		err = ErrUnexpectedType
	}
	return
}
