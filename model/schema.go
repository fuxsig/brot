// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"strings"
)

// Schema describes the structure of an object.
type Schema struct {
	Name     string     `brot:"name"`
	Plural   string     `brot:"plural"`
	Elements []*Element `brot:"elements"`
	elements map[string]*Element
	butter   *Butter
	Indexed  []string
}

func (s *Schema) initialize(butter *Butter) error {
	s.butter = butter
	schemas := butter.schemas

	if s.Plural == "" {
		s.Plural = strings.ToLower(s.Name) + "s"
	} else {
		s.Plural = strings.ToLower(s.Plural)
	}

	// initialize elements
	indexed := make([]string, 0, len(s.Elements))
	s.elements = make(map[string]*Element, len(s.Elements))
	for _, element := range s.Elements {
		err := element.initialize(schemas)
		if err != nil {
			return err
		}
		s.elements[element.Name] = element
		if element.Index {
			indexed = append(indexed, element.Name)
		}
	}
	s.Indexed = make([]string, len(indexed))
	copy(s.Indexed, indexed)
	return nil
}

// CheckValue accepts a name value pair and checks the validity of the value.
func (s *Schema) CheckValue(name, value string) (err error) {
	e := s.elements[name]
	if e == nil {
		switch name {
		case "_schema":
			if s.butter.schemas[value] == nil {
				err = fmt.Errorf("schema value %s is not registered", value)
			}
		case "_label", "_id":
			if len(value) > 256 {
				err = fmt.Errorf("%s has more than 256 bytes", name)
			}
		default:
			err = fmt.Errorf("schema %s does not contain element with name %s", s.Name, name)
		}
		return
	}
	if e.slice {
		var err error
		values := strings.Split(value, ",")
		for _, current := range values {
			err = e.check(current)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return e.check(value)
}
