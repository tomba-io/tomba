package output

import "os"

// CreateOutput creates a new file with the given filename and writes the data to it.
func CreateOutput(filename string, data string) error {
	// Create a new file with the given filename
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the data to the file
	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}
