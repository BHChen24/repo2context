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
			fmt.Fprintf(os.Stderr, "error getting absolute path for '%s': %v\n", path, err)
			continue
		}
		// Check if the path exists
		if _, err := os.Stat(absPath); err != nil {
			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "path does not exist: %s\n", absPath)
				continue
			}
			// TODO: Handle other errors
			fmt.Fprintf(os.Stderr, "error checking path '%s': %v\n", absPath, err)
			continue
		}
		fmt.Printf("Read path at: %s\n", absPath)
	}
	return nil
}
