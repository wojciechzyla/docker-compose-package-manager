package main

import (
	"fmt"
	"os"
)

func run() error {
	const templateFile = "template/template.yaml"
	const dataFile = "template/values.yaml"
	const outputFile = "parsed/parsed.yaml"
	if err := Parse(templateFile, dataFile, outputFile); err != nil {
		return err
	}
	fmt.Printf("file %s was generated.\n", outputFile)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
