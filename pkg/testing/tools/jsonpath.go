package tools

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

// TestJSONPath wrapper around https://github.com/tidwall/gjson
func TestJSONPath(t *testing.T, expected interface{}, jsonPath, json string) {
	assert.Equal(t, expected, gjson.Get(json, jsonPath).String(), fmt.Sprintf("jsonPath %q expected to be %q", jsonPath, expected))
}

func TestJSONPathNotEmpty(t *testing.T, jsonPath, json string) {
	assert.NotEmpty(t, gjson.Get(json, jsonPath).String(), fmt.Sprintf("jsonPath %q expected to be not empty", jsonPath))
}
