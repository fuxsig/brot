// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package di

import (
	"errors"
	"math"
	"strconv"
)

// ErrUnexpectedType is returned if a different type was expected.
var ErrUnexpectedType = errors.New("Non-numeric type could not be converted")

// GetBool accepts a value and tries to return a boolean represenation.
// If a conversion is not possible, then this function returns an error.
func GetBool(value interface{}) (result bool, err error) {
	switch b := value.(type) {
	case nil:
		result = false
	case bool:
		result = b
	case string:
		result, err = strconv.ParseBool(b)
	case int64:
		result = b != 0
	case int32:
		result = b != 0
	case int16:
		result = b != 0
	case int8:
		result = b != 0
	case int:
		result = b != 0
	case uint64:
		result = b != 0
	case uint32:
		result = b != 0
	case uint16:
		result = b != 0
	case uint8:
		result = b != 0
	case uint:
		result = b != 0
	default:
		err = ErrUnexpectedType
	}
	return
}

// GetString accepts a value and tries to return a string representation.
// If a conversion is not possible, then this function returns an error.
func GetString(value interface{}) (result string, err error) {
	switch i := value.(type) {
	case nil:
		result = ""
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
		err = ErrUnexpectedType
	}
	return
}

// GetFloat64 accepts a value and tries to return a float64 representation.
// If a conversion is not possible, then this function returns an error.
func GetFloat64(value interface{}) (result float64, err error) {
	switch i := value.(type) {
	case nil:
		result = 0
	case float64:
		result = i
	case float32:
		result = float64(i)
	case int64:
		result = float64(i)
	case int:
		result = float64(i)
	case int32:
		result = float64(i)
	case int16:
		result = float64(i)
	case int8:
		result = float64(i)
	case uint64:
		result = float64(i)
	case uint:
		result = float64(i)
	case uint32:
		result = float64(i)
	case uint16:
		result = float64(i)
	case uint8:
		result = float64(i)
	case string:
		result, err = strconv.ParseFloat(i, 64)
	default:
		result, err = math.NaN(), ErrUnexpectedType
	}
	return
}

// GetFloat32 accepts a value and tries to return a float64 representation.
// If a conversion is not possible, then this function returns an error.
func GetFloat32(value interface{}) (result float32, err error) {
	switch i := value.(type) {
	case nil:
		result = 0
	case float32:
		result = i
	case float64:
		result = float32(i)
	case int64:
		result = float32(i)
	case int:
		result = float32(i)
	case int32:
		result = float32(i)
	case int16:
		result = float32(i)
	case int8:
		result = float32(i)
	case uint64:
		result = float32(i)
	case uint:
		result = float32(i)
	case uint32:
		result = float32(i)
	case uint16:
		result = float32(i)
	case uint8:
		result = float32(i)
	case string:
		var aux float64
		if aux, err = strconv.ParseFloat(i, 64); err == nil {
			result = float32(aux)
		}
	default:
		result, err = 0, ErrUnexpectedType
	}
	return
}
