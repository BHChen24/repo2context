package core

import (
	"fmt"
	"os"
	"path/filepath"
)

// Run is the main function for the repo2context application logic.
func Run(paths []string) error {
	for _, path := range paths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("error getting absolute path for '%s': %w", path, err)
		}
		// Check if the path exists
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", absPath)
		}
		fmt.Printf("Read path at: %s\n", absPath)
	}
	return nil
}
