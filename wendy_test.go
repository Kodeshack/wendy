package wendy

import (
	"os"
	"path"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestFSGenerator_Generate(t *testing.T) {
	tmpdir := t.TempDir()

	g := &FSGenerator{OutputDir: tmpdir}

	err := g.Generate(
		PlainFile("README.md", "This is the package"),
		Dir("bin",
			Dir("cli",
				PlainFile("main.go", "package main"),
			),
		),
		Dir("pkg",
			PlainFile("README.md", "how to use this thing"),
			Dir("cli",
				PlainFile("cli.go", "package cli..."),
				PlainFile("run.go", "package cli...run..."),
			),
		),
	)

	assert.NoError(t, err)

	rootDir, err := os.ReadDir(tmpdir)
	assert.NoError(t, err)
	assert.Len(t, rootDir, 3)

	binDir, err := os.ReadDir(path.Join(tmpdir, "bin"))
	assert.NoError(t, err)
	assert.Len(t, binDir, 1)

	pkgDir, err := os.ReadDir(path.Join(tmpdir, "pkg"))
	assert.NoError(t, err)
	assert.Len(t, pkgDir, 2)

	pkgCliDir, err := os.ReadDir(path.Join(tmpdir, "pkg", "cli"))
	assert.NoError(t, err)
	assert.Len(t, pkgCliDir, 2)

	readme, err := os.ReadFile(path.Join(tmpdir, "README.md"))
	assert.NoError(t, err)
	assert.Equal(t, "This is the package", string(readme))

	binCliMainGo, err := os.ReadFile(path.Join(tmpdir, "bin", "cli", "main.go"))
	assert.NoError(t, err)
	assert.Equal(t, "package main", string(binCliMainGo))

	pkgReadme, err := os.ReadFile(path.Join(tmpdir, "pkg", "README.md"))
	assert.NoError(t, err)
	assert.Equal(t, "how to use this thing", string(pkgReadme))

	pkgCliCliGo, err := os.ReadFile(path.Join(tmpdir, "pkg", "cli", "cli.go"))
	assert.NoError(t, err)
	assert.Equal(t, "package cli...", string(pkgCliCliGo))

	pkgCliRunGo, err := os.ReadFile(path.Join(tmpdir, "pkg", "cli", "run.go"))
	assert.NoError(t, err)
	assert.Equal(t, "package cli...run...", string(pkgCliRunGo))
}

func TestFSGenerator_Generate_ErrorOnExistingDir(t *testing.T) {
	tmpdir := t.TempDir()

	g := &FSGenerator{OutputDir: tmpdir}

	err := g.Generate(Dir("will_exist"))
	assert.NoError(t, err)

	err = g.Generate(Dir("will_exist"))
	assert.NoError(t, err)

	g = &FSGenerator{OutputDir: tmpdir, ErrorOnExistingDir: true}

	err = g.Generate(Dir("will_exist"))
	assert.Error(t, err)
}

func TestFSGenerator_Generate_ErrorOnExistingFile(t *testing.T) {
	tmpdir := t.TempDir()

	g := &FSGenerator{OutputDir: tmpdir, ErrorOnExistingFile: true}

	err := g.Generate(Dir("will_exist", PlainFile("test", "contents")))
	assert.NoError(t, err)

	err = g.Generate(Dir("will_exist", PlainFile("test", "contents")))
	assert.ErrorIs(t, err, os.ErrExist)
}

func TestFSGenerator_Generate_CleanDir(t *testing.T) {
	tmpdir := t.TempDir()

	g := &FSGenerator{
		OutputDir:          tmpdir,
		CleanDir:           true,
		ErrorOnExistingDir: true,
	}

	err := g.Generate(Dir("will_exist"))
	assert.NoError(t, err)

	err = g.Generate(Dir("will_exist"))
	assert.NoError(t, err)

	err = g.Generate(Dir("different_dir"))
	assert.NoError(t, err)

	entries, err := os.ReadDir(tmpdir)
	assert.NoError(t, err)

	assert.Len(t, entries, 1)
	assert.Equal(t, "different_dir", entries[0].Name())
	assert.True(t, entries[0].IsDir())
}

func TestFSGenerator_Generate_CleanDir_DirNotExist(t *testing.T) {
	tmpdir := t.TempDir()

	g := &FSGenerator{
		OutputDir: path.Join(tmpdir, "non_existent"),
		CleanDir:  true,
	}

	err := g.Generate(PlainFile("test.txt", "contents"))
	assert.NoError(t, err)
}

func TestFSGenerator_Generate_FileFromTemplate(t *testing.T) {
	tmpdir := t.TempDir()

	g := &FSGenerator{ //nolint:varnamelen // This is just a test
		OutputDir:          tmpdir,
		CleanDir:           true,
		ErrorOnExistingDir: true,
	}

	data := map[string]any{
		"foo": "bar",
		"baz": "bat",
	}

	tmplt, err := template.New("").Parse(`{
		"foo": "{{ .foo }}",
		"baz": "{{ .baz }}",
}`)
	assert.NoError(t, err)

	err = g.Generate(
		FileFromTemplate("test.json", tmplt, data),
	)
	assert.NoError(t, err)

	expected := `{
		"foo": "bar",
		"baz": "bat",
}`

	testJsonContents, err := os.ReadFile(path.Join(tmpdir, "test.json"))
	assert.NoError(t, err)
	assert.Equal(t, expected, string(testJsonContents))

	data2 := struct{ Foo string }{"Bar"}

	err = g.Generate(
		FileFromTemplate("test2.json", tmplt, data2),
	)
	assert.Error(t, err)
}

func TestFSGenerator_Generate_NoCreateOutputDir(t *testing.T) {
	tmpdir := t.TempDir()

	g := &FSGenerator{
		OutputDir:         path.Join(tmpdir, "will_not_create"),
		NoCreateOutputDir: true,
	}

	err := g.Generate(PlainFile("test.txt", "contents"))
	assert.Error(t, err)

	err = os.Mkdir(g.OutputDir, 0755)
	assert.NoError(t, err)

	err = g.Generate(PlainFile("test.txt", "contents"))
	assert.NoError(t, err)

	actual, err := os.ReadFile(path.Join(tmpdir, "will_not_create", "test.txt"))
	assert.NoError(t, err)
	assert.Equal(t, "contents", string(actual))
}
