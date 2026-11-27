// Command rename changes the Go module path of this template and rewrites
// every import that uses it. Run it once after cloning the template:
//
//	make rename NAME=github.com/you/new_project
//	go run ./tools/rename -name github.com/you/new_project
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// skipDirs are never scanned for Go sources.
var skipDirs = map[string]bool{
	".git":         true,
	"node_modules": true,
	"webui":        true,
	"build":        true,
	".next":        true,
}

var validModuleChars = regexp.MustCompile(`^[A-Za-z0-9._~\-/]+$`)

// ValidateModulePath performs a conservative sanity check on a module path.
func ValidateModulePath(p string) error {
	if p == "" {
		return errors.New("module path must not be empty")
	}
	if !validModuleChars.MatchString(p) {
		return fmt.Errorf("module path %q contains invalid characters (allowed: letters, digits, . _ ~ - /)", p)
	}
	if strings.HasPrefix(p, "/") || strings.HasSuffix(p, "/") {
		return fmt.Errorf("module path %q must not start or end with '/'", p)
	}
	for _, elem := range strings.Split(p, "/") {
		if elem == "" || elem == "." || elem == ".." {
			return fmt.Errorf("module path %q has an invalid element %q", p, elem)
		}
	}
	return nil
}

// ReadModulePath returns the module path declared in root/go.mod.
func ReadModulePath(root string) (string, error) {
	f, err := os.Open(filepath.Join(root, "go.mod"))
	if err != nil {
		return "", err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if rest, ok := strings.CutPrefix(line, "module "); ok {
			return strings.Trim(strings.TrimSpace(rest), `"`), nil
		}
	}
	if err := sc.Err(); err != nil {
		return "", err
	}
	return "", errors.New("no module directive found in go.mod")
}

// Rename rewrites the module path in root/go.mod and in every import string of
// the .go files below root. It returns the old module path and the list of
// files (relative to root) that were modified.
func Rename(root, newModule string) (string, []string, error) {
	if err := ValidateModulePath(newModule); err != nil {
		return "", nil, err
	}
	oldModule, err := ReadModulePath(root)
	if err != nil {
		return "", nil, err
	}
	if oldModule == newModule {
		return oldModule, nil, nil
	}

	var changed []string

	// go.mod
	goModPath := filepath.Join(root, "go.mod")
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return oldModule, nil, err
	}
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "module ") {
			lines[i] = "module " + newModule
			break
		}
	}
	if err := os.WriteFile(goModPath, []byte(strings.Join(lines, "\n")), 0o644); err != nil {
		return oldModule, nil, err
	}
	changed = append(changed, "go.mod")

	// Go sources: only touch quoted import paths, i.e. "old" or "old/...".
	replacer := strings.NewReplacer(
		`"`+oldModule+`"`, `"`+newModule+`"`,
		`"`+oldModule+`/`, `"`+newModule+`/`,
	)
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			if path != root && skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		out := replacer.Replace(string(src))
		if out == string(src) {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(out), info.Mode().Perm()); err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		changed = append(changed, filepath.ToSlash(rel))
		return nil
	})
	return oldModule, changed, err
}

func main() {
	name := flag.String("name", "", "new Go module path, e.g. github.com/you/new_project")
	root := flag.String("root", ".", "project root containing go.mod")
	flag.Parse()

	if *name == "" {
		fmt.Fprintln(os.Stderr, "usage: go run ./tools/rename -name <new/module/path>")
		os.Exit(2)
	}

	oldModule, changed, err := Rename(*root, *name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "rename failed:", err)
		os.Exit(1)
	}
	if len(changed) == 0 {
		fmt.Printf("Module is already %q, nothing to do.\n", *name)
		return
	}
	fmt.Printf("Module renamed: %s -> %s\n", oldModule, *name)
	for _, f := range changed {
		fmt.Println("  updated", f)
	}
}
