package output

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
)

// DisplayCSV output CSV format.
func DisplayCSV(raw string) (string, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Find the array data - could be under "data" key or at top level
	var items []interface{}
	if d, ok := data["data"]; ok {
		switch v := d.(type) {
		case []interface{}:
			items = v
		case map[string]interface{}:
			// Single object, wrap in array
			items = []interface{}{v}
		}
	}

	// If no "data" key found, try the top-level object itself
	if items == nil {
		items = []interface{}{data}
	}

	if len(items) == 0 {
		return "", nil
	}

	// Extract headers from the first item
	firstItem, ok := items[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected data format")
	}

	headers := make([]string, 0, len(firstItem))
	for k := range firstItem {
		headers = append(headers, k)
	}
	sort.Strings(headers)

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write rows
	for _, item := range items {
		row, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		record := make([]string, len(headers))
		for i, h := range headers {
			if val, exists := row[h]; exists {
				record[i] = fmt.Sprintf("%v", val)
			}
		}
		if err := writer.Write(record); err != nil {
			return "", fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("CSV writer error: %w", err)
	}

	return buf.String(), nil
}
