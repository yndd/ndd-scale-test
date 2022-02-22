package generator

import (
	"bytes"
	"errors"
	"io"
	"os"
	"text/template"

	"gihub.com/yndd/ndd-scale-test/pkg/output"
	"gihub.com/yndd/ndd-scale-test/pkg/templ"
)

type Generator interface {
	WithIndexes(o, c int)
	WithTemplate(t string)

	Generate() ([]*bytes.Buffer, error)
}

func WithTemplate(t string) Option {
	return func(g Generator) {
		g.WithTemplate(t)
	}
}

func WithIndexes(o, c int) Option {
	return func(g Generator) {
		g.WithIndexes(o, c)
	}
}

// Option can be used to manipulate Options.
type Option func(g Generator)

func NewGenerator(opts ...Option) (Generator, error) {
	g := &generator{}
	for _, opt := range opts {
		opt(g)
	}
	if err := g.initTemplates(); err != nil {
		return nil, err
	}
	return g, nil
}

type generator struct {
	offset       int
	count        int
	templateFile string
	outputDir    string
	template     *template.Template
	outputPlugin output.OutputPlugin
}

func (g *generator) WithTemplate(t string) {
	g.templateFile = t
}

func (g *generator) WithOutputDir(dir string) {
	g.outputDir = dir
}

func (g *generator) WithIndexes(o, c int) {
	g.offset = o
	g.count = c
}

func (g *generator) WithOutputPlugin(to output.OutputPlugin) {
	g.outputPlugin = to
}

func (g *generator) initTemplates() error {
	if !fileExists(g.templateFile) {
		return errors.New("template file does not exist")
	}
	var err error
	g.template, err = templ.ParseTemplate(g.templateFile)
	if err != nil {
		return err
	}
	return nil
}

// FileExists function
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (g *generator) Generate() ([]*bytes.Buffer, error) {
	result := []*bytes.Buffer{}

	for i := g.offset; i <= g.offset+g.count; i++ {

		b := bytes.NewBuffer([]byte{})

		if err := g.renderTemplate(b, i); err != nil {
			return nil, err
		}
		result = append(result, b)

	}
	return result, nil
}

func (g *generator) renderTemplate(f io.Writer, i int) error {
	s := struct {
		Index int
	}{
		Index: i,
	}
	if err := g.template.Execute(f, s); err != nil {
		return err
	}
	return nil
}
