package urlshort

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	var a http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		redirect, ok := pathsToUrls[url]
		if !ok {
			fallback.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, redirect, http.StatusMovedPermanently)
	}
	return a
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
type redirects struct {
	Path string
	URL  string
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// TODO: Implement this...

	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, fmt.Errorf("Parsing YAML: %s", err)
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func parseYAML(yml []byte) ([]redirects, error) {
	var s []redirects
	err := yaml.Unmarshal(yml, &s)
	if err != nil {
		return nil, fmt.Errorf("Unmarshaling YAML: %s", err)
	}
	return s, nil
}

func buildMap(r []redirects) map[string]string {
	m := make(map[string]string)
	for _, v := range r {
		m[v.Path] = v.URL
	}
	return m
}

func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJson, err := parseJSON(json)
	if err != nil {
		return nil, fmt.Errorf("Parsing JSON: %s", err)
	}
	pathMap := buildMap(parsedJson)
	return MapHandler(pathMap, fallback), nil
}

func parseJSON(jsn []byte) ([]redirects, error) {
	var r []redirects
	err := json.Unmarshal(jsn, &r)
	if err != nil {
		return nil, fmt.Errorf("Unmarshaling JSON: %s", err)
	}
	return r, nil
}

func BoltHandler(fallback http.Handler) (http.HandlerFunc, error) {
	// get all of the key values from the database and put them into a map
	db, err := bolt.Open(filepath.Join("..", "data", "paths.db"), 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("Opening db: %s", err)
	}
	defer db.Close()

	pathMap := make(map[string]string)

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Paths"))

		if err := b.ForEach(func(k, v []byte) error {
			pathMap[string(k)] = string(v)
			return nil
		}); err != nil {
			return fmt.Errorf("Getting db items: %s", err)
		}
		return nil
	})

	return MapHandler(pathMap, fallback), nil
}
