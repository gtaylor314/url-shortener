package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gtaylor314/url-shortener/urlshort"
)

func main() {
	// declare flag for yaml file
	yamlPtr := flag.String("yaml", "", "used to provide the filename of the yaml file you wish to use")
	jsonPtr := flag.String("json", "", "used to provide the filename of the json file you wish to use")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	// assign yaml the default configuration and update it if the yaml flag is set
	yamlString := `
    - path: /urlshort
      url: https://github.com/gophercises/urlshort
    - path: /urlshort-final
      url: https://github.com/gophercises/urlshort/tree/solution
    `
	yaml := []byte(yamlString)

	// updating yaml if the yaml flag is set
	if *yamlPtr != "" {
		yaml = readFile(yamlPtr)
	}

	yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
	if err != nil {
		panic(err)
	}

	// Build the JSONHandler using the yamlHandler as the
	// fallback
	// assign json the default configuration and update it if the json flag is set
	jsonString := `[{
		"path": "/watch",
		"url": "https://www.youtube.com"
	}, {
		"path": "/play",
		"url": "https://www.playstation.com"
	}]`
	json := []byte(jsonString)
	// update if json flag is set
	if *jsonPtr != "" {
		json = readFile(jsonPtr)
	}

	jsonHandler, err := urlshort.JSONHandler(json, yamlHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

// readFile() can be used for either yaml or json
func readFile(s *string) []byte {
	// read from file and store bytes in bSlice
	bSlice, err := os.ReadFile(*s)
	if err != nil {
		fmt.Println("error reading file " + *s + " -- " + err.Error())
		panic(err)
	}
	return bSlice
}
