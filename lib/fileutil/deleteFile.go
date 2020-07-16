package fileutil

import "os"

// DeleteFile : Delete a file from given fileLocation
func DeleteFile(fileLocation string) error {
	err := os.Remove(fileLocation)
	if err != nil {
		return err
	}

	return nil
}
