package main

import (
	"flag"
	"fmt"
	"log"
	"nuconv/internal/pkgs"
)

func run(patterns []string) (err error) {
	parser := pkgs.NewParser()

	err = parser.Load(patterns...)
	if err != nil {
		err = fmt.Errorf("load packages: %w", err)
		return
	}

	err = parser.Generate()
	if err != nil {
		err = fmt.Errorf("generate packages: %w", err)
		return
	}

	return
}

func main() {
	flag.Parse()

	patterns := flag.Args()

	err := run(patterns)
	if err != nil {
		log.Fatal(err)
	}
}
