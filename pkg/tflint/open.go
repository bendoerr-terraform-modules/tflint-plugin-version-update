package tflint

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/afero"
)

// OpenConfig loads a TFLint config file following the logic from tflint except
// it will not auto load from the home directory, to update that provide the
// path explictly.
//
// tflint's LoadConfig @ tflint/tflint/config.go
//
// The priority of the configuration files is as follows:
//
// 1. file passed by the --config option
// 2. file set by the TFLINT_CONFIG_FILE environment variable
// 3. current directory (./.tflint.hcl)
//
// For 1 and 2, if the file does not exist, an error will be returned immediately.
// If 3 fails then an error will be returned.
func OpenConfig(fs afero.Afero, file string) (afero.File, error) {
	// Load the file passed by the --config option
	if file != "" {
		log.Printf("[INFO] Load config: %s", file)
		f, err := fs.Open(file)
		if err != nil {
			return nil, fmt.Errorf("unable to open file='%s': %w", file, err)
		}
		return f, nil
	}

	// Load the file set by the environment variable
	envFile := os.Getenv("TFLINT_CONFIG_FILE")
	if envFile != "" {
		log.Printf("[INFO] Load config: %s", envFile)
		f, err := fs.Open(envFile)
		if err != nil {
			return nil, fmt.Errorf("unable to open TFLINT_CONFIG_FILE='%s': %w", envFile, err)
		}
		return f, nil
	}

	// Load the default config file
	var defaultConfigFile = ".tflint.hcl"
	log.Printf("[INFO] Load default config: %s", defaultConfigFile)
	if f, err := fs.Open(defaultConfigFile); err == nil {
		return f, nil
	}
	log.Printf("[INFO] Default config not found")

	return nil, errors.New("no config file found")
}
