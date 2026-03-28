package parser

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

var (
	annotationRegex = regexp.MustCompile(`^\s*@[A-Za-z0-9_]+(\(.*?\))?\s*$`)
	scriptRegex     = regexp.MustCompile(`^(\s*)script\s+([A-Za-z0-9_]+)`)
	propertyRegex   = regexp.MustCompile(`^(\s*)property\s+([A-Za-z0-9_]+)\s+([A-Za-z0-9_]+)\s*(?:=\s*(.*))?$`)
	methodRegex     = regexp.MustCompile(`^(\s*)method\s+([A-Za-z0-9_]+)\s+([A-Za-z0-9_]+)\s*\((.*)\)\s*$`)
)

type MethodParam struct {
	Name string
	Type string
}

// parseArgs takes a string like "number amount, string name" and returns a slice of MethodParam
func parseArgs(argString string) []MethodParam {
	if strings.TrimSpace(argString) == "" {
		return nil
	}
	params := []MethodParam{}
	parts := strings.Split(argString, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		subParts := strings.Fields(part)
		if len(subParts) >= 2 {
			paramType := strings.Join(subParts[:len(subParts)-1], " ")
			paramName := subParts[len(subParts)-1]
			params = append(params, MethodParam{Name: paramName, Type: paramType})
		} else if len(subParts) == 1 {
			params = append(params, MethodParam{Name: subParts[0], Type: "any"})
		}
	}
	return params
}

// Transpile converts mLua syntax into standard Lua 5.3 syntax
func Transpile(mluaCode string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(mluaCode))
	var out []string
	currentClass := ""

	for scanner.Scan() {
		line := scanner.Text()

		if annotationRegex.MatchString(line) {
			continue
		}

		if matches := scriptRegex.FindStringSubmatch(line); len(matches) > 0 {
			indent := matches[1]
			currentClass = matches[2]
			out = append(out, fmt.Sprintf("%s%s = {}", indent, currentClass))
			continue
		}

		if matches := propertyRegex.FindStringSubmatch(line); len(matches) > 0 {
			indent := matches[1]
			name := matches[3]
			val := matches[4]
			if val == "" {
				val = "nil"
			}
			if currentClass != "" {
				out = append(out, fmt.Sprintf("%s%s.%s = %s", indent, currentClass, name, val))
			} else {
				out = append(out, fmt.Sprintf("%slocal %s = %s", indent, name, val))
			}
			continue
		}

		if matches := methodRegex.FindStringSubmatch(line); len(matches) > 0 {
			indent := matches[1]
			methodName := matches[3]
			argsString := matches[4]

			params := parseArgs(argsString)
			var paramNames []string
			var sigParts []string
			for _, p := range params {
				paramNames = append(paramNames, p.Name)
				sigParts = append(sigParts, fmt.Sprintf("%s:%s", p.Name, p.Type))
			}

			signature := fmt.Sprintf("%s-- @method_signature: %s(%s)", indent, methodName, strings.Join(sigParts, ","))
			out = append(out, signature)

			if currentClass != "" {
				out = append(out, fmt.Sprintf("%sfunction %s:%s(%s)", indent, currentClass, methodName, strings.Join(paramNames, ", ")))
			} else {
				out = append(out, fmt.Sprintf("%sfunction %s(%s)", indent, methodName, strings.Join(paramNames, ", ")))
			}
			continue
		}

		out = append(out, line)
	}

	return strings.Join(out, "\n"), scanner.Err()
}
