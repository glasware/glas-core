package config

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/spf13/afero"
)

func (c *Config) walk(abs string, parser *hclparse.Parser) error {
	ok, err := afero.DirExists(c.afs, abs)
	if err != nil {
		return fmt.Errorf("afero.DirExists(%s) -- %w", abs, err)
	}

	if !ok {
		return nil
	}

	return afero.Walk(c.afs, abs, func(path string, info fs.FileInfo, err error) error {
		if err != nil /*&& !errors.Is(err, os.ErrNotExist)*/ {
			return err
		}

		if info.Name() == filepath.Base(abs) {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		byt, err := afero.ReadFile(c.afs, path)
		if err != nil {
			return fmt.Errorf("afero.ReadFile -- %w", err)
		}

		var diags hcl.Diagnostics
		switch {
		case strings.HasSuffix(info.Name(), ".hcl"):
			_, diags = parser.ParseHCL(byt, path)
		case strings.HasSuffix(info.Name(), ".json"):
			_, diags = parser.ParseJSON(byt, path)
		}
		if diags.HasErrors() {
			return diags
		}
		return nil
	})
}
