package main

import (
	"fmt"
	"os"

	"mlua-evaluator/parser"
	"mlua-evaluator/runner"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: mlua-evaluator <mlua_file> <dataset_json>")
		os.Exit(1)
	}

	mluaPath := os.Args[1]
	datasetPath := os.Args[2]

	mluaData, err := os.ReadFile(mluaPath)
	if err != nil {
		fmt.Printf("Error reading mLua file: %v\n", err)
		os.Exit(1)
	}

	// Transpile
	transpiledCode, err := parser.Transpile(string(mluaData))
	if err != nil {
		fmt.Printf("Error transpiling mLua file: %v\n", err)
		os.Exit(1)
	}

	// Run tests
	if err := runner.RunTests(transpiledCode, datasetPath); err != nil {
		fmt.Println("Tests failed.")
		os.Exit(1)
	}
	fmt.Println("All tests passed successfully.")
}
