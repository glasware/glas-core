package config

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hclparse"
)

func (c Config) walk(abs string, parser *hclparse.Parser) error {
	return filepath.WalkDir(abs, func(path string, d fs.DirEntry, err error) error {
		if err != nil /*&& !errors.Is(err, os.ErrNotExist)*/ {
			return err
		}

		if d.Name() == filepath.Base(abs) {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		switch {
		case strings.HasSuffix(d.Name(), ".hcl"):
			_, err = parser.ParseHCLFile(path)
		case strings.HasSuffix(d.Name(), ".json"):
			_, err = parser.ParseJSONFile(path)
		}
		return err
	})
}
