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
		err := os.Mkdir(rootDir, 0755)
		if err != nil {
			if !errors.Is(err, os.ErrExist) {
				return err
			}
		}
	}

	for _, f := range files {
		err := g.generate(ctx, rootDir, f)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *FSGenerator) generate(ctx context.Context, parentDir string, file File) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if dir, ok := file.(Directory); ok {
		return g.generateDir(ctx, parentDir, dir)
	}

	if wt, ok := file.(WriterToFile); ok {
		return g.generateRealFile(ctx, path.Join(parentDir, file.Name()), wt)
	}

	if wt, ok := file.(io.WriterTo); ok {
		return g.generateRealFile(ctx, path.Join(parentDir, file.Name()), &writerToAdapter{wt})
	}

	return nil
}

func (g *FSGenerator) generateDir(ctx context.Context, parentDir string, dir Directory) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	dirpath := path.Join(parentDir, dir.Name())

	err := os.Mkdir(dirpath, 0755)
	if err != nil {
		if !errors.Is(err, os.ErrExist) || g.ErrorOnExistingDir {
			return err
		}
	}

	entries, err := dir.Entries()
	if err != nil {
		return err
	}

	for _, f := range entries {
		err = g.generate(ctx, dirpath, f)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *FSGenerator) generateRealFile(ctx context.Context, filepath string, wt WriterToFile) (err error) {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if g.ErrorOnExistingFile {
		stat, statErr := os.Stat(filepath)
		if statErr != nil {
			if !errors.Is(statErr, os.ErrNotExist) {
				return statErr
			}
		}

		if stat != nil {
			return fmt.Errorf("file already exits %s: %w", filepath, os.ErrExist)
		}
	}

	fh, err := os.OpenFile(filepath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, fh.Close()) }()

	_, err = wt.WriteToFile(filepath, fh)

	return
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
