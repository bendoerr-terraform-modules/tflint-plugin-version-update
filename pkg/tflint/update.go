package tflint

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/spf13/afero"

	"github.com/bendoerr-terraform-modules/tflint-plugin-version-update/pkg/ui"
)

type TFLint struct {
	file     afero.File
	hclread  *hcl.File
	hclwrite *hclwrite.File
	plugins  []*PluginConfig
	updated  bool
}

func NewTFLint(file afero.File) *TFLint {
	return &TFLint{file: file}
}

func (tfl *TFLint) ParseHCL(ctx context.Context) error {
	ui.Update(ctx, "🔍️ Parsing HCL")

	data, err := NewData(tfl.file)
	if err != nil {
		ui.Error(ctx, "🚨 Failed to read HCL config!")
		return err
	}

	tfl.hclread, err = data.ParseForRead()
	if err != nil {
		ui.Error(ctx, "🚨 Failed to parse HCL config!")
		return err
	}

	tfl.hclwrite, err = data.ParseForWrite()
	if err != nil {
		ui.Error(ctx, "🚨 Failed to parse HCL config!")
		return err
	}

	tfl.plugins, err = FindPluginVersions(tfl.hclread)
	if err != nil {
		ui.Error(ctx, "🚨 Failed to find plugin config!")
		return err
	}

	ui.Update(ctx, "")
	return nil
}

type LatestVersionFunc func(owner, repo string) (tag, sha, desc string, err error)

func (tfl *TFLint) UpdatePlugins(ctx context.Context, freezeSHA bool, latestVersionFunc LatestVersionFunc) error {
	if tfl.plugins == nil {
		panic("programmer error: parse the HCL first")
	}

	ui.Update(ctx, "🚀 Updating plugins to the latest versions")

	for _, plugin := range tfl.plugins {
		latestTag, latestSHA, _, err := latestVersionFunc(plugin.SourceOwner, plugin.SourceRepo)
		if err != nil {
			ui.Error(ctx, "🚨 Failed get latest plugin version!")
			return err
		}

		if isUpToDate(plugin.Version, latestTag, latestSHA, freezeSHA) {
			ui.Info(ctx, fmt.Sprintf("✅ %s/%s@%s", plugin.SourceOwner, plugin.SourceRepo, plugin.Version))
			continue
		}

		var newVersion string
		var newComment string
		if freezeSHA {
			newVersion = latestSHA
			newComment = latestTag
		} else {
			// Stylistically tflint drops the 'v' in their documentation,
			// so we'll follow that as well.
			newVersion = strings.TrimPrefix(latestTag, "v")
		}

		err = UpdatePluginVersion(plugin.Name, newVersion, newComment, tfl.hclwrite)
		if err != nil {
			ui.Error(ctx, "🚨 Failed to update the HCL plugin version")
			return err
		}
		tfl.updated = true

		ui.Info(ctx, fmt.Sprintf("⬆️ %s/%s@%s → %s", plugin.SourceOwner, plugin.SourceRepo, plugin.Version, newVersion))
	}

	ui.Update(ctx, "")
	return nil
}

func (tfl *TFLint) Write(ctx context.Context) error {
	if tfl.hclwrite == nil {
		panic("programmer error: parse the HCL first")
	}

	if !tfl.updated {
		return nil
	}

	ui.Update(ctx, "✏️ Writing new TFLint config")

	_ = tfl.file.Truncate(0)
	_, _ = tfl.file.Seek(0, 0)

	_, err := tfl.hclwrite.WriteTo(tfl.file)
	if err != nil {
		ui.Error(ctx, "🚨 Failed to write new TFLint config!")
		return err
	}

	ui.Update(ctx, "")
	return nil
}

func isUpToDate(currentVersion, latestTag, latestSHA string, useSHA bool) bool {
	return (useSHA && currentVersion == latestSHA) ||
		(!useSHA && (currentVersion == latestTag || "v"+currentVersion == latestTag))
}
