package lib

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePath(t *testing.T) {
	testcases := []struct {
		input        string
		expectedPath Path
		expectedErr  error
	}{
		{
			"",
			Path{OriginalString: ""},
			errors.New("Path cannot be empty"),
		},
		{
			"hi/there",
			Path{OriginalString: "hi/there"},
			errors.New("Path must begin with slash"),
		},
		{
			"/hi//there",
			Path{OriginalString: "/hi//there"},
			errors.New("Invalid path format"),
		},
		{
			"/hi/there/",
			Path{OriginalString: "/hi/there/"},
			errors.New("Invalid path format"),
		},
		{
			"/",
			Path{OriginalString: "/", TableName: "", RemainingSegments: []string{}},
			nil,
		},
		{
			"/hi/there",
			Path{
				OriginalString:    "/hi/there",
				TableName:         "hi",
				RemainingSegments: []string{"there"}},
			nil,
		},
		{
			"/hi/there/sup",
			Path{
				OriginalString:    "/hi/there/sup",
				TableName:         "hi",
				RemainingSegments: []string{"there", "sup"}},
			nil,
		},
	}

	for _, testcase := range testcases {
		path, err := ParsePath(testcase.input)
		assert.Equal(t, testcase.expectedPath, path)
		assert.Equal(t, testcase.expectedErr, err)

		// Test Path.IsRoot.
		if testcase.expectedErr == nil {
			if testcase.input == "/" {
				assert.True(t, path.IsRoot())
			} else {
				assert.False(t, path.IsRoot())
			}
		}
	}
}
