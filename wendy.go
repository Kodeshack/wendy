package wendy

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

type FSGenerator struct {
	RootDir string
}

func (g *FSGenerator) Generate(files ...File) error {
	rootDir, err := filepath.Abs(g.RootDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		err := g.generate(rootDir, f)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *FSGenerator) generate(parentDir string, file File) error {
	if dir, ok := file.(Directory); ok {
		return g.generateDir(parentDir, dir)
	}

	if wt, ok := file.(io.WriterTo); ok {
		return g.generateRealFile(path.Join(parentDir, file.Name()), wt)
	}

	return nil
}

func (g *FSGenerator) generateDir(parentDir string, dir Directory) error {
	dirpath := path.Join(parentDir, dir.Name())

	err := os.Mkdir(dirpath, 0755)
	if err != nil {
		return err
	}

	entries, err := dir.Entries()
	if err != nil {
		return err
	}

	for _, f := range entries {
		err = g.generate(dirpath, f)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *FSGenerator) generateRealFile(filepath string, wt io.WriterTo) (err error) {
	fh, err := os.OpenFile(filepath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() { err = fh.Close() }()

	_, err = wt.WriteTo(fh)

	return
}

type File interface {
	Name() string
}

type Directory interface {
	File
	Entries() ([]File, error)
}

func Generate(f File) error {
	return nil
}

func Dir(name string, entries ...File) Directory {
	return &dir{name: name, entries: entries}
}

type dir struct {
	name    string
	entries []File
}

func (d *dir) Name() string {
	return d.name
}

func (d *dir) Entries() ([]File, error) {
	return d.entries, nil
}

func PlainFile(name string, contents string) File {
	return &plainFile{name: name, contents: []byte(contents)}
}

type plainFile struct {
	name     string
	contents []byte
}

func (f *plainFile) Name() string {
	return f.name
}

// WriteTo implements [io.WriterTo]
func (f *plainFile) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(f.contents)
	if err != nil {
		return 0, err
	}

	return int64(n), nil
}

func TemplateFile(name string, template string, data any) File {
	return &tmplFile{name: name, template: template, data: data}
}

type tmplFile struct {
	name     string
	template string
	data     any
}

func (f *tmplFile) Name() string {
	return f.name
}

// WriteTo implements [io.WriterTo]
func (f *tmplFile) WriteTo(w io.Writer) (int64, error) {
	t, err := template.New(f.name).Parse(f.template)
	if err != nil {
		return 0, err
	}

	return 0, t.Execute(w, f.data)
}

type Template interface {
	Execute(w io.Writer, data any) error
}

func FileFromTemplate(name string, template Template, data any) File {
	return &fileFromTmpl{name: name, template: template, data: data}
}

type fileFromTmpl struct {
	name     string
	template Template
	data     any
}

func (f *fileFromTmpl) Name() string {
	return f.name
}

// WriteTo implements [io.WriterTo]
func (f *fileFromTmpl) WriteTo(w io.Writer) (int64, error) {
	return 0, f.template.Execute(w, f.data)
}
