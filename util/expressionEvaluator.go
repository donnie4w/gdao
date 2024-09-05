// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package util

import (
	"github.com/donnie4w/gdao/gdaoStruct"
)

func EvaluateExpression(expression string, context *gdaoStruct.ParamContext) (bool, error) {
	return Evaluate(resolveExpression(expression, context))
}

func resolveExpression(expression string, context *gdaoStruct.ParamContext) string {
	if expression == "" {
		return ""
	}
	return extractVariables(expression, context)
}
