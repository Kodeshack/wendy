package wendy

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFSGenerator_Generate(t *testing.T) {
	tmpdir := t.TempDir()

	g := &FSGenerator{RootDir: tmpdir}

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

	g := &FSGenerator{RootDir: tmpdir}

	err := g.Generate(Dir("will_exist"))
	assert.NoError(t, err)

	err = g.Generate(Dir("will_exist"))
	assert.NoError(t, err)

	g = &FSGenerator{RootDir: tmpdir, ErrorOnExistingDir: true}

	err = g.Generate(Dir("will_exist"))
	assert.Error(t, err)
}

func TestFSGenerator_Generate_CleanDir(t *testing.T) {
	tmpdir := t.TempDir()

	g := &FSGenerator{
		RootDir:            tmpdir,
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
