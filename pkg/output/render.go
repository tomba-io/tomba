package output

import "fmt"

// CSVFlag is set globally from the --csv flag
var CSVFlag bool

// Render outputs data in the appropriate format based on flags.
func Render(raw string, jsonFlag, yamlFlag bool, outputFile string, command string) {
	if jsonFlag {
		json, _ := DisplayJSON(raw)
		fmt.Println(json)
	} else if yamlFlag {
		yaml, _ := DisplayYAML(raw)
		fmt.Println(yaml)
	} else if CSVFlag {
		csv, _ := DisplayCSV(raw)
		fmt.Print(csv)
	} else {
		text := DisplayText(raw, command)
		fmt.Print(text)
	}
	if outputFile != "" {
		err := CreateOutput(outputFile, raw)
		if err != nil {
			fmt.Println("Error creating file:", err)
		}
	}
}
