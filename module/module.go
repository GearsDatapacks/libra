package module

import (
	"os"
	"path"
	"strings"

	"github.com/gearsdatapacks/libra/diagnostics"
	"github.com/gearsdatapacks/libra/lexer"
	"github.com/gearsdatapacks/libra/parser"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/text"
)

type File struct {
	Path string
	Ast  *ast.Program
}

func loadFile(path string) (*File, diagnostics.Manager) {
	file := &File{Path: path}
	l := lexer.New(text.LoadFile(path))
	tokens := l.Tokenise()
	if len(l.Diagnostics) != 0 {
		return file, l.Diagnostics
	}

	p := parser.New(tokens, l.Diagnostics)
	file.Ast = p.Parse()
	return file, p.Diagnostics
}

func loadModule(modPath string) ([]File, diagnostics.Manager) {
	dir, err := os.ReadDir(modPath)
	if err != nil {
		panic(err)
	}

	files := []File{}
	diagnostics := diagnostics.Manager{}
	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".lb") {
			file, diags := loadFile(path.Join(modPath, entry.Name()))
			files = append(files, *file)
			diagnostics = append(diagnostics, diags...)
		}
	}

	return files, diagnostics
}

type Module struct {
	Name     string
	Files    []File
	Imported map[string]*Module
}

var fetchedModules = map[string]*Module{}

func Load(filePath string) (*Module, diagnostics.Manager) {
	modPath := filePath
	if !isDir(modPath) {
		modPath = path.Dir(modPath)
	}
	if fetched, ok := fetchedModules[modPath]; ok {
		return fetched, diagnostics.Manager{}
	}

	files, diagnostics := loadModule(modPath)
	_, name := path.Split(modPath)
	mod := &Module{
		Name:     name,
		Files:    files,
		Imported: map[string]*Module{},
	}
	fetchedModules[modPath] = mod

	for _, file := range files {
		for _, stmt := range file.Ast.Statements {
			if importStmt, ok := stmt.(*ast.ImportStatement); ok {
				importedPath := path.Join(modPath, importStmt.Module.Value)
				imported, diags := Load(importedPath)
				diagnostics = append(diagnostics, diags...)
				mod.Imported[importStmt.Module.Value] = imported
			}
		}
	}

	return mod, diagnostics
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return info.IsDir()
}
