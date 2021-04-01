package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/alekseysychev/PixelatorGo/pkg/pixelator"
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
	start := time.Now()

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

	err = pixelator.Compile(inputFile, outputFile, pixelator.Settings{
		ClusterSize: clusterSize,
		Quality:     quality,
	})

	if err != nil {
		log.Fatalln(err)
	}
	stop := time.Now()

	log.Printf("Options:\n")
	log.Printf("  -  cluster size : %d\n", clusterSize)
	log.Printf("  -  quality      : %d\n", quality)
	log.Printf("Script working time: %s\n", stop.Sub(start))
}
