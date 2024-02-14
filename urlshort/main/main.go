package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/jacobomantilla10/gophercises/urlshort"
)

func main() {
	var yamlFile = flag.String("yamlfile", "default.yaml", "Declare the name of the YAML file containing the paths")
	var jsonFile = flag.String("jsonfile", "default.json", "Declare the name of the JSON file containing the paths")
	flag.Parse()
	fmt.Println(*yamlFile)
	fmt.Println(*jsonFile)

	mux := defaultMux()

	_, err := SetupDB()
	if err != nil {
		panic(err)
	}

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

// Creates the db buckets and returns the db
func SetupDB() (*bolt.DB, error) {
	db, err := bolt.Open("paths.db", 0600, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Paths"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		err = b.Put([]byte("/houston"), []byte("https://www.houstontx.gov/"))
		err = b.Put([]byte("/atlanta"), []byte("https://www.atlantaga.gov/"))
		err = b.Put([]byte("/nashville"), []byte("https://www.visitmusiccity.com/"))
		err = b.Put([]byte("/birmingham"), []byte("https://www.birminghamal.gov/"))
		err = b.Put([]byte("/austin"), []byte("https://www.austintexas.org/"))
		if err != nil {
			return fmt.Errorf("put item: %s", err)
		}
		return nil
	})
	return db, nil
}
