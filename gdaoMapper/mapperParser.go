// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoMapper

import (
	"encoding/xml"
	"fmt"
	"github.com/donnie4w/gdao/base"
	. "github.com/donnie4w/gofer/hashmap"
	"os"
	"strings"
)

type mapperParser struct {
	mapper          *MapL[string, *paramBean]
	namespaceMapper *MapL[string, []string]
}

func newMapperParser() *mapperParser {
	return &mapperParser{mapper: NewMapL[string, *paramBean](), namespaceMapper: NewMapL[string, []string]()}
}

var mapperparser *mapperParser

func init() {
	mapperparser = newMapperParser()
	base.GetMapperIds = mapperparser.getMapperIds
	base.HasMapperId = mapperparser.hasMapperId
}

func (m *mapperParser) getMapperIds(namespace string) []string {
	if s, ok := m.namespaceMapper.Get(namespace); ok {
		return s
	}
	return nil
}

func (m *mapperParser) hasMapperId(mapperId string) bool {
	return m.mapper.Has(mapperId)
}

func (m *mapperParser) hasMapper() bool {
	return m.mapper.Len() > 0 || m.namespaceMapper.Len() > 0
}

func (m *mapperParser) parser(xmlpath string) (err error) {
	file, err := os.Open(xmlpath)
	if err != nil {
		return err
	}
	defer file.Close()
	var root rootElement
	decoder := xml.NewDecoder(file)
	if err = decoder.Decode(&root); err != nil {
		return err
	}
	if root.XMLName.Local == "configuration" {
		file.Seek(0, 0)
		var configuration Configuration
		if err = decoder.Decode(&configuration); err != nil {
			return err
		}
		for _, mapper := range configuration.Mappers.MapperList {
			m.parseMapper(mapper.Resource)
		}
	} else if root.XMLName.Local == "mappers" {
		file.Seek(0, 0)
		var mappers mappers
		if err = decoder.Decode(&mappers); err != nil {
			return err
		}
		for _, mapper := range mappers.MapperList {
			m.parseMapper(mapper.Resource)
		}
	} else if root.XMLName.Local == "mapper" {
		m.parseMapper(xmlpath)
	} else {
		return fmt.Errorf("Unknown root element: %s\n", root.XMLName.Local)
	}
	return nil
}

func (m *mapperParser) parseMapper(xmlpath string) (err error) {
	file, err := os.Open(xmlpath)
	if err != nil {
		return
	}
	defer file.Close()

	var mapper mapper
	decoder := xml.NewDecoder(file)
	if err = decoder.Decode(&mapper); err != nil {
		return err
	}

	for _, crudNode := range mapper.CrudNodes {
		pb := newParamBean(mapper.Namespace, crudNode.ID, crudNode.XMLName.Local, crudNode.Query, crudNode.ParameterType, crudNode.ResultType)
		m.mapperAdd(mapper.Namespace, crudNode.ID, pb)
		if node := sqlnode(crudNode); node != nil {
			pb.sqlNode = node
		}
	}
	return nil
}

func sqlnode(crudnode CrudNode) (r sqlNode) {
	if len(crudnode.Dynamics) == 0 {
		return nil
	}
	switch strings.ToLower(crudnode.XMLName.Local) {
	case "select", "update", "delete", "insert":
		r = newCrudNode(strings.TrimSpace(crudnode.Query))
	default:
		panic("Unsupport tag: " + crudnode.XMLName.Local)
	}
	for _, dynamicNode := range crudnode.Dynamics {
		if node := getSqlNode(dynamicNode); node != nil {
			r.addSqlNode(node)
		}
	}
	return
}

func getSqlNode(dynamicNode DynamicXml) (r sqlNode) {
	trim := strings.TrimSpace
	switch strings.ToLower(dynamicNode.XMLName.Local) {
	case "if":
		r = newIfNode(trim(dynamicNode.Test), trim(dynamicNode.Query))
	case "where":
		r = newWhereNode(trim(dynamicNode.Query))
	case "trim":
		r = newTrimNode(dynamicNode.Prefix, dynamicNode.Suffix, dynamicNode.PrefixOverrides, dynamicNode.SuffixOverrides)
	case "set":
		r = newSetNode()
	case "foreach":
		r = newForeach(trim(dynamicNode.Query), dynamicNode.Collection, dynamicNode.Item, dynamicNode.Index, dynamicNode.Open, dynamicNode.Close, dynamicNode.Separator)
	case "choose":
		r = newChooseNode(trim(dynamicNode.Test), trim(dynamicNode.Query))
	case "when":
		r = newWhenNode(trim(dynamicNode.Test), trim(dynamicNode.Query))
	case "otherwise":
		r = newOtherWiseNode(trim(dynamicNode.Test), trim(dynamicNode.Query))
	default:
		panic("Unsupport tag: " + dynamicNode.XMLName.Local)
	}
	for _, dn := range dynamicNode.Child {
		r.addSqlNode(getSqlNode(dn))
	}
	return
}

func (m *mapperParser) mapperAdd(namespace, id string, pb *paramBean) {
	if id == "" {
		panic("Mapper id is empty . [namespace]:" + namespace)
	}
	mapperId := namespace + "." + id
	if namespace == "" {
		mapperId = id
	}
	if m.mapper.Has(mapperId) {
		panic(fmt.Sprintf("namespace and id are defined repeatedly: [namespace:%s][id:%s]", namespace, id))
	}
	if pb.sqltype == _SELECT {
		m.namespaceMapperAdd(namespace, id)
	}
	m.mapper.Put(mapperId, pb)
}

func (m *mapperParser) getParamBean(mapperId string) (*paramBean, bool) {
	return m.mapper.Get(mapperId)
}

func (m *mapperParser) namespaceMapperAdd(namespace, id string) {
	if ids, _ := m.namespaceMapper.Get(namespace); ids != nil {
		m.namespaceMapper.Put(namespace, append(ids, id))
	} else {
		m.namespaceMapper.Put(namespace, []string{id})
	}
}

// Builder parses an XML mapping file located at the specified path to configure the application.
// Parameters:
//
//	xmlPath: A string representing the file path of the XML mapping file.
//
// This function is responsible for reading the XML file and extracting configuration information for the application.
func Builder(xmlPath string) error {
	return mapperparser.parser(xmlPath)
}
