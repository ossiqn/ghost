package languages

import (
	"bufio"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/ossiqn/ghost/internal/scanner"
)

type GoLanguage struct{}

func (l *GoLanguage) Name() string { return "go" }
func (l *GoLanguage) Extensions() []string { return []string{".go"} }

func (l *GoLanguage) Scan(path string) *scanner.Result {
	result := &scanner.Result{}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		return result
	}

	defined := make(map[string]int)
	called := make(map[string]bool)

	ast.Inspect(f, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.FuncDecl:
			if node.Name != nil {
				pos := fset.Position(node.Pos())
				defined[node.Name.Name] = pos.Line
			}
		case *ast.CallExpr:
			switch fn := node.Fun.(type) {
			case *ast.Ident:
				called[fn.Name] = true
			case *ast.SelectorExpr:
				called[fn.Sel.Name] = true
			}
		}
		return true
	})

	for name, line := range defined {
		if name == "main" || name == "init" {
			continue
		}
		if isGoExported(name) {
			continue
		}
		if !called[name] {
			result.DeadFunctions = append(result.DeadFunctions, scanner.DeadCode{
				File:     path,
				Line:     line,
				Name:     name,
				Type:     "function",
				Language: "go",
			})
		}
	}

	result.UnusedImports = append(result.UnusedImports, scanGoUnusedImports(path, f, fset)...)
	result.UnusedVariables = append(result.UnusedVariables, scanGoUnusedVars(path, f, fset)...)

	return result
}

func scanGoUnusedImports(path string, f *ast.File, fset *token.FileSet) []scanner.DeadCode {
	var unused []scanner.DeadCode

	usedPkgs := make(map[string]bool)
	ast.Inspect(f, func(n ast.Node) bool {
		if sel, ok := n.(*ast.SelectorExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok {
				usedPkgs[ident.Name] = true
			}
		}
		return true
	})

	for _, imp := range f.Imports {
		if imp.Name != nil && imp.Name.Name == "_" {
			continue
		}

		var pkgName string
		if imp.Name != nil {
			pkgName = imp.Name.Name
		} else {
			parts := strings.Split(strings.Trim(imp.Path.Value, `"`), "/")
			pkgName = parts[len(parts)-1]
		}

		if !usedPkgs[pkgName] {
			pos := fset.Position(imp.Pos())
			unused = append(unused, scanner.DeadCode{
				File:     path,
				Line:     pos.Line,
				Name:     strings.Trim(imp.Path.Value, `"`),
				Type:     "import",
				Language: "go",
			})
		}
	}

	return unused
}

func scanGoUnusedVars(path string, f *ast.File, fset *token.FileSet) []scanner.DeadCode {
	var unused []scanner.DeadCode

	ast.Inspect(f, func(n ast.Node) bool {
		vs, ok := n.(*ast.ValueSpec)
		if !ok {
			return true
		}

		for _, name := range vs.Names {
			if name.Name == "_" {
				continue
			}
			if isGoExported(name.Name) {
				continue
			}
			if name.Obj != nil && !isGoUsed(f, name.Name) {
				pos := fset.Position(name.Pos())
				unused = append(unused, scanner.DeadCode{
					File:     path,
					Line:     pos.Line,
					Name:     name.Name,
					Type:     "variable",
					Language: "go",
				})
			}
		}
		return true
	})

	return unused
}

func (l *GoLanguage) ScanDeps(root string) ([]scanner.GhostDep, error) {
	modPath := filepath.Join(root, "go.mod")
	f, err := os.Open(modPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var declared []scanner.GhostDep
	sc := bufio.NewScanner(f)
	inRequire := false

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())

		if line == "require (" {
			inRequire = true
			continue
		}
		if line == ")" {
			inRequire = false
			continue
		}

		if inRequire && line != "" {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				declared = append(declared, scanner.GhostDep{
					Name:    parts[0],
					Version: parts[1],
					File:    modPath,
				})
			}
		}
	}

	usedImports := make(map[string]bool)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return nil
		}

		for _, imp := range f.Imports {
			importPath := strings.Trim(imp.Path.Value, `"`)
			usedImports[importPath] = true
		}

		return nil
	})

	var ghosts []scanner.GhostDep
	for _, dep := range declared {
		isUsed := false
		for usedImport := range usedImports {
			if strings.HasPrefix(usedImport, dep.Name) {
				isUsed = true
				break
			}
		}
		if !isUsed {
			ghosts = append(ghosts, dep)
		}
	}

	return ghosts, nil
}

func isGoExported(name string) bool {
	if len(name) == 0 {
		return false
	}
	return name[0] >= 'A' && name[0] <= 'Z'
}

func isGoUsed(f *ast.File, name string) bool {
	used := false
	ast.Inspect(f, func(n ast.Node) bool {
		if ident, ok := n.(*ast.Ident); ok {
			if ident.Name == name && ident.Obj == nil {
				used = true
			}
		}
		return !used
	})
	return used
}
