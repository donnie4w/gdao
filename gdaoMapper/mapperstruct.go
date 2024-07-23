// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoMapper

import (
	"encoding/xml"
)

type rootElement struct {
	XMLName xml.Name
}

type Configuration struct {
	XMLName xml.Name `xml:"configuration"`
	Mappers *mappers `xml:"mappers"`
}

type mappers struct {
	XMLName    xml.Name     `xml:"mappers"`
	MapperList []mapperCell `xml:"mapper"`
}

// mapper represents a mapper element in the XML file
type mapperCell struct {
	Resource string `xml:"resource,attr"`
}

// mapper represents the root element of the  mapper file
type mapper struct {
	XMLName   xml.Name  `xml:"mapper"`
	Namespace string    `xml:"namespace,attr"`
	Selects   []select_ `xml:"select"`
	Inserts   []insert  `xml:"insert"`
	Updates   []update  `xml:"update"`
	Deletes   []delete  `xml:"delete"`
}

// Select represents a select statement in the mapper file
type select_ struct {
	ID            string `xml:"id,attr"`
	ResultType    string `xml:"resultType,attr"`
	ParameterType string `xml:"parameterType,attr,omitempty"`
	Query         string `xml:",chardata"`
}

// Insert represents an insert statement in the mapper file
type insert struct {
	ID            string `xml:"id,attr"`
	ParameterType string `xml:"parameterType,attr"`
	Query         string `xml:",chardata"`
}

// Update represents an update statement in the mapper file
type update struct {
	ID            string `xml:"id,attr"`
	ParameterType string `xml:"parameterType,attr"`
	Query         string `xml:",chardata"`
}

// Delete represents a delete statement in the mapper file
type delete struct {
	ID            string `xml:"id,attr"`
	ParameterType string `xml:"parameterType,attr"`
	Query         string `xml:",chardata"`
}
