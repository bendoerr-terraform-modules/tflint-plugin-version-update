package main

import (
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/spf13/afero"

	"github.com/bendoerr-terraform-modules/tflint-plugin-version-update/pkg/github"
	"github.com/bendoerr-terraform-modules/tflint-plugin-version-update/pkg/tflint"
)

type Config struct {
	Freeze bool   `name:"freeze"`
	Path   string `name:"path" arg:"" type:"path"`
}

func main() {
	var err error
	var cfg Config

	_ = kong.Parse(&cfg)

	tflFile, err := tflint.OpenConfig(afero.Afero{Fs: afero.NewOsFs()}, cfg.Path)
	if err != nil {
		panic(err)
	}

	tflData, err := tflint.NewData(tflFile)
	if err != nil {
		panic(err)
	}

	tflHcl, err := tflData.ParseForRead()
	if err != nil {
		panic(err)
	}

	tflHclW, err := tflData.ParseForWrite()
	if err != nil {
		panic(err)
	}

	plugins, err := tflint.FindPluginVersions(tflHcl)
	if err != nil {
		panic(err)
	}

	runUpdate(plugins, cfg, tflHclW)

	_, _ = tflHclW.WriteTo(os.Stdout)
}

func runUpdate(plugins []*tflint.PluginConfig, cfg Config, tflHclW *hclwrite.File) {
	for _, plugin := range plugins {
		latestVersion, err := github.LatestVersion(plugin.SourceOwner, plugin.SourceRepo)
		if err != nil {
			panic(err)
		}

		if cfg.Freeze {
			if plugin.Version == latestVersion.ReleaseSHA {
				continue
			}
		} else {
			if plugin.Version == latestVersion.ReleaseTag || "v"+plugin.Version == latestVersion.ReleaseTag {
				continue
			}
		}

		if cfg.Freeze {
			err = tflint.UpdatePluginVersion(plugin.Name, latestVersion.ReleaseSHA, latestVersion.ReleaseTag, tflHclW)
			if err != nil {
				panic(err)
			}
		} else {
			// Stylistically tflint drops the 'v' in their documentation,
			// so we'll follow that as well.
			version := strings.TrimPrefix(latestVersion.ReleaseTag, "v")
			err = tflint.UpdatePluginVersion(plugin.Name, version, "", tflHclW)
			if err != nil {
				panic(err)
			}
		}
	}
}
