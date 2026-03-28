package runner

import (
	"fmt"
	"mlua-evaluator/ast"
	"strings"

	"github.com/yuin/gopher-lua"
)

// RunTests executes the tests defined in the parsed mLua file.
func RunTests(parsedFile *ast.ParsedMluaFile) error {
	if parsedFile == nil {
		return fmt.Errorf("cannot run tests on a nil parsed file")
	}

	fmt.Printf("Testing script: %s\n", parsedFile.ClassName)
	fmt.Println("========================================")

	for _, tc := range parsedFile.TestCases {
		fmt.Printf("Running test: %s (%s)\n", tc.Name, tc.Description)

		L := lua.NewState()
		defer L.Close()

		script, err := generateTestScript(parsedFile, tc)
		if err != nil {
			fmt.Printf("❌ FAILED to generate script: %v\n", err)
			continue // Move to the next test case
		}

		if err := L.DoString(script); err != nil {
			fmt.Printf("❌ FAILED: %v\n", err)
			// In a real test runner, we would collect failures and report at the end.
			// For now, we continue to the next test.
			continue
		}

		fmt.Printf("✅ PASSED\n")
	}

	fmt.Println("========================================")
	return nil
}

// generateTestScript creates a self-contained Lua script for a single test case.
func generateTestScript(fileAst *ast.ParsedMluaFile, testCase ast.TestCase) (string, error) {
	var sb strings.Builder

	// 1. Define the class and its methods based on the AST
	sb.WriteString(fmt.Sprintf("%s = {}\n", fileAst.ClassName))
	for _, prop := range fileAst.Properties {
		// A default 'nil' value is assumed for now.
		sb.WriteString(fmt.Sprintf("%s.%s = nil\n", fileAst.ClassName, prop.Name))
	}

	for _, method := range fileAst.Methods {
		var paramNames []string
		for _, p := range method.Parameters {
			paramNames = append(paramNames, p.Name)
		}
		// The body of the method is empty, as we are mocking/testing behavior.
		sb.WriteString(fmt.Sprintf("function %s:%s(%s)\n", fileAst.ClassName, method.Name, strings.Join(paramNames, ", ")))
		sb.WriteString("  -- Method body is not executed in this test setup\n")
		// If methods returned values, we would need a mock setup here.
		sb.WriteString("end\n\n")
	}

	// 2. Find the target method from the AST
	targetMethod, err := findMethod(fileAst, testCase.Target)
	if err != nil {
		return "", err
	}

	// 3. Set up input arguments
	if len(testCase.Input) != len(targetMethod.Parameters) {
		return "", fmt.Errorf("argument count mismatch for test '%s': expected %d, got %d", testCase.Name, len(targetMethod.Parameters), len(testCase.Input))
	}

	var callArgs []string
	for i, param := range targetMethod.Parameters {
		argVar := fmt.Sprintf("arg_%s", param.Name)
		sb.WriteString(fmt.Sprintf("local %s = %s\n", argVar, testCase.Input[i]))
		callArgs = append(callArgs, argVar)
	}

	// 4. Call the function
	// For now, we assume methods are called on the class itself.
	actualValueVar := "actual_value"
	sb.WriteString(fmt.Sprintf("local %s = %s:%s(%s)\n", actualValueVar, fileAst.ClassName, testCase.Target, strings.Join(callArgs, ", ")))

	// 5. Assert the result
	expectedValue := testCase.Expected
	errorMsg := fmt.Sprintf("Assertion failed in test '%s': expected ' .. tostring(%s) .. ' but got ' .. tostring(%s)", testCase.Name, expectedValue, actualValueVar)
	sb.WriteString(fmt.Sprintf(`if tostring(%s) ~= tostring(%s) then error("%s") end`, actualValueVar, expectedValue, errorMsg))

	return sb.String(), nil
}

func findMethod(fileAst *ast.ParsedMluaFile, methodName string) (*ast.Method, error) {
	for _, m := range fileAst.Methods {
		if m.Name == methodName {
			return &m, nil
		}
	}
	return nil, fmt.Errorf("method '%s' not found in script '%s'", methodName, fileAst.ClassName)
}
