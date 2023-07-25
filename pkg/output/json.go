package output

import (
	"bytes"
	"encoding/json"

	"github.com/alecthomas/chroma/quick"
)

// DisplayJSON Pretty print the JSON with highlight syntax.
func DisplayJSON(jsonString string) (string, error) {
	// Pretty print the JSON
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(jsonString), "", "  "); err != nil {
		return "", err
	}

	// Use chroma to highlight JSON syntax
	var highlightedJSON bytes.Buffer
	err := quick.Highlight(&highlightedJSON, prettyJSON.String(), "json", "terminal", "monokai")
	if err != nil {
		return "", err
	}

	// return the highlighted JSON to the terminal
	return highlightedJSON.String(), nil
}
