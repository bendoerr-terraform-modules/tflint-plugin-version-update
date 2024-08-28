package tflint

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/spf13/afero"
)

type Data struct {
	Bytes    []byte
	Filename string
}

func NewData(file afero.File) (*Data, error) {
	d := &Data{}
	var err error
	d.Filename = file.Name()
	d.Bytes, err = afero.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Data) ParseForRead() (*hcl.File, error) {
	f, diag := hclsyntax.ParseConfig(d.Bytes, d.Filename, hcl.InitialPos)
	if diag.HasErrors() {
		return nil, diag
	}
	return f, nil
}

func (d *Data) ParseForWrite() (*hclwrite.File, error) {
	f, diag := hclwrite.ParseConfig(d.Bytes, d.Filename, hcl.InitialPos)
	if diag.HasErrors() {
		return nil, diag
	}
	return f, nil
}
