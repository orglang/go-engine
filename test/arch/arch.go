package arch

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

type Kind int

const (
	Unknown Kind = iota
	Agnostic
	Specific
	Framework
	Shared
)

func (k Kind) String() string {
	switch k {
	case Agnostic:
		return "agnostic core"
	case Specific:
		return "toolkit-specific adapter"
	case Framework:
		return "framework layer"
	case Shared:
		return "shared kernel"
	default:
		return "not a hexagon"
	}
}

const modulePrefix = "orglang/go-engine/"

func Classify(pkgPath string) Kind {
	rest, ok := strings.CutPrefix(pkgPath, modulePrefix)
	if !ok {
		return Unknown
	}
	if rest == "adt" || strings.HasPrefix(rest, "adt/") {
		return Shared
	}
	if rest == "lib" || strings.HasPrefix(rest, "lib/") {
		return Framework
	}
	for _, seg := range strings.Split(rest, "/") {
		switch seg {
		case "core", "domain":
			return Agnostic
		case "adapter", "adapters":
			return Specific
		}
	}
	return Unknown
}

type Violation struct {
	From string
	To   string
}

func (v Violation) String() string {
	return fmt.Sprintf("%s imports %s", v.From, v.To)
}

func Check(imports map[string][]string) []Violation {
	var violations []Violation
	for from, tos := range imports {
		if Classify(from) != Agnostic {
			continue
		}
		for _, to := range tos {
			if Classify(to) == Specific {
				violations = append(violations, Violation{From: from, To: to})
			}
		}
	}
	sort.Slice(violations, func(i, j int) bool {
		if violations[i].From != violations[j].From {
			return violations[i].From < violations[j].From
		}
		return violations[i].To < violations[j].To
	})
	return violations
}

func Load(dirs ...string) (map[string][]string, error) {
	cfg := &packages.Config{Mode: packages.NeedName | packages.NeedImports}
	pkgs, err := packages.Load(cfg, dirs...)
	if err != nil {
		return nil, err
	}
	graph := make(map[string][]string, len(pkgs))
	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			return nil, fmt.Errorf("load %s: %v", pkg.PkgPath, pkg.Errors[0])
		}
		paths := make([]string, 0, len(pkg.Imports))
		for _, imported := range pkg.Imports {
			paths = append(paths, imported.PkgPath)
		}
		graph[pkg.PkgPath] = paths
	}
	return graph, nil
}
