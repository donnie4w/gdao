// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package gdaoMapper

import (
	. "github.com/donnie4w/gdao/gdaoStruct"
	"strings"
	"unicode"
)

type trimNode struct {
	Prefix          string
	Suffix          string
	PrefixOverrides []string
	SuffixOverrides []string
	nodes           []sqlNode
}

func newTrimNode(prefix, suffix, prefixOverrides, suffixOverrides string) *trimNode {
	var prefixOverrideArr, suffixOverrideArr []string
	if prefixOverrides != "" {
		prefixOverrideArr = strings.Split(prefixOverrides, "|")
	}
	if suffixOverrides != "" {
		suffixOverrideArr = strings.Split(suffixOverrides, "|")
	}
	return &trimNode{
		Prefix:          prefix,
		Suffix:          suffix,
		PrefixOverrides: prefixOverrideArr,
		SuffixOverrides: suffixOverrideArr,
	}
}

func (tn *trimNode) addSqlNode(node sqlNode) {
	tn.nodes = append(tn.nodes, node)
}

func (tn *trimNode) apply(context *ParamContext) *ackContext {
	ac, err := newAckContext2("", context)
	if err != nil {
		panic(err)
	}

	for _, node := range tn.nodes {
		if node != nil {
			ac.Append(node.apply(context))
		}
	}

	sqlstr := ac.sqlbuilder.String()
	if sqlstr == "" {
		return ac
	}

	for len(sqlstr) > 0 && unicode.IsSpace(rune(sqlstr[0])) {
		sqlstr = sqlstr[1:]
	}

	if tn.PrefixOverrides != nil {
		for _, toRemove := range tn.PrefixOverrides {
			if strings.HasPrefix(sqlstr, toRemove) {
				sqlstr = sqlstr[len(toRemove):]
				break
			}
		}
	}

	if tn.SuffixOverrides != nil {
		for _, toRemove := range tn.SuffixOverrides {
			if strings.HasSuffix(sqlstr, toRemove) {
				sqlstr = sqlstr[:len(sqlstr)-len(toRemove)]
				break
			}
		}
	}

	if tn.Prefix != "" {
		sqlstr = tn.Prefix + sqlstr
	}
	if tn.Suffix != "" {
		sqlstr = sqlstr + tn.Suffix
	}
	ac.sqlbuilder.Reset()
	ac.sqlbuilder.WriteString(sqlstr)
	return ac
}
