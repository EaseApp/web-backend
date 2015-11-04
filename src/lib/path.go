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
func (p Path) IsRoot() bool {
	return p.TopLevelDocName == ""
}
