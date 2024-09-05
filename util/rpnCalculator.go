// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package util

import (
	"fmt"
	"github.com/donnie4w/gdao/gdaoStruct"
	"regexp"
	"strconv"
	"strings"
)

func Evaluate(expression string) (r bool, err error) {
	defer Recover(&err)
	if expression == "" {
		return false, nil
	}
	if rpn, err := toRPN(expression); err == nil {
		return evaluateRPN(rpn), nil
	} else {
		return false, err
	}
}

func EvaluateWithAny(expression string, paramer any) (bool, error) {
	return Evaluate(resolveExpression(expression, gdaoStruct.NewParamContext(paramer)))
}

func EvaluateWithContext(expression string, context map[string]interface{}) (bool, error) {
	return Evaluate(resolveExpression(expression, gdaoStruct.NewParamContext(context)))
}

func toRPN(expression string) ([]string, error) {
	regex := regexp.MustCompile(`\s*(\(|\)|&&|\|\||==|!=|<=|>=|\+|\-|\*|\/|%|&|\||~|>|<|\w+|"[^"]*"|'[^']*')\s*`)
	matches := regex.FindAllString(expression, -1)

	operatorStack := []string{}
	output := []string{}

	for _, token := range matches {
		token = strings.TrimSpace(token)
		if isNumeric(token) || isStringLiteral(token) || token == "nil" {
			output = append(output, token)
		} else if isOperator(token) {
			for len(operatorStack) > 0 && precedence(operatorStack[len(operatorStack)-1]) >= precedence(token) {
				output = append(output, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			operatorStack = append(operatorStack, token)
		} else if token == "(" {
			operatorStack = append(operatorStack, token)
		} else if token == ")" {
			for len(operatorStack) > 0 && operatorStack[len(operatorStack)-1] != "(" {
				output = append(output, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			if len(operatorStack) > 0 {
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
		} else {
			return nil, fmt.Errorf("Unparsed characters found:%s", token)
		}
	}

	for len(operatorStack) > 0 {
		output = append(output, operatorStack[len(operatorStack)-1])
		operatorStack = operatorStack[:len(operatorStack)-1]
	}

	return output, nil
}

func evaluateRPN(tokens []string) bool {
	var stack []interface{}
	for _, token := range tokens {
		if isNumeric(token) {
			value, _ := strconv.ParseFloat(token, 64)
			stack = append(stack, value)
		} else if isStringLiteral(token) {
			stack = append(stack, token[1:len(token)-1])
		} else if token == "nil" {
			stack = append(stack, nil)
		} else if isOperator(token) {
			var b, a interface{}
			if token != "~" {
				b = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				a = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			} else {
				a = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, applyOperator(token, a, b))
		} else {
			stack = append(stack, token)
		}
	}

	return toBoolean(stack[len(stack)-1])
}

func applyOperator(operator string, a, b interface{}) interface{} {
	switch operator {
	case "+":
		return a.(float64) + b.(float64)
	case "-":
		return a.(float64) - b.(float64)
	case "*":
		return a.(float64) * b.(float64)
	case "/":
		return a.(float64) / b.(float64)
	case "%":
		return float64(int(a.(float64)) % int(b.(float64)))
	case "&":
		return float64(int(a.(float64)) & int(b.(float64)))
	case "|":
		return float64(int(a.(float64)) | int(b.(float64)))
	case "~":
		return float64(^int(a.(float64)))
	case "&&":
		return toBoolean(a) && toBoolean(b)
	case "||":
		return toBoolean(a) || toBoolean(b)
	case "==":
		return a == b
	case "!=":
		return a != b
	case ">":
		if a == nil || b == nil {
			return false
		}
		return a.(float64) > b.(float64)
	case "<":
		if a == nil || b == nil {
			return false
		}
		return a.(float64) < b.(float64)
	case ">=":
		if a == nil || b == nil {
			return false
		}
		return a.(float64) >= b.(float64)
	case "<=":
		if a == nil || b == nil {
			return false
		}
		return a.(float64) <= b.(float64)
	default:
		panic("Unsupported operator: " + operator)
	}
}

func evaluateRPN1(tokens []string) bool {
	stack := []interface{}{}

	for _, token := range tokens {
		fmt.Println("token:", token)
		if isNumeric(token) {
			value, _ := strconv.ParseFloat(token, 64)
			stack = append(stack, value)
		} else if isStringLiteral(token) {
			stack = append(stack, token[1:len(token)-1])
		} else if token == "nil" {
			stack = append(stack, nil)
		} else if isOperator(token) {
			b := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			var a interface{}
			if token != "~" {
				a = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}

			switch token {
			case "+":
				stack = append(stack, a.(float64)+b.(float64))
			case "-":
				stack = append(stack, a.(float64)-b.(float64))
			case "*":
				stack = append(stack, a.(float64)*b.(float64))
			case "/":
				stack = append(stack, a.(float64)/b.(float64))
			case "%":
				stack = append(stack, float64(int(a.(float64))%int(b.(float64))))
			case "&":
				stack = append(stack, int(a.(float64))&int(b.(float64)))
			case "|":
				stack = append(stack, int(a.(float64))|int(b.(float64)))
			case "~":
				stack = append(stack, ^int(b.(float64)))
			case "&&":
				stack = append(stack, toBoolean(a) && toBoolean(b))
			case "||":
				stack = append(stack, toBoolean(a) || toBoolean(b))
			case "==":
				stack = append(stack, equals(a, b))
			case "!=":
				stack = append(stack, !equals(a, b))
			case ">":
				if a == nil || b == nil {
					return false
				}
				stack = append(stack, compare(a, b) > 0)
			case "<":
				if a == nil || b == nil {
					return false
				}
				stack = append(stack, compare(a, b) < 0)
			case ">=":
				if a == nil || b == nil {
					return false
				}
				stack = append(stack, compare(a, b) >= 0)
			case "<=":
				if a == nil || b == nil {
					return false
				}
				stack = append(stack, compare(a, b) <= 0)
			default:
				panic(fmt.Sprintf("Unsupported operator: %s", token))
			}
		} else {
			stack = append(stack, token)
		}
	}

	return toBoolean(stack[len(stack)-1])
}

func isNumeric(str string) bool {
	_, err := strconv.ParseFloat(str, 64)
	return err == nil
}

func isStringLiteral(str string) bool {
	return len(str) > 1 && (strings.HasPrefix(str, "\"") && strings.HasSuffix(str, "\"") || strings.HasPrefix(str, "'") && strings.HasSuffix(str, "'"))
}

func isOperator(token string) bool {
	return strings.Contains("+-*/%&&||==!=><>=<=&|~", token)
}

func precedence(operator string) int {
	switch operator {
	case "||", "&&":
		return 1
	case "==", "!=", ">", "<", ">=", "<=":
		return 2
	case "+", "-":
		return 3
	case "*", "/", "%":
		return 4
	case "|", "&", "~":
		return 5
	default:
		return 0
	}
}

func toBoolean(obj interface{}) bool {
	switch v := obj.(type) {
	case bool:
		return v
	case float64:
		return v != 0
	case string:
		return v != ""
	default:
		return obj != nil
	}
}

func equals(a, b interface{}) bool {
	return a == b
}

func compare(a, b interface{}) int {
	aNum, aIsNum := toFloat64(a)
	bNum, bIsNum := toFloat64(b)

	if aIsNum && bIsNum {
		if aNum < bNum {
			return -1
		} else if aNum > bNum {
			return 1
		} else {
			return 0
		}
	}
	aStr, aIsStr := a.(string)
	bStr, bIsStr := b.(string)

	if aIsStr && bIsStr {
		return strings.Compare(aStr, bStr)
	}
	panic(fmt.Sprintf("Cannot compare %v with %v", a, b))
}

func toFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case int:
		return float64(v), true
	case uint:
		return float64(v), true
	case int8:
		return float64(v), true
	case uint8:
		return float64(v), true
	case int16:
		return float64(v), true
	case uint16:
		return float64(v), true
	case int32:
		return float64(v), true
	case uint32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}
