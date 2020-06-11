package main

import (
	"fmt"
	"golang.org/x/tools/go/packages"
	"sort"
)

type context struct {
	soFar map[string]struct{}
}

func isInPackages(id string) bool {
	for pkg := range packagesList {
		if pkg == id {
			return true
		}
	}
	return false
}

func (c *context) find(name string) (err error) {
	if name == "C" {
		return nil
	}
	var pkg *packages.Package
	pkg, ok := cache[name]
	if !ok {
		cfg := &packages.Config{
			Mode: packages.NeedImports,
		}
		if modDirectory != nil && *modDirectory != "" {
			cfg.BuildFlags = []string{fmt.Sprintf("-mod=%s", *modDirectory)}
		}
		pkgs, err := packages.Load(cfg, name)
		if err != nil {
			return err
		}
		if len(pkgs) != 1 {
			return fmt.Errorf("expected 1 package but got %d", len(pkgs))
		}
		pkg = pkgs[0]
	}
	cache[name] = pkg
	if !isInPackages(pkg.ID) {
		return nil
	}

	if name != "." {
		c.soFar[pkg.ID] = struct{}{}
	}
	imports := pkg.Imports

	for _, imp := range imports {
		if _, ok := c.soFar[imp.ID]; !ok {
			if err := c.find(imp.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

func findDeps(name string) ([]string, error) {
	c := &context{
		soFar: make(map[string]struct{}),
	}
	if err := c.find(name); err != nil {
		return nil, err
	}
	deps := make([]string, 0, len(c.soFar))
	for p := range c.soFar {
		deps = append(deps, p)
	}
	sort.Strings(deps)
	return deps, nil
}
