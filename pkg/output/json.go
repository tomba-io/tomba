package output

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/alecthomas/chroma/quick"
	"github.com/tomba-io/tomba/pkg/util"
	"golang.org/x/term"
)

// DisplayJSON Pretty print the JSON with highlight syntax.
// When output is piped (non-TTY) or --no-color is set, returns plain JSON.
func DisplayJSON(jsonString string) (string, error) {
	// Pretty print the JSON
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(jsonString), "", "  "); err != nil {
		return "", err
	}

	// Skip highlighting if not a terminal or no-color is set
	if util.NoColor || !term.IsTerminal(int(os.Stdout.Fd())) {
		return prettyJSON.String(), nil
	}

	// Use chroma to highlight JSON syntax
	var highlightedJSON bytes.Buffer
	err := quick.Highlight(&highlightedJSON, prettyJSON.String(), "json", "terminal", "monokai")
	if err != nil {
		return prettyJSON.String(), nil
	}

	return highlightedJSON.String(), nil
}
