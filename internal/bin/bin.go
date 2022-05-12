package bin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Locate tries to find the specified binary on the machine.
// It will look within directories specified in the PATH environment variable then the current directory.
func Locate(bin string) (string, error) {
	// Try to find the binary within directories in PATH
	path, err := exec.LookPath(bin)
	if err == nil {
		return path, nil
	}

	// Try to find the binary within current directory
	return exec.LookPath(fmt.Sprintf("./%s", bin))
}

// Get returns the contents of the specified embedded binary
func Get(bin string) ([]byte, error) {
	return binFS.ReadFile(bin)
}

// Write writes the embedded binary to the current directory
func Write(bin string) error {
	return WriteAs(bin, filepath.Base(bin))
}

// WriteAs writes the embedded binary to the current directory as the new name
func WriteAs(bin string, newName string) error {
	return WriteTo(bin, filepath.Join(".", newName))
}

// WriteTo writes the embedded binary to the destination
func WriteTo(bin string, dest string) error {
	data, err := Get(bin)
	if err != nil {
		return err
	}

	return os.WriteFile(dest, data, 755)
}
