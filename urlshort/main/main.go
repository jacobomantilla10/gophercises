package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jacobomantilla10/gophercises/urlshort"
)

func main() {
	var yamlFile = flag.String("yamlfile", "default.yaml", "Declare the name of the YAML file containing the paths")
	var jsonFile = flag.String("jsonfile", "default.json", "Declare the name of the JSON file containing the paths")
	flag.Parse()
	fmt.Println(*yamlFile)
	fmt.Println(*jsonFile)

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml, err := os.ReadFile(filepath.Join("yaml", *yamlFile))
	if err != nil {
		if os.IsNotExist(err) {
			yaml, _ = os.ReadFile(filepath.Join("yaml", "default.yaml"))
		} else {
			panic(err)
		}
	}
	yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
	if err != nil {
		panic(err)
	}

	//Build the JSONHandler using the YAMLHandler as the fallback
	json, err := os.ReadFile(filepath.Join("json", *jsonFile))
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Testing")
			json, _ = os.ReadFile(filepath.Join("json", "default.json"))
		} else {
			panic(err)
		}
	}
	jsonHandler, err := urlshort.JSONHandler([]byte(json), yamlHandler)
	if err != nil {
		panic(err)
	}

	//Build the BoltdbHandler using the JSONHandler as the fallback
	boltHandler, err := urlshort.BoltHandler(jsonHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", boltHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
