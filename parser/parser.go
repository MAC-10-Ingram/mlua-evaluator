package parser

import (
	"bufio"
	"fmt"
	"mlua-evaluator/ast"
	"regexp"
	"strings"
)

var (
	classRegex    = regexp.MustCompile(`^\s*script\s+([A-Za-z0-9_]+)`)
	propertyRegex = regexp.MustCompile(`^\s*property\s+([A-Za-z0-9_]+)\s+([A-Za-z0-9_]+)\s*(?:=\s*(.*))?$`)
	methodRegex   = regexp.MustCompile(`^\s*method\s+([A-Za-z0-9_]+)\s+([A-Za-z0-9_]+)\s*\((.*)\)\s*.*$`)
	testCaseRegex = regexp.MustCompile(`^\s*test\s*"([^"]+)"\s*"([^"]+)"\s*=\s*{([^}]+)}`)
	metaRegex     = regexp.MustCompile(`^\s*@([A-Za-z0-9_]+)\((.*)\)`)
	commentRegex  = regexp.MustCompile(`^\s*--\s*(.*)`)
)

// Parse takes the content of a .mlua file and returns a structured *ast.ParsedMluaFile.
func Parse(content string) (*ast.ParsedMluaFile, error) {
	parsedFile := &ast.ParsedMluaFile{
		Metadata: make(map[string]string),
	}
	scanner := bufio.NewScanner(strings.NewReader(content))
	var currentComments []string

	for scanner.Scan() {
		line := scanner.Text()

		if matches := commentRegex.FindStringSubmatch(line); len(matches) > 1 {
			currentComments = append(currentComments, matches[1])
			continue
		}

		if matches := metaRegex.FindStringSubmatch(line); len(matches) > 2 {
			key := matches[1]
			value := strings.Trim(matches[2], `"`)
			parsedFile.Metadata[key] = value
			parsedFile.Comments = append(parsedFile.Comments, currentComments...)
			currentComments = nil
			continue
		}

		if matches := classRegex.FindStringSubmatch(line); len(matches) > 0 {
			parsedFile.ClassName = matches[1]
			parsedFile.Comments = append(parsedFile.Comments, currentComments...)
			currentComments = nil
			continue
		}

		if matches := propertyRegex.FindStringSubmatch(line); len(matches) > 0 {
			prop := ast.Property{
				Type:     matches[1],
				Name:     matches[2],
				Comments: currentComments,
			}
			parsedFile.Properties = append(parsedFile.Properties, prop)
			currentComments = nil
			continue
		}

		if matches := methodRegex.FindStringSubmatch(line); len(matches) > 0 {
			method := ast.Method{
				ReturnType: matches[1],
				Name:       matches[2],
				Parameters: parseMethodParams(matches[3]),
				Comments:   currentComments,
			}
			parsedFile.Methods = append(parsedFile.Methods, method)
			currentComments = nil
			continue
		}

		if matches := testCaseRegex.FindStringSubmatch(line); len(matches) > 3 {
			testCase := ast.TestCase{
				Name:        matches[1],
				Description: matches[2],
				Comments:    currentComments,
			}

			parts := strings.Split(matches[3], ",")
			if len(parts) >= 2 {
				testCase.Target = strings.TrimSpace(parts[0])
				testCase.Input = parseTestInput(strings.Join(parts[1:len(parts)-1], ","))
				testCase.Expected = strings.TrimSpace(parts[len(parts)-1])
			}

			parsedFile.TestCases = append(parsedFile.TestCases, testCase)
			currentComments = nil
			continue
		}

		// If a line is not a comment and not a recognized structure, and there are pending comments,
		// associate them with the file-level comments if no class has been defined yet.
		if strings.TrimSpace(line) != "" && len(currentComments) > 0 && parsedFile.ClassName == "" {
			parsedFile.Comments = append(parsedFile.Comments, currentComments...)
			currentComments = nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading mLua content: %w", err)
	}

	return parsedFile, nil
}

func parseMethodParams(argString string) []ast.MethodParam {
	if strings.TrimSpace(argString) == "" {
		return nil
	}
	params := []ast.MethodParam{}
	parts := strings.Split(argString, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		subParts := strings.Fields(part)
		if len(subParts) >= 2 {
			paramType := subParts[0]
			paramName := subParts[1]
			params = append(params, ast.MethodParam{Name: paramName, Type: paramType})
		}
	}
	return params
}

func parseTestInput(inputStr string) []string {
	inputStr = strings.TrimSpace(inputStr)
	if inputStr == "" {
		return nil
	}
	parts := strings.Split(inputStr, ",")
	trimmedParts := make([]string, len(parts))
	for i, part := range parts {
		trimmedParts[i] = strings.TrimSpace(part)
	}
	return trimmedParts
}
