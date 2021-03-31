package main

import (
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	pixelator "github.com/alekseysychev/PixelatorGo/pkg/pixelator"
)

func main() {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			var path string = req.URL.Path
			if path == "/" {
				path = "/index.html"
			}
			http.ServeFile(w, req, filepath.Join("./web", path))
		case http.MethodPost:

			req.ParseMultipartForm(32 << 20)
			inputFile, _, err := req.FormFile("file")
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				log.Println("err file")
				return
			}
			defer inputFile.Close()

			clusterSize, err := strconv.Atoi(req.FormValue("clusterSize"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				log.Println("err s")
				return
			}

			quality, err := strconv.Atoi(req.FormValue("quality"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				log.Println("err q")
				return
			}

			err = pixelator.Init(inputFile, w).Compile(clusterSize, quality)
			if err != nil {
				log.Println("err pixelator")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	http.HandleFunc("/", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
