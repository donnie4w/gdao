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
	XMLName   xml.Name   `xml:"mapper"`
	Namespace string     `xml:"namespace,attr"`
	CrudNodes []CrudNode `xml:",any"`
}

type CrudNode struct {
	XMLName       xml.Name     `xml:""`
	ID            string       `xml:"id,attr"`
	ResultType    string       `xml:"resultType,attr"`
	ParameterType string       `xml:"parameterType,attr,omitempty"`
	Query         string       `xml:",chardata"`
	Dynamics      []DynamicXml `xml:",any"`
}

type DynamicXml struct {
	XMLName         xml.Name     `xml:""`
	Test            string       `xml:"test,attr"`
	Query           string       `xml:",chardata"`
	Collection      string       `xml:"collection,attr"`
	Item            string       `xml:"item,attr"`
	Index           string       `xml:"index,attr"`
	Open            string       `xml:"open,attr,omitempty"`
	Close           string       `xml:"close,attr,omitempty"`
	Separator       string       `xml:"separator,attr,omitempty"`
	Prefix          string       `xml:"prefix,attr,omitempty"`
	Suffix          string       `xml:"suffix,attr,omitempty"`
	PrefixOverrides string       `xml:"prefixOverrides,attr,omitempty"`
	SuffixOverrides string       `xml:"suffixOverrides,attr,omitempty"`
	Child           []DynamicXml `xml:",any"`
}
