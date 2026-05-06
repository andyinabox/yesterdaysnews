package main

import (
	"flag"
	"fmt"
	"os"

	"code.andydayton.com/andy/yesterdaysnews/domain/logger"
	"code.andydayton.com/andy/yesterdaysnews/pkg/markov/basicchain"
)

var prefixLength int
var verbose bool
var inputFile string

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.IntVar(&prefixLength, "p", 2, "prefix length")
	flag.StringVar(&inputFile, "i", "", "input file")
	flag.Parse()

	logger.SetDefault(&logger.Config{
		Verbose: true,
	})

	if inputFile == "" {
		panic("input file is required")
	}
}

func main() {

	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	chain := basicchain.New(prefixLength)
	chain.Build(file)

	model, err := chain.Save()
	if err != nil {
		panic(err)
	}

	fmt.Print(string(model))

}
