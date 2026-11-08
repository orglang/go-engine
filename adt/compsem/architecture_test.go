package compsem

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestADTCoreDoesNotDependOnInfrastructure(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not locate architecture test source")
	}
	adtRoot := filepath.Dir(filepath.Dir(currentFile))

	err := filepath.WalkDir(adtRoot, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if entry.Name() == "adapter" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		file, parseErr := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if parseErr != nil {
			return parseErr
		}
		for _, spec := range file.Imports {
			importPath, unquoteErr := strconv.Unquote(spec.Path.Value)
			if unquoteErr != nil {
				return unquoteErr
			}
			if importPath == "github.com/jackc/pgx/v5" ||
				strings.HasPrefix(importPath, "github.com/jackc/pgx/v5/") ||
				importPath == "github.com/huandu/go-sqlbuilder" ||
				strings.Contains(importPath, "/adapter") {
				t.Errorf("%s imports infrastructure package %q; move the implementation behind an adapter", path, importPath)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("scan ADT package dependencies: %v", err)
	}
}
