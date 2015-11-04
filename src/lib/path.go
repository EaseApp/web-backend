package lib

import (
	"errors"
	"regexp"
	"strings"
)

// Path represents an application data path.
type Path struct {
	OriginalString    string
	TopLevelDocName   string
	RemainingSegments []string
}

var pathRegex = regexp.MustCompile("^(/|/[^/]+(/[^/]+)*)$")

// ParsePath parses a string into a path.
func ParsePath(pathStr string) (Path, error) {
	path := Path{OriginalString: pathStr}
	if len(pathStr) < 1 {
		return path, errors.New("Path cannot be empty")
	}
	if pathStr[0] != '/' {
		return path, errors.New("Path must begin with slash")
	}
	if !pathRegex.MatchString(pathStr) {
		return path, errors.New("Invalid path format")
	}
	segments := strings.Split(pathStr, "/")
	path.TopLevelDocName = segments[1]
	path.RemainingSegments = segments[2:]
	return path, nil
}

// IsRoot returns true iff a path refers to the root of all docs.
func (path Path) IsRoot() bool {
	return path.TopLevelDocName == ""
}

// ToNestedQuery converts the path to a nested map for a rethinkdb query.
// The map eventually points to the given data.
func (path Path) ToNestedQuery(data interface{}) map[string]interface{} {

	// Generate the nested data query.
	nestedDataQuery := make(map[string]interface{})

	if len(path.RemainingSegments) == 0 {
		nestedDataQuery["data"] = data
	} else {
		nestedDataQuery["data"] = make(map[string]interface{})
		lastNestedEntry := nestedDataQuery["data"].(map[string]interface{})
		for idx, segment := range path.RemainingSegments {
			// For the last part of the query, set it to the data, else nest further.
			if idx == len(path.RemainingSegments)-1 {
				lastNestedEntry[segment] = data
			} else {
				lastNestedEntry[segment] = make(map[string]interface{})
				lastNestedEntry = lastNestedEntry[segment].(map[string]interface{})
			}
		}
	}

	return nestedDataQuery
}
