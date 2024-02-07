package urlshort

import (
	"fmt"
	"net/http"

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
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	fmt.Println("First: ", pathMap["/urlshort"])
	fmt.Println("Second: ", pathMap["/urlshort-final"])
	return MapHandler(pathMap, fallback), nil
}

func parseYAML(yml []byte) ([]redirects, error) {
	var s []redirects
	err := yaml.Unmarshal(yml, &s)
	if err != nil {
		fmt.Println("Unmarshalling YAML: ", err)
	}
	return s, nil
}

func buildMap(yml []redirects) map[string]string {
	m := make(map[string]string)
	for _, v := range yml {
		m[v.Path] = v.URL
	}
	return m
}
