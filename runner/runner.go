package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// --- Structs for dataset.json ---
type TestCase struct {
	Name         string                 `json:"name"`
	TargetMethod string                 `json:"target_method"`
	Mocks        []string               `json:"mocks"`
	Inputs       map[string]interface{} `json:"inputs"`
	Asserts      []Assert               `json:"asserts"`
}

type Assert struct {
	Actual   string      `json:"actual"`
	Expected interface{} `json:"expected"`
}

type TestDataset struct {
	TestCases []TestCase `json:"test_cases"`
}

// --- Structs for signature parsing ---
type MethodParam struct {
	Name string
	Type string
}

var signatureRegex = regexp.MustCompile(`-- @method_signature:\s*([A-Za-z0-9_]+)\((.*)\)`)

// --- Helper Functions ---

// findAndParseSignature finds and parses the method signature comment.
func findAndParseSignature(code, methodName string) ([]MethodParam, error) {
	lines := strings.Split(code, "\n")
	for _, line := range lines {
		if matches := signatureRegex.FindStringSubmatch(line); len(matches) > 0 {
			parsedMethodName := matches[1]
			if parsedMethodName == methodName {
				argString := matches[2]
				if strings.TrimSpace(argString) == "" {
					return nil, nil
				}
				var params []MethodParam
				parts := strings.Split(argString, ",")
				for _, part := range parts {
					argParts := strings.Split(strings.TrimSpace(part), ":")
					if len(argParts) == 2 {
						params = append(params, MethodParam{Name: argParts[0], Type: argParts[1]})
					}
				}
				return params, nil
			}
		}
	}
	return nil, fmt.Errorf("signature for method '%s' not found", methodName)
}

// getLuaType maps a mLua type to a Lua primitive type string for the `type()` function.
func getLuaType(mluaType string) string {
	switch strings.ToLower(mluaType) {
	case "number", "integer", "float":
		return "number"
	case "string":
		return "string"
	case "boolean":
		return "boolean"
	case "table", "entity", "vector2", "vector3":
		return "table"
	default:
		// For custom classes or complex types, we assume they'll be tables.
		return "table"
	}
}

func convertJSONToLua(val interface{}) string {
	switch v := val.(type) {
	case string:
		return fmt.Sprintf("%q", v)
	case float64:
		return fmt.Sprintf("%g", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case nil:
		return "nil"
	case []interface{}:
		items := make([]string, len(v))
		for i, item := range v {
			items[i] = convertJSONToLua(item)
		}
		return "{" + strings.Join(items, ", ") + "}"
	case map[string]interface{}:
		items := make([]string, 0, len(v))
		for key, item := range v {
			items = append(items, fmt.Sprintf("[%q] = %s", key, convertJSONToLua(item)))
		}
		return "{" + strings.Join(items, ", ") + "}"
	default:
		// Fallback for other types like int, which get read as float64 by json.Unmarshal
		return fmt.Sprintf("%v", v)
	}
}

// --- Main Runner Logic ---

// RunTests executes the transpiled mLua code against the dataset.
func RunTests(mluaCode, datasetPath string) error {
	data, err := os.ReadFile(datasetPath)
	if err != nil {
		return fmt.Errorf("failed to read dataset: %v", err)
	}

	var dataset TestDataset
	if err := json.Unmarshal(data, &dataset); err != nil {
		return fmt.Errorf("failed to parse dataset json: %v", err)
	}

	for _, tc := range dataset.TestCases {
		fmt.Printf("Running test: %s\n", tc.Name)

		L := lua.NewState()
		defer L.Close()

		// Extract the method name from "Class:Method"
		targetParts := strings.Split(tc.TargetMethod, ":")
		if len(targetParts) != 2 {
			return fmt.Errorf("invalid target_method format: %s. Expected 'Class:Method'", tc.TargetMethod)
		}
		methodName := targetParts[1]

		// Find method signature
		params, err := findAndParseSignature(mluaCode, methodName)
		if err != nil {
			return err
		}

		var scriptBuilder strings.Builder
		scriptBuilder.WriteString(mluaCode + "\n")

		for _, mock := range tc.Mocks {
			scriptBuilder.WriteString(mock + "\n")
		}

		// --- Argument and Type Validation ---
		var orderedArgNames []string
		for _, p := range params {
			luaVal, ok := tc.Inputs[p.Name]
			if !ok {
				return fmt.Errorf("missing input for parameter '%s' in test '%s'", p.Name, tc.Name)
			}
			argVarName := "arg_" + p.Name
			orderedArgNames = append(orderedArgNames, argVarName)

			// 1. Declare local variable for the argument
			scriptBuilder.WriteString(fmt.Sprintf("local %s = %s\n", argVarName, convertJSONToLua(luaVal)))

			// 2. Add type check
			if p.Type != "any" {
				expectedLuaType := getLuaType(p.Type)
				errorMsg := fmt.Sprintf("Type mismatch for param '%s': expected %s, got ' .. type(%s)", p.Name, p.Type, argVarName)
				scriptBuilder.WriteString(fmt.Sprintf(`if type(%s) ~= "%s" then error("%s") end`, argVarName, expectedLuaType, errorMsg) + "\n")
			}
		}

		// --- Call Target Method ---
		scriptBuilder.WriteString(fmt.Sprintf("%s(%s)\n", tc.TargetMethod, strings.Join(orderedArgNames, ", ")))

		// --- Assertions ---
		for i, assert := range tc.Asserts {
			expectedLua := convertJSONToLua(assert.Expected)
			errorMsg := fmt.Sprintf("Assertion %d failed: For '%s', expected ' .. tostring(%s) .. ' but got ' .. tostring(%s)", i+1, assert.Actual, expectedLua, assert.Actual)
			scriptBuilder.WriteString(fmt.Sprintf(`if %s ~= %s then error("%s") end`, assert.Actual, expectedLua, errorMsg) + "\n")
		}

		if err := L.DoString(scriptBuilder.String()); err != nil {
			fmt.Printf("❌ FAILED: %v\n", err)
			return err
		}
		fmt.Printf("✅ PASSED\n")
	}

	return nil
}
