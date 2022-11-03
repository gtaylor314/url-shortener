package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// grab the redirect URL from the map based on the path given by the request
		redirectURL := pathsToUrls[r.URL.Path]
		// if the provided path doesn't exist in our pathsToUrls map, the redirectURL will be an empty string
		// call the fallback http.Handler
		if redirectURL == "" {
			// call the fallback http.Handler
			// if fallback stores a ServeMux, ServeHTTP dispatches the request to the handler which most closely matches
			// the path in the request
			// if fallback stores a http.HandlerFunc, ServeHTTP calls the function with w and r
			fallback.ServeHTTP(w, r)
			return
		}
		// otherwise, the path does exist in our pathsToUrls map
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	})
}

// pathToUrls struct for unmarshalling yaml file
type pathsToUrls struct {
	Path string
	Url  string
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
// The only errors that can be returned all relate to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// parse YAML into a slice of pathsToUrl struct with Path and Url as elements
	parsedYAML, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	// create map from parsed YAML
	mappedYAML := buildMap(parsedYAML)
	return MapHandler(mappedYAML, fallback), nil
}

func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJSON, err := parseJSON(jsn)
	if err != nil {
		return nil, err
	}
	mappedJSON := buildMap(parsedJSON)
	return MapHandler(mappedJSON, fallback), nil
}

func parseYAML(yml []byte) ([]pathsToUrls, error) {
	// create slice of pathToUrls struct for unmarshalled yaml
	p := []pathsToUrls{}
	err := yaml.Unmarshal(yml, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func parseJSON(jsn []byte) ([]pathsToUrls, error) {
	// create slice of pathToUrls struct for unmarshalled json
	p := []pathsToUrls{}
	err := json.Unmarshal(jsn, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func buildMap(pathsToUrls []pathsToUrls) map[string]string {
	// create map to populate from pathToUrls
	m := make(map[string]string)
	for _, v := range pathsToUrls {
		m[v.Path] = v.Url
	}
	return m
}
