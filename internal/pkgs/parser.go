package pkgs

import (
	"fmt"
	"nuconv/internal/gen"
	"os"
	"path/filepath"

	"github.com/dave/jennifer/jen"
	"golang.org/x/tools/go/packages"
)

type ParsedPackage struct {
	Dir string
	Gen gen.Package
}

type Parser struct {
	LogicalPkgs map[gen.PkgPath]ParsedPackage
}

func NewParser() *Parser {
	return &Parser{
		LogicalPkgs: map[gen.PkgPath]ParsedPackage{},
	}
}

func (p *Parser) Load(patterns ...string) (err error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo,
		Dir: ".",
	}, patterns...)
	if err != nil {
		err = fmt.Errorf("load packages (%v): %w", patterns, err)
		return
	}

	for _, pkg := range pkgs {
		err = p.parsePkg(pkg)
		if err != nil {
			err = fmt.Errorf("parse pkg (%s): %w", pkg.PkgPath, err)
			return
		}
	}

	return
}

func (p *Parser) Generate() (err error) {
	for path, pkg := range p.LogicalPkgs {
		err = p.genPkg(pkg)
		if err != nil {
			err = fmt.Errorf("gen pkg (%s): %w", path, err)
			return
		}
	}
	return
}

func (p *Parser) parsePkg(pkg *packages.Package) (err error) {
	parser, err := newPkgParser(pkg)
	if err != nil {
		err = fmt.Errorf("new pkg parser (%s): %w", pkg.PkgPath, err)
		return
	}
	generated, err := parser.Parse()
	if err != nil {
		err = fmt.Errorf("parse pkg (%s): %w", pkg.PkgPath, err)
		return
	}

	p.LogicalPkgs[generated.Path] = ParsedPackage{
		Dir: pkg.Dir,
		Gen: generated,
	}

	return
}

func (p *Parser) genPkg(pkg ParsedPackage) (err error) {
	g := gen.GenContext{
		Out: jen.NewFilePath(string(pkg.Gen.Path)),
	}

	gen.RenderHelpers(g)

	for _, t := range pkg.Gen.Types {
		_, err = t.Type.Definition(g)
		if err != nil {
			err = fmt.Errorf("render type definition (%s): %w", t.Type.TypeID().Name, err)
			return
		}
		_, err = t.Type.Parser(g)
		if err != nil {
			err = fmt.Errorf("render parser (%s): %w", t.Type.TypeID().Name, err)
			return
		}
		_, err = t.Type.Serializer(g)
		if err != nil {
			err = fmt.Errorf("render serializer (%s): %w", t.Type.TypeID().Name, err)
			return
		}
	}

	targetPath := filepath.Join(pkg.Dir, "nuconv.gen.go")
	f, err := os.Create(targetPath)
	if err != nil {
		err = fmt.Errorf("create file (%s): %w", targetPath, err)
		return
	}
	defer f.Close()

	err = g.Out.Render(f)
	if err != nil {
		err = fmt.Errorf("render file (%s): %w", targetPath, err)
		return
	}

	return
}
