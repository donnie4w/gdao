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
	"github.com/donnie4w/gdao/util"
	. "github.com/donnie4w/gofer/hashmap"
	"os"
	"reflect"
	"regexp"
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

	for _, s := range mapper.Selects {
		pb := newParamBean(mapper.Namespace, s.ID, "select", s.Query, s.ParameterType, s.ResultType)
		m.mapperAdd(mapper.Namespace, s.ID, pb)
	}

	for _, ins := range mapper.Inserts {
		pb := newParamBean(mapper.Namespace, ins.ID, "insert", ins.Query, ins.ParameterType, "")
		m.mapperAdd(mapper.Namespace, ins.ID, pb)
	}

	for _, upd := range mapper.Updates {
		pb := newParamBean(mapper.Namespace, upd.ID, "update", upd.Query, upd.ParameterType, "")
		m.mapperAdd(mapper.Namespace, upd.ID, pb)
	}

	for _, del := range mapper.Deletes {
		pb := newParamBean(mapper.Namespace, del.ID, "delete", del.Query, del.ParameterType, "")
		m.mapperAdd(mapper.Namespace, del.ID, pb)
	}
	return nil
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

type paramBean struct {
	namespace      string
	id             string
	sqltype        sqlType
	sql            string
	parameterNames []string
	inputType      string
	outputType     string
}

func newParamBean(namespace, id, sqltype, sql, inputType, outputType string) (r *paramBean) {
	r = &paramBean{namespace: namespace, id: id, inputType: inputType, outputType: outputType}
	switch sqltype {
	case "select":
		r.sqltype = _SELECT
	case "insert":
		r.sqltype = _INSERT
	case "update":
		r.sqltype = _UPDATE
	case "delete":
		r.sqltype = _DELETE
	}
	r.sql, r.parameterNames = parseSql(sql)
	return
}

func parseSql(sqlstr string) (sql string, parameterNames []string) {
	pattern := `#\{(\w+)\}`
	re := regexp.MustCompile(pattern)
	sql = re.ReplaceAllStringFunc(sqlstr, func(match string) string {
		parameterNames = append(parameterNames, strings.Trim(match[2:len(match)-1], " "))
		return "?"
	})
	return strings.TrimSpace(sql), parameterNames
}

func (p *paramBean) err_num_no_match(expectvalue, gotvalue int) error {
	return fmt.Errorf("the parameter number does not match the configuration:Expected %d but got %d. [namespace:%s][mapper id:%s]", expectvalue, gotvalue, p.namespace, p.id)
}

func (p *paramBean) err_invalid_parameter(getType string) error {
	return fmt.Errorf("invalid parameter type: Expected %s but get %s [namespace:%s][mapper id:%s]", p.inputType, getType, p.namespace, p.id)
}

func (p *paramBean) err_params_no_found(param string) error {
	return fmt.Errorf("parameter not found:  %s  [namespace:%s][mapper id:%s]", param, p.namespace, p.id)
}

func (p *paramBean) setParameter(parameter any) (args []any, err error) {
	defer util.Recover(&err)
	if p.inputType == "" || parameter == nil || len(p.parameterNames) == 0 {
		return args, nil
	}
	if isDBType(p.inputType) {
		if len(p.parameterNames) == 1 {
			return []any{parameter}, nil
		} else {
			return nil, p.err_num_no_match(len(p.parameterNames), 1)
		}
	}
	if isSlice(p.inputType) {
		typ := reflect.TypeOf(parameter)
		if typ.Kind() == reflect.Slice {
			return p.toAnySlice(parameter, typ.Elem().String())
		} else {
			fmt.Println("The variable is not a slice.")
		}
	}
	if p.inputType == "map" {
		if m, ok := parameter.(map[string]any); ok {
			args = make([]any, 0)
			for _, name := range p.parameterNames {
				if v, ok := m[name]; ok {
					args = append(args, v)
				} else {
					return nil, p.err_params_no_found(name)
				}
			}
			if len(args) == len(p.parameterNames) {
				return args, nil
			} else {
				return nil, p.err_num_no_match(len(p.parameterNames), len(args))
			}
		} else {
			return nil, p.err_invalid_parameter(reflect.TypeOf(parameter).String())
		}
	}

	val := reflect.ValueOf(parameter)
	typ := reflect.TypeOf(parameter)
	if val.Kind() == reflect.Struct {
		ptr := reflect.New(val.Type())
		ptr.Elem().Set(val)
		val = ptr
		typ = val.Type()
	}
	kind := typ.Kind()
	switch kind {
	case reflect.Struct, reflect.Ptr:
		args = make([]any, 0)
		for _, name := range p.parameterNames {
			method := val.MethodByName("Get" + upperFirstLetter(name))
			if method.IsValid() {
				//in := []reflect.Value{}
				//if method.Type().NumIn() == 1 {
				//	in = []reflect.Value{val}
				//}
				if vals := method.Call(nil); len(vals) > 0 {
					args = append(args, vals[0].Interface())
				}
			} else {
				return nil, p.err_params_no_found(name)
			}
		}
		if len(args) == len(p.parameterNames) {
			return args, nil
		} else {
			return nil, p.err_num_no_match(len(p.parameterNames), len(args))
		}
	}
	return nil, p.err_invalid_parameter(reflect.TypeOf(parameter).String())
}

func (p *paramBean) toAnySlice(parameter any, typename string) (r []any, err error) {
	length := len(p.parameterNames)
	switch typename {
	case "string":
		if slice, ok := parameter.([]string); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	case "int64":
		if slice, ok := parameter.([]int64); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	case "int":
		if slice, ok := parameter.([]int); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	case "int32":
		if slice, ok := parameter.([]int32); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	case "int16":
		if slice, ok := parameter.([]int16); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	case "int8":
		if slice, ok := parameter.([]int8); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	case "uint64":
		if slice, ok := parameter.([]uint64); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	case "uint32":
		if slice, ok := parameter.([]uint32); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	case "uint16":
		if slice, ok := parameter.([]uint16); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	case "uint8":
		if slice, ok := parameter.([]uint8); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	case "float64":
		if slice, ok := parameter.([]float64); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	case "float32":
		if slice, ok := parameter.([]float32); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	case "bool":
		if slice, ok := parameter.([]bool); ok {
			if length == len(slice) {
				r = make([]any, length)
				for i := range slice {
					r[i] = slice[i]
				}
			} else {
				err = p.err_num_no_match(length, len(slice))
			}
		}
	default:
		return nil, p.err_invalid_parameter(reflect.TypeOf(parameter).String())
	}
	return
}

func isDBType(putType string) bool {
	if putType[0:1] == "*" {
		putType = putType[1:]
	}
	switch putType {
	case "int64", "int32", "int", "int16", "int8":
		return true
	case "float64", "float32":
		return true
	case "uint64", "uint32", "uint", "uint16", "uint8":
		return true
	case "string", "byte", "rune", "[]byte", "bool", "time", "Time", "time.Time":
		return true
	}
	return false
}

func isSlice(putType string) bool {
	if putType[:2] == "[]" {
		return true
	}
	return false
}

func upperFirstLetter(s string) string {
	return strings.ToUpper(string(s[0])) + s[1:]
}
