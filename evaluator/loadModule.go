package evaluator

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	. "github.com/BEN00262/simpleLang/exceptions"
	. "github.com/BEN00262/simpleLang/lexer"
	. "github.com/BEN00262/simpleLang/parser"
	. "github.com/BEN00262/simpleLang/symbolstable"
)

func (eval *Evaluator) _evaluateProgramNode(nodes []interface{}) ExceptionNode {
	for _, node := range nodes {
		_, exception := eval.walkTree(node)

		if exception.Type != NO_EXCEPTION {
			return exception
		}
	}

	return ExceptionNode{Type: NO_EXCEPTION}
}

// create a dependancy graph
type ImportModule struct {
	context ContextValue
}

// how tf do we work with imports to the same file across board
func (eval *Evaluator) LoadModule(module Import) ExceptionNode {
	// ensure the filename exists --> also check for errors in the lexer and the parser too
	// have a * we dump to the global scope
	// otherwise we namespace

	// create push our own context then use it later
	isSystemImport := true

	// this is a local import
	if strings.HasPrefix(module.FileName, "./") || strings.HasPrefix(module.FileName, "../") {
		isSystemImport = false
	}

	// check if the module has a .happ extension if not add it
	if filepath.Ext(module.FileName) != ".happ" {
		if module.Alias == "" {
			module.Alias = module.FileName
		}

		module.FileName += ".happ"
	}

	importPath := ""

	if isSystemImport {
		_basePath, err := os.Executable()

		if err != nil {
			// we get the error code for not working here
			// what happens for now we use the path in the current directory

			return ExceptionNode{
				Type:    MODULE_IMPORT_EXCEPTION,
				Message: err.Error(),
			}
		}

		importPath = filepath.Join(
			filepath.Dir(_basePath),
			"includes",
			module.FileName,
		)
	} else {
		// local import
		importPath = filepath.Join(
			eval.baseFilePath,
			module.FileName,
		)
	}

	if importPath == "" {
		return ExceptionNode{
			Type:    MODULE_IMPORT_EXCEPTION,
			Message: "Failed to resolve the import path",
		}
	}

	importedModule, err := ioutil.ReadFile(importPath)

	if err != nil {
		return ExceptionNode{
			Type:    MODULE_IMPORT_EXCEPTION,
			Message: err.Error(),
		}
	}

	lexer := InitLexer(string(importedModule))
	parser := InitParser(lexer.Lex(), lexer.SplitCode)

	if module.Alias != "*" {
		eval.symbolsTable.PushContext()

		_exception := eval._evaluateProgramNode(parser.Parse().Nodes)

		if _exception.Type != NO_EXCEPTION {
			return _exception
		}

		_module_context := eval.symbolsTable.GetTopContext()

		eval.symbolsTable.PushToContext(module.Alias, SymbolTableValue{
			Type: IMPORTED_MODULE,
			Value: ImportModule{
				context: _module_context,
			},
		})

		return _exception
	}

	return eval._evaluateProgramNode(parser.Parse().Nodes)
}

func (eval *Evaluator) _eval(codeString string) (result interface{}, exception ExceptionNode) {
	var lexer *Lexer
	var program *ProgramNode

	if eval.evalCache.IsFresh(codeString) {
		lexer = eval.evalCache.Lexed
		program = eval.evalCache.Program
	} else {
		lexer = InitLexer(codeString)
		program = InitParser(lexer.Lex(), lexer.SplitCode).Parse()

		go (func() {
			eval.evalCache.UpdateCache(codeString, lexer, program)
		})()
	}

	for _, node := range program.Nodes {

		switch node.(type) {
		case ExpressionNode:
			{
				result, exception = eval.walkTree(node)

				if exception.Type != NO_EXCEPTION {
					return nil, exception
				}
			}
		default:
			return nil, ExceptionNode{
				Type:    INVALID_OPERATION_EXCEPTION,
				Message: "Expected an expression",
			}
		}
		break
	}
	return
}
