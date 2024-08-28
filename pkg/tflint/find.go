package tflint

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
)

func FindPluginVersions(file *hcl.File) ([]*PluginConfig, error) {
	var plugins []*PluginConfig

	var configSchema = &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "plugin",
				LabelNames: []string{"name"},
			},
		},
	}

	content, diag := file.Body.Content(configSchema)
	if diag.HasErrors() {
		return nil, diag
	}

	for _, block := range content.Blocks {
		if block.Type == "plugin" {
			pluginConfig := &PluginConfig{Name: block.Labels[0]}
			if err := gohcl.DecodeBody(block.Body, nil, pluginConfig); err != nil {
				return nil, err
			}
			if err := pluginConfig.Validate(); err != nil {
				return nil, err
			}
			plugins = append(plugins, pluginConfig)
		}
	}

	return plugins, nil
}
