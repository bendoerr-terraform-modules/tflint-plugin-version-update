package tflint

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/afero"

	"github.com/bendoerr-terraform-modules/tflint-plugin-version-update/pkg/ui"
)

// OpenConfig loads a TFLint config file following the logic from tflint except
// it will not auto load from the home directory, to update that provide the
// path explicitly.
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
func OpenConfig(ctx context.Context, fs afero.Afero, file string) (afero.File, error) {
	ui.Update(ctx, "🔧 Finding config")

	// Load the file passed by the --config option
	if file != "" {
		ui.Info(ctx, "🔧 Using provided config "+file)
		f, err := fs.OpenFile(file, os.O_RDWR, os.ModePerm)
		if err != nil {
			ui.Error(ctx, "🚨 Couldn't open provided config file!")
			return nil, fmt.Errorf("unable to open config='%s': %w", file, err)
		}
		return f, nil
	}

	// Load the file set by the environment variable
	envFile := os.Getenv("TFLINT_CONFIG_FILE")
	if envFile != "" {
		ui.Info(ctx, "🔧 Using env.TFLINT_CONFIG_FILE "+envFile)
		f, err := fs.OpenFile(envFile, os.O_RDWR, os.ModePerm)
		if err != nil {
			ui.Error(ctx, "🚨 Couldn't open env.TFLINT_CONFIG_FILE")
			return nil, fmt.Errorf("unable to open TFLINT_CONFIG_FILE='%s': %w", envFile, err)
		}
		return f, nil
	}

	// Load the default config file
	var defaultConfigFile = ".tflint.hcl"
	ui.Info(ctx, "🔧 Using default config "+defaultConfigFile)
	f, err := fs.OpenFile(defaultConfigFile, os.O_RDWR, os.ModePerm)
	if err != nil {
		ui.Error(ctx, "🚨 Couldn't open default config")
		return nil, errors.New("no config file found")
	}

	ui.Update(ctx, "")
	return f, nil
}
