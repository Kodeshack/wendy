package wendy

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

type FSGenerator struct {
	OutputDir           string
	ErrorOnExistingDir  bool
	CleanDir            bool
	NoCreateOutputDir   bool
	ErrorOnExistingFile bool
}

type genfile struct {
	path     string
	contents WriterToFile
}

type gendir struct {
	path string
}

func (g *FSGenerator) Generate(files ...File) error {
	return g.GenerateCtx(context.Background(), files...)
}

func (g *FSGenerator) GenerateCtx(ctx context.Context, files ...File) error {
	if g.CleanDir {
		err := cleanDir(g.OutputDir)
		if err != nil {
			return err
		}
	}

	rootDir, err := filepath.Abs(g.OutputDir)
	if err != nil {
		return err
	}

	if !g.NoCreateOutputDir {
		err = os.Mkdir(rootDir, 0755)
		if err != nil {
			if !errors.Is(err, os.ErrExist) {
				return err
			}
		}
	}

	gendirs := make([]*gendir, 0, len(files))
	genfiles := make([]*genfile, 0, len(files))

	for _, f := range files {
		dirs, files, err := g.generate(ctx, rootDir, f)
		if err != nil {
			return err
		}

		gendirs = append(gendirs, dirs...)
		genfiles = append(genfiles, files...)
	}

	for _, d := range gendirs {
		err = g.generateRealDir(d.path)
		if err != nil {
			return err
		}
	}

	for _, f := range genfiles {
		err = g.generateRealFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *FSGenerator) generateRealDir(dir string) error {
	err := os.Mkdir(dir, 0755)
	if err != nil {
		if !errors.Is(err, os.ErrExist) || g.ErrorOnExistingDir {
			return err
		}
	}

	return nil
}

func (g *FSGenerator) generateRealFile(file *genfile) error {
	if g.ErrorOnExistingFile {
		stat, statErr := os.Stat(file.path)
		if statErr != nil {
			if !errors.Is(statErr, os.ErrNotExist) {
				return statErr
			}
		}

		if stat != nil {
			return fmt.Errorf("file already exits %s: %w", file.path, os.ErrExist)
		}
	}

	tmpfile, err := os.CreateTemp("", path.Base(file.path))
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(
			err,
			tmpfile.Close(),
		)

		if err != nil {
			err = errors.Join(err, os.Remove(tmpfile.Name()))
		}
	}()

	_, err = file.contents.WriteToFile(file.path, tmpfile)
	if err != nil {
		return err
	}

	err = os.Rename(tmpfile.Name(), file.path)

	return err
}

func (g *FSGenerator) generate(ctx context.Context, parentDir string, file File) ([]*gendir, []*genfile, error) {
	select {
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	default:
	}

	if dir, ok := file.(Directory); ok {
		return g.generateDir(ctx, parentDir, dir)
	}

	if wt, ok := file.(WriterToFile); ok {
		return nil, []*genfile{{path: path.Join(parentDir, file.Name()), contents: wt}}, nil
	}

	if wt, ok := file.(io.WriterTo); ok {
		return nil, []*genfile{{path: path.Join(parentDir, file.Name()), contents: &writerToAdapter{wt}}}, nil
	}

	return nil, nil, nil
}

func (g *FSGenerator) generateDir(ctx context.Context, parentDir string, dir Directory) ([]*gendir, []*genfile, error) {
	select {
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	default:
	}

	dirpath := path.Join(parentDir, dir.Name())

	entries, err := dir.Entries()
	if err != nil {
		return nil, nil, err
	}

	gendirs := []*gendir{{path: dirpath}}
	genfiles := make([]*genfile, 0, len(entries))

	for _, f := range entries {
		dirs, files, err := g.generate(ctx, dirpath, f)
		if err != nil {
			return nil, nil, err
		}

		gendirs = append(gendirs, dirs...)
		genfiles = append(genfiles, files...)
	}

	return gendirs, genfiles, nil
}

type File interface {
	Name() string
}

type Directory interface {
	File
	Entries() ([]File, error)
}

type WriterToFile interface {
	WriteToFile(filename string, w io.Writer) (n int64, err error)
}

type writerToAdapter struct {
	io.WriterTo
}

func (a *writerToAdapter) WriteToFile(_ string, w io.Writer) (n int64, err error) {
	return a.WriteTo(w)
}

func cleanDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	for _, e := range entries {
		err := os.RemoveAll(path.Join(dir, e.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}
