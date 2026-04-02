package output

import "fmt"

// Render outputs data in the appropriate format based on flags.
// Returns true if output was rendered, false if no format matched.
func Render(raw string, jsonFlag, yamlFlag bool, outputFile string, command string) {
	if jsonFlag {
		json, _ := DisplayJSON(raw)
		fmt.Println(json)
	} else if yamlFlag {
		yaml, _ := DisplayYAML(raw)
		fmt.Println(yaml)
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
