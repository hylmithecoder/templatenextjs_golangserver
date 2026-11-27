package main

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestValidateModulePath(t *testing.T) {
	valid := []string{"new_project", "github.com/you/new_project", "example.com/a-b/c.d~e"}
	for _, p := range valid {
		if err := ValidateModulePath(p); err != nil {
			t.Errorf("ValidateModulePath(%q) unexpected error: %v", p, err)
		}
	}
	invalid := []string{"", "has space", "/leading", "trailing/", "a//b", "a/../b", "a/./b", `quo"te`}
	for _, p := range invalid {
		if err := ValidateModulePath(p); err == nil {
			t.Errorf("ValidateModulePath(%q) expected error, got nil", p)
		}
	}
}

func TestReadModulePath(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "go.mod"), "// comment\nmodule example.com/old\n\ngo 1.27\n")
	got, err := ReadModulePath(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got != "example.com/old" {
		t.Fatalf("got %q, want %q", got, "example.com/old")
	}

	empty := t.TempDir()
	writeFile(t, filepath.Join(empty, "go.mod"), "go 1.27\n")
	if _, err := ReadModulePath(empty); err == nil {
		t.Fatal("expected error for go.mod without module directive")
	}
}

func TestRenameRewritesGoModAndImports(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "go.mod"), "module oldmod\n\ngo 1.27\n")
	writeFile(t, filepath.Join(dir, "src", "main.go"), `package main

import (
	"fmt"

	"oldmod/src/server"
	alias "oldmod/src/utils"
)

// oldmod mentioned in a comment stays untouched; so does "oldmodule/x".
var _ = "oldmodule/x"

func main() { fmt.Println(server.X, alias.Y) }
`)
	writeFile(t, filepath.Join(dir, "src", "noimport.go"), "package main\n")
	// Files under skipped directories must not be touched.
	skipped := filepath.Join(dir, "webui", "node_modules", "pkg", "x.go")
	writeFile(t, skipped, `package x; import _ "oldmod/src/server"`)

	oldModule, changed, err := Rename(dir, "github.com/you/new_project")
	if err != nil {
		t.Fatal(err)
	}
	if oldModule != "oldmod" {
		t.Errorf("old module = %q, want oldmod", oldModule)
	}
	if strings.Join(changed, ",") != "go.mod,src/main.go" {
		t.Errorf("changed = %v, want [go.mod src/main.go]", changed)
	}

	if got := readFile(t, filepath.Join(dir, "go.mod")); !strings.HasPrefix(got, "module github.com/you/new_project\n\ngo 1.27") {
		t.Errorf("go.mod not rewritten correctly:\n%s", got)
	}
	main := readFile(t, filepath.Join(dir, "src", "main.go"))
	for _, want := range []string{
		`"github.com/you/new_project/src/server"`,
		`alias "github.com/you/new_project/src/utils"`,
		`// oldmod mentioned in a comment`,
		`"oldmodule/x"`,
	} {
		if !strings.Contains(main, want) {
			t.Errorf("main.go missing %q:\n%s", want, main)
		}
	}
	if strings.Contains(readFile(t, skipped), "new_project") {
		t.Error("file inside skipped directory was modified")
	}
}

func TestRenameSameNameIsNoop(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "go.mod"), "module same\n")
	_, changed, err := Rename(dir, "same")
	if err != nil {
		t.Fatal(err)
	}
	if len(changed) != 0 {
		t.Fatalf("expected no changes, got %v", changed)
	}
}

func TestRenameRejectsInvalidName(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "go.mod"), "module oldmod\n")
	if _, _, err := Rename(dir, "bad name"); err == nil {
		t.Fatal("expected error for invalid module path")
	}
	if got := readFile(t, filepath.Join(dir, "go.mod")); got != "module oldmod\n" {
		t.Fatalf("go.mod modified despite invalid name: %q", got)
	}
}

// TestRenameRealProjectBuilds copies this repository's Go sources into a temp
// dir, renames the module and verifies the result still builds and vets.
func TestRenameRealProjectBuilds(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in -short mode")
	}
	repoRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	dst := t.TempDir()
	for _, f := range []string{"go.mod", "go.sum"} {
		writeFile(t, filepath.Join(dst, f), readFile(t, filepath.Join(repoRoot, f)))
	}
	err = filepath.WalkDir(filepath.Join(repoRoot, "src"), func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(repoRoot, path)
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		writeFile(t, filepath.Join(dst, rel), string(b))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	const newName = "github.com/example/new_project"
	if _, _, err := Rename(dst, newName); err != nil {
		t.Fatal(err)
	}
	if got, _ := ReadModulePath(dst); got != newName {
		t.Fatalf("module = %q, want %q", got, newName)
	}

	for _, args := range [][]string{{"build", "./..."}, {"vet", "./..."}} {
		cmd := exec.Command("go", args...)
		cmd.Dir = dst
		cmd.Env = append(os.Environ(), "GOFLAGS=-buildvcs=false")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("go %s failed after rename: %v\n%s", strings.Join(args, " "), err, out)
		}
	}
}
