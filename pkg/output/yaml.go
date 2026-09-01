package output

import (
	"bytes"
	"encoding/json"

	"github.com/alecthomas/chroma/quick"
	"gopkg.in/yaml.v2"
)

// DisplayYAML output YAML format.
func DisplayYAML(jsonString string) (string, error) {
	var data interface{}

	// Decode JSON into a Go structure
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		return "", err
	}

	// Encode Go structure into YAML
	yamlBytes, err := yaml.Marshal(data)
	if err != nil {
		return "", err
	}

	// Use chroma to highlight YAML syntax
	var highlightedYAML bytes.Buffer
	er := quick.Highlight(&highlightedYAML, string(yamlBytes), "yaml", "terminal", "monokai")
	if er != nil {
		return "", er
	}

	return highlightedYAML.String(), nil
}
