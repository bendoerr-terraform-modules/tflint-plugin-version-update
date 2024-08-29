package main

import (
	"context"
	"os"

	"github.com/alecthomas/kong"
	"github.com/spf13/afero"

	"github.com/bendoerr-terraform-modules/tflint-plugin-version-update/pkg/github"
	"github.com/bendoerr-terraform-modules/tflint-plugin-version-update/pkg/tflint"
	"github.com/bendoerr-terraform-modules/tflint-plugin-version-update/pkg/ui"
)

type Config struct {
	Freeze  bool   `name:"freeze"`
	Verbose bool   `name:"verbose"`
	Path    string `name:"path" arg:"" type:"path" optional:"true"`
}

func main() {
	var err error
	var cfg Config

	_ = kong.Parse(&cfg)

	ctx := context.Background()
	ctx = ui.ToContext(ctx, ui.NewUI(os.Stdout, cfg.Verbose))

	tflFile, err := tflint.OpenConfig(ctx, afero.Afero{Fs: afero.NewOsFs()}, cfg.Path)
	if err != nil {
		if cfg.Verbose {
			panic(err)
		}
		os.Exit(1)
	}

	tfl := tflint.NewTFLint(tflFile)

	err = tfl.ParseHCL(ctx)
	if err != nil {
		if cfg.Verbose {
			panic(err)
		}
		os.Exit(1)
	}

	err = tfl.UpdatePlugins(ctx, cfg.Freeze, github.LatestVersion)
	if err != nil {
		if cfg.Verbose {
			panic(err)
		}
		os.Exit(1)
	}

	err = tfl.Write(ctx)
	if err != nil {
		if cfg.Verbose {
			panic(err)
		}
		os.Exit(1)
	}

	ui.Stop(ctx)
	ui.Info(ctx, "✨ Done")
}
