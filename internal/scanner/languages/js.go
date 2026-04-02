package languages

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ossiqn/ghost/internal/scanner"
)

type JavaScriptLanguage struct{}

func (l *JavaScriptLanguage) Name() string { return "javascript" }
func (l *JavaScriptLanguage) Extensions() []string { return []string{".js", ".ts", ".jsx", ".tsx"} }

var (
	jsFuncDecl    = regexp.MustCompile(`(?:^|\s)(?:function\s+(\w+)|const\s+(\w+)\s*=\s*(?:async\s*)?\(|let\s+(\w+)\s*=\s*(?:async\s*)?\()`)
	jsArrowFunc   = regexp.MustCompile(`(?:export\s+)?(?:const|let|var)\s+(\w+)\s*=\s*(?:async\s*)?\(.*?\)\s*=>`)
	jsImportDecl  = regexp.MustCompile(`import\s+(?:{([^}]+)}|(\w+)|\*\s+as\s+(\w+))\s+from\s+['"]([^'"]+)['"]`)
	jsUsage       = regexp.MustCompile(`\b(\w+)\b`)
)

func (l *JavaScriptLanguage) Scan(path string) *scanner.Result {
	result := &scanner.Result{}

	f, err := os.Open(path)
	if err != nil {
		return result
	}
	defer f.Close()

	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	content := strings.Join(lines, "\n")

	defined := make(map[string]int)
	imports := make(map[string]int)
	used := make(map[string]bool)

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if matches := jsFuncDecl.FindStringSubmatch(trimmed); matches != nil {
			for _, m := range matches[1:] {
				if m != "" {
					defined[m] = i + 1
				}
			}
		}

		if matches := jsArrowFunc.FindStringSubmatch(trimmed); matches != nil {
			if matches[1] != "" {
				defined[matches[1]] = i + 1
			}
		}

		if matches := jsImportDecl.FindStringSubmatch(trimmed); matches != nil {
			if matches[1] != "" {
				for _, imp := range strings.Split(matches[1], ",") {
					name := strings.TrimSpace(imp)
					name = strings.Split(name, " as ")[0]
					name = strings.TrimSpace(name)
					if name != "" {
						imports[name] = i + 1
					}
				}
			}
			for _, m := range matches[2:4] {
				if m != "" {
					imports[m] = i + 1
				}
			}
		}
	}

	for _, match := range jsUsage.FindAllString(content, -1) {
		used[match] = true
	}

	for name, line := range defined {
		if !used[name] && !strings.HasPrefix(name, "_") {
			result.DeadFunctions = append(result.DeadFunctions, scanner.DeadCode{
				File:     path,
				Line:     line,
				Name:     name,
				Type:     "function",
				Language: "javascript",
			})
		}
	}

	for name, line := range imports {
		if !used[name] {
			result.UnusedImports = append(result.UnusedImports, scanner.DeadCode{
				File:     path,
				Line:     line,
				Name:     name,
				Type:     "import",
				Language: "javascript",
			})
		}
	}

	return result
}

func (l *JavaScriptLanguage) ScanDeps(root string) ([]scanner.GhostDep, error) {
	pkgPath := filepath.Join(root, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return nil, err
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}

	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	usedImports := make(map[string]bool)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".js" && ext != ".ts" && ext != ".jsx" && ext != ".tsx" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		re := regexp.MustCompile(`(?:import|require)\s*\(?['"]([^'"./][^'"]*)['"]\)?`)
		for _, match := range re.FindAllSubmatch(content, -1) {
			if len(match) > 1 {
				pkg := strings.Split(string(match[1]), "/")[0]
				usedImports[pkg] = true
			}
		}

		return nil
	})

	var ghosts []scanner.GhostDep
	allDeps := make(map[string]string)
	for k, v := range pkg.Dependencies {
		allDeps[k] = v
	}
	for k, v := range pkg.DevDependencies {
		allDeps[k] = v
	}

	for dep, version := range allDeps {
		if !usedImports[dep] {
			ghosts = append(ghosts, scanner.GhostDep{
				Name:    dep,
				Version: version,
				File:    pkgPath,
			})
		}
	}

	return ghosts, nil
}
