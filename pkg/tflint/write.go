package tflint

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func UpdatePluginVersion(name, version, comment string, file *hclwrite.File) error {
	for _, block := range file.Body().Blocks() {
		if block.Type() == "plugin" && block.Labels()[0] == name {
			if comment == "" {
				block.Body().SetAttributeValue("version", cty.StringVal(version))
			} else {
				tokens := append(
					hclwrite.TokensForValue(cty.StringVal(version)),
					&hclwrite.Token{
						Type:  hclsyntax.TokenComment,
						Bytes: []byte("# " + comment),
					},
				)
				block.Body().SetAttributeRaw("version", tokens)
			}
		}
	}
	return nil
}
