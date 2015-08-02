package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/variadico/decompose"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) != 2 {
		log.Println("usage: decompose docker-compose.yml")
		os.Exit(1)
	}

	inputFile := os.Args[1]

	// prefix is the directory that inputFile is in
	finAbs, err := filepath.Abs(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	decompose.Prefix = filepath.Base(filepath.Dir(finAbs)) + "_"

	services, err := decompose.ParseComposeFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range services {
		fmt.Println(s)
	}
}
