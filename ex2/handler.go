package urlshort

import (
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...
	return func(resp http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if val, ok := pathsToUrls[path]; ok {
			http.Redirect(resp, r, val, http.StatusSeeOther)
		} else {
			fallback.ServeHTTP(resp, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	urlPaths, err := parseUrlPath(yml)

	if err != nil {
		return nil, err
	}

	urlsMap := urlsPathsMap(urlPaths)

	return MapHandler(urlsMap, fallback), nil
}

func parseUrlPath(yml []byte) ([]urlPath, error) {
	var urlPaths []urlPath

	err := yaml.Unmarshal(yml, &urlPaths)

	return urlPaths, err
}

func urlsPathsMap(urlPaths []urlPath) map[string]string {
	result := make(map[string]string)

	for _, v := range urlPaths {
		result[v.Path] = v.Url
	}

	return result
}

type urlPath struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}
