package ast

// MethodParam represents a parameter of a method.
type MethodParam struct {
	Name string
	Type string
}

// Method represents a method of a class.
type Method struct {
	Name       string
	Parameters []MethodParam
	ReturnType string
	Comments   []string // Comments associated with the method
}

// Property represents a property of a class.
type Property struct {
	Name     string
	Type     string
	Comments []string // Comments associated with the property
}

// ParsedMluaFile represents the structured content of a .mlua file.
type ParsedMluaFile struct {
	ClassName  string
	Comments   []string // File-level comments
	Metadata   map[string]string
	Properties []Property
	Methods    []Method
	TestCases  []TestCase
}

// TestCase represents a test case defined in the mLua file.
type TestCase struct {
	Name        string
	Description string
	Target      string   // Target method or function
	Input       []string // Input values for the test
	Expected    string   // Expected output
	Comments    []string
}
