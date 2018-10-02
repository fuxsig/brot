// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package aux

import (
	"bytes"
	"errors"
)

type MultiError struct {
	Errors []error
}

func (e *MultiError) Error() string {
	var buffer bytes.Buffer
	l := len(e.Errors) - 1
	if l > 0 {
		buffer.WriteString("Multiple errors: ")
	}

	for i, act := range e.Errors {
		buffer.WriteString(act.Error())
		if i < l {
			buffer.WriteString(", ")
		}
	}
	return buffer.String()
}

func (e *MultiError) Append(err error) *MultiError {
	if err != nil {
		e.Errors = append(e.Errors, err)
	}
	return e
}

func (e *MultiError) AppendString(message string) *MultiError {
	e.Append(errors.New(message))
	return e
}

func (e *MultiError) Merge(err error) {
	if err, ok := err.(*MultiError); ok {
		e.Errors = append(e.Errors, err.Errors...)
	}
}

func (e *MultiError) ErrorOrNil() error {
	if e.Errors == nil {
		return nil
	}
	return e
}
