package config

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

type spec struct {
	Aliases []Action `hcl:"alias,block"`
}

func (c *Config) load(path string) (spec, error) {
	parser := hclparse.NewParser()

	if err := c.walk(path, parser); err != nil {
		return spec{}, fmt.Errorf("cfg.walk(%s) -- %w", path, err)
	}

	files := make([]*hcl.File, 0, len(parser.Files()))
	for _, file := range parser.Files() {
		files = append(files, file)
	}

	var s spec
	diags := gohcl.DecodeBody(
		hcl.MergeFiles(files), nil, &s)
	if diags.HasErrors() {
		return spec{}, diags
	}

	return s, nil
}
