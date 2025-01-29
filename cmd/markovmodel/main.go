package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdaysnews/pkg/markov"
)

var prefixLength int
var verbose bool
var inputFile string

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.IntVar(&prefixLength, "p", 2, "prefix length")
	flag.StringVar(&inputFile, "i", "", "input file")
	flag.Parse()

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	if inputFile == "" {
		log.Fatal("input file is required")
	}
}

func main() {

	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	chain := markov.NewBasicChain(prefixLength)
	chain.Build(file)

	model, err := chain.Save()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(string(model))

}
