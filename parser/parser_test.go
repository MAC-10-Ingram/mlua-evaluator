package parser

import (
	"mlua-evaluator/ast"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	input := `
-- This is a file comment
@Author("Gemini")
script MyClass

-- This is a property comment
property number value = 10

-- This is a method comment
method string GetValue(number input, string name)
end

test "My first test" "should return 10" = {GetValue, 5, "testName", 10}
`

	expected := &ast.ParsedMluaFile{
		ClassName: "MyClass",
		Comments:  []string{"This is a file comment"},
		Metadata: map[string]string{
			"Author": "Gemini",
		},
		Properties: []ast.Property{
			{
				Name:     "value",
				Type:     "number",
				Comments: []string{"This is a property comment"},
			},
		},
		Methods: []ast.Method{
			{
				Name:       "GetValue",
				ReturnType: "string",
				Parameters: []ast.MethodParam{
					{Name: "input", Type: "number"},
					{Name: "name", Type: "string"},
				},
				Comments: []string{"This is a method comment"},
			},
		},
		TestCases: []ast.TestCase{
			{
				Name:        "My first test",
				Description: "should return 10",
				Target:      "GetValue",
				Input:       []string{"5", `"testName"`},
				Expected:    "10",
				Comments:    []string{},
			},
		},
	}

	parsed, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse returned an error: %v", err)
	}

	// Compare ClassName
	if parsed.ClassName != expected.ClassName {
		t.Errorf("Expected ClassName %q, got %q", expected.ClassName, parsed.ClassName)
	}

	// Compare Comments
	if !reflect.DeepEqual(parsed.Comments, expected.Comments) {
		t.Errorf("Expected Comments %v, got %v", expected.Comments, parsed.Comments)
	}

	// Compare Metadata
	if !reflect.DeepEqual(parsed.Metadata, expected.Metadata) {
		t.Errorf("Expected Metadata %v, got %v", expected.Metadata, parsed.Metadata)
	}

	// Compare Properties
	if len(parsed.Properties) != len(expected.Properties) {
		t.Fatalf("Expected %d properties, got %d", len(expected.Properties), len(parsed.Properties))
	}
	for i, p := range parsed.Properties {
		if p.Name != expected.Properties[i].Name {
			t.Errorf("Property %d: Expected Name %q, got %q", i, expected.Properties[i].Name, p.Name)
		}
		if p.Type != expected.Properties[i].Type {
			t.Errorf("Property %d: Expected Type %q, got %q", i, expected.Properties[i].Type, p.Type)
		}
		if !reflect.DeepEqual(p.Comments, expected.Properties[i].Comments) {
			t.Errorf("Property %d: Expected Comments %v, got %v", i, expected.Properties[i].Comments, p.Comments)
		}
	}

	// Compare Methods
	if len(parsed.Methods) != len(expected.Methods) {
		t.Fatalf("Expected %d methods, got %d", len(expected.Methods), len(parsed.Methods))
	}
	for i, m := range parsed.Methods {
		if m.Name != expected.Methods[i].Name {
			t.Errorf("Method %d: Expected Name %q, got %q", i, expected.Methods[i].Name, m.Name)
		}
		if m.ReturnType != expected.Methods[i].ReturnType {
			t.Errorf("Method %d: Expected ReturnType %q, got %q", i, expected.Methods[i].ReturnType, m.ReturnType)
		}
		if !reflect.DeepEqual(m.Parameters, expected.Methods[i].Parameters) {
			t.Errorf("Method %d: Expected Parameters %v, got %v", i, expected.Methods[i].Parameters, m.Parameters)
		}
		if !reflect.DeepEqual(m.Comments, expected.Methods[i].Comments) {
			t.Errorf("Method %d: Expected Comments %v, got %v", i, expected.Methods[i].Comments, m.Comments)
		}
	}

	// Compare TestCases
	if len(parsed.TestCases) != len(expected.TestCases) {
		t.Fatalf("Expected %d test cases, got %d", len(expected.TestCases), len(parsed.TestCases))
	}
	for i, tc := range parsed.TestCases {
		if tc.Name != expected.TestCases[i].Name {
			t.Errorf("TestCase %d: Expected Name %q, got %q", i, expected.TestCases[i].Name, tc.Name)
		}
		if tc.Description != expected.TestCases[i].Description {
			t.Errorf("TestCase %d: Expected Description %q, got %q", i, expected.TestCases[i].Description, tc.Description)
		}
		if tc.Target != expected.TestCases[i].Target {
			t.Errorf("TestCase %d: Expected Target %q, got %q", i, expected.TestCases[i].Target, tc.Target)
		}
		if !reflect.DeepEqual(tc.Input, expected.TestCases[i].Input) {
			t.Errorf("TestCase %d: Expected Input %v, got %v", i, expected.TestCases[i].Input, tc.Input)
		}
		if tc.Expected != expected.TestCases[i].Expected {
			t.Errorf("TestCase %d: Expected Expected %q, got %q", i, expected.TestCases[i].Expected, tc.Expected)
		}
	}
}
