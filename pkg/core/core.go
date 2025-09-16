package core

import (
	"fmt"
	"path/filepath"
)

// Run is the main function for the repo2context application logic.
func Run(paths []string) error {
	for _, path := range paths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("error getting absolute path for '%s': %w", path, err)
		}
		fmt.Printf("Read path at: %s\n", absPath)
	}
	return nil
}
