package main

import (
	"flag"
	"log"
	"os"

	pixelator "github.com/alekseysychev/PixelatorGo/pkg/pixelator"
)

var (
	clusterSize int
	quality     int
)

func init() {
	flag.IntVar(&clusterSize, "n", 1, "cluster size. default : 1")
	flag.IntVar(&quality, "q", 100, "quality. default : 100")
	flag.Parse()
}

func main() {
	inputFile, err := os.Open("./examples/input.jpg")
	if err != nil {
		log.Fatalln(err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create("./examples/output.jpg")
	if err != nil {
		log.Fatalln(err)
	}
	defer outputFile.Close()

	err = pixelator.Init(inputFile, outputFile).Compile(clusterSize, quality)
	if err != nil {
		log.Fatalln(err)
	}
}
