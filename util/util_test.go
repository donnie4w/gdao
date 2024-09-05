// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package util

import (
	"fmt"
	"testing"
)

func Test_evaluate(t *testing.T) {
	context := map[string]interface{}{
		"x": 10,
		"y": "hello",
		"z": nil,
		"s": "11",
	}

	expressions := []string{
		"nil==s",
		"( x|1 < 3|15) && (y== 'hello') && (z == nil)",
		"( x|1 < 3|15) && (y== \"hello\") && (z != nil)",
		"( x|1 < 3|15) || (y== \"hello1\") && (z == nil)",
		"( s==\"s\") && (y== \"hello\") && (z == nil)",
		"'www>>>>2'!= null",
	}

	for _, expression := range expressions {
		result, err := EvaluateWithContext(expression, context)
		fmt.Println("Result:", result, ",Error:", err)
	}
}
