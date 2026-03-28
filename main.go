package main

import (
	"fmt"
	"os"

	"mlua-evaluator/parser"
	"mlua-evaluator/runner"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mlua-evaluator <mlua_file>")
		os.Exit(1)
	}

	mluaPath := os.Args[1]

	mluaData, err := os.ReadFile(mluaPath)
	if err != nil {
		fmt.Printf("Error reading mLua file: %v\n", err)
		os.Exit(1)
	}

	parsedFile, err := parser.Parse(string(mluaData))
	if err != nil {
		fmt.Printf("Error parsing mLua file: %v\n", err)
		os.Exit(1)
	}

	// Run tests
	if err := runner.RunTests(parsedFile); err != nil {
		fmt.Println("Tests encountered failures.")
		os.Exit(1)
	}
	fmt.Println("All tests processed.")
}
