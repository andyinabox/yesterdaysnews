package main

// import (
// 	"encoding/json"
// 	"os"
// 	"path/filepath"

// 	"github.com/charmbracelet/log"
// 	"gitlab.com/andyinabox/yesterdaysnews/domain/manifest"
// )

// func main() {

// 	outputDir := "dist"

// 	manifest, err := manifest.Create(outputDir)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	outFile := filepath.Join(outputDir, "manifest.json")

// 	b, err := json.Marshal(manifest)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Info("outputting manifest file", "file", outFile)
// 	err = os.WriteFile(outFile, b, os.ModePerm)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
