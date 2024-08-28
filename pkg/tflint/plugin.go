package tflint

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
)

// PluginConfig is a TFLint's plugin config.
type PluginConfig struct {
	Name       string `hcl:"name,label"`
	Enabled    bool   `hcl:"enabled"`
	Version    string `hcl:"version,optional"`
	Source     string `hcl:"source,optional"`
	SigningKey string `hcl:"signing_key,optional"`

	Body hcl.Body `hcl:",remain"`

	// Parsed source attributes
	SourceHost  string
	SourceOwner string
	SourceRepo  string
}

func (c *PluginConfig) Validate() error {
	if c.Version != "" && c.Source == "" {
		return fmt.Errorf(`plugin "%s": "source" attribute cannot be omitted when specifying "version"`, c.Name)
	}

	if c.Source != "" {
		if c.Version == "" {
			return fmt.Errorf(`plugin "%s": "version" attribute cannot be omitted when specifying "source"`, c.Name)
		}

		parts := strings.Split(c.Source, "/")
		if len(parts) != 3 { //nolint:mnd // Expected `github.com/owner/repo` format
			return fmt.Errorf(`plugin "%s": "source" is invalid. Must be in the format "${host}/${owner}/${repo}"`, c.Name)
		}

		c.SourceHost = parts[0]
		c.SourceOwner = parts[1]
		c.SourceRepo = parts[2]
	}

	return nil
}
