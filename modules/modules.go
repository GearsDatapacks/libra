package modules

import (
	"os"
	"path"

	"github.com/gearsdatapacks/libra/interpreter/environment"
	"github.com/gearsdatapacks/libra/lexer"
	"github.com/gearsdatapacks/libra/parser"
	"github.com/gearsdatapacks/libra/parser/ast"
	"github.com/gearsdatapacks/libra/type_checker/symbols"
)

type Module struct {
	Path string
	Ast  ast.Program
}

func modFromFile(file string) (*Module, error) {
	code, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	lexer := lexer.New(code)
	tokens, err := lexer.Tokenise()
	if err != nil {
		return nil, err
	}

	parser := parser.New()
	program, err := parser.Parse(tokens)
	if err != nil {
		return nil, err
	}

	return &Module{
		Ast:  program,
		Path: file,
	}, nil
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func Get(file string) ([]Module, error) {
	if !isDir(file) {
		mod, err := modFromFile(file)
		if err != nil {
			return nil, err
		}
		return []Module{*mod}, nil
	}

	dir, err := os.ReadDir(file)
	if err != nil {
		return nil, err
	}
	mods := []Module{}
	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}
		mod, err := modFromFile(path.Join(file, entry.Name()))
		if err != nil {
			return nil, err
		}
		mods = append(mods, *mod)
	}

	return mods, nil
}

type ModuleManager struct {
	Name           string
	Files          []Module
	SymbolTable    *symbols.SymbolTable
	Env            *environment.Environment
	Imported       map[string]*ModuleManager
	TypeCheckStage int
	InterpretStage int
}

var fetchedModules = map[string]*ModuleManager{}

func NewManager(file string, table *symbols.SymbolTable, env *environment.Environment) (*ModuleManager, error) {
	mods, err := Get(file)
	if err != nil {
		return nil, err
	}

	var basePath string
	if isDir(file) {
		basePath = file
	} else {
		basePath = path.Dir(file)
	}

	_, name := path.Split(basePath)
	m := &ModuleManager{
		Files:       mods,
		SymbolTable: table,
		Env:         env,
		Imported:    map[string]*ModuleManager{},
		Name:        name,
	}
	fetchedModules[file] = m

	for _, file := range m.Files {
		for _, stmt := range file.Ast.Body {
			if importStmt, ok := stmt.(*ast.ImportStatement); ok {
				modPath := path.Clean(path.Join(basePath, importStmt.Module))
				if modManager, loaded := fetchedModules[modPath]; loaded {
					m.Imported[importStmt.Module] = modManager
					continue
				}

				modManager, err := NewManager(modPath, symbols.New(), environment.New())
				if err != nil {
					return nil, err
				}
				m.Imported[importStmt.Module] = modManager
			}
		}
	}

	return m, nil
}

func NewDetatched(table *symbols.SymbolTable, env *environment.Environment) *ModuleManager {
	return &ModuleManager{
		Name: "main",
		Files: []Module{{
			Path: ".",
			Ast:  ast.Program{},
		}},
		SymbolTable: table,
		Env:         env,
		Imported:    map[string]*ModuleManager{},
	}
}

func (m *ModuleManager) EnterScope(scope *symbols.SymbolTable) {
	m.SymbolTable = scope
}

func (m *ModuleManager) ExitScope() {
	if m.SymbolTable.Parent == nil {
		panic("Cannot exit global scope")
	}
	m.SymbolTable = m.SymbolTable.Parent
}

func (m *ModuleManager) EnterEnv(env *environment.Environment) {
	m.Env = env
}

func (m *ModuleManager) ExitEnv() {
	if m.Env.Parent == nil {
		panic("Cannot exit global scope")
	}
	m.Env = m.Env.Parent
}
