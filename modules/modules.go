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

func Get(file string) (*Module, error) {
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

type ModuleManager struct {
	Name           string
	Main           *Module
	SymbolTable    *symbols.SymbolTable
	Env            *environment.Environment
	Modules        map[string]*ModuleManager
	TypeCheckStage int
	InterpretStage int
}

var fetchedModules = map[string]*ModuleManager{}

func NewManager(file string, table *symbols.SymbolTable, env *environment.Environment) (*ModuleManager, error) {
	mod, err := Get(file)
	if err != nil {
		return nil, err
	}
	_, name := path.Split(path.Dir(file))
	m := &ModuleManager{
		Main:        mod,
		SymbolTable: table,
		Env:         env,
		Modules:     map[string]*ModuleManager{},
		Name:        name,
	}
	fetchedModules[file] = m

	basePath, _ := path.Split(file)

	for _, stmt := range m.Main.Ast.Body {
		if importStmt, ok := stmt.(*ast.ImportStatement); ok {
			modPath := path.Clean(basePath + "/" + importStmt.Module + "/main.lb")
			if modManager, loaded := fetchedModules[modPath]; loaded {
				m.Modules[importStmt.Module] = modManager
				continue
			}

			modManager, err := NewManager(modPath, symbols.New(), environment.New())
			if err != nil {
				return nil, err
			}
			m.Modules[importStmt.Module] = modManager
		}
	}

	return m, nil
}

func NewDetatched(table *symbols.SymbolTable, env *environment.Environment) *ModuleManager {
	return &ModuleManager{
		Name:           "main",
		Main:           &Module{
			Path: ".",
			Ast:  ast.Program{},
		},
		SymbolTable:    table,
		Env:            env,
		Modules: map[string]*ModuleManager{},
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
