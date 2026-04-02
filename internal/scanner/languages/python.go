package languages

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ossiqn/ghost/internal/scanner"
)

type PythonLanguage struct{}

func (l *PythonLanguage) Name() string { return "python" }
func (l *PythonLanguage) Extensions() []string { return []string{".py"} }

var (
	pyFuncDecl   = regexp.MustCompile(`^def\s+(\w+)\s*\(`)
	pyClassDecl  = regexp.MustCompile(`^class\s+(\w+)\s*[:(]`)
	pyImport     = regexp.MustCompile(`^import\s+(\w+)|^from\s+(\w+)\s+import\s+(.+)`)
)

func (l *PythonLanguage) Scan(path string) *scanner.Result {
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

	usageRe := regexp.MustCompile(`\b(\w+)\b`)
	for _, match := range usageRe.FindAllString(content, -1) {
		used[match] = true
	}

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if matches := pyFuncDecl.FindStringSubmatch(trimmed); matches != nil {
			defined[matches[1]] = i + 1
		}

		if matches := pyClassDecl.FindStringSubmatch(trimmed); matches != nil {
			defined[matches[1]] = i + 1
		}

		if matches := pyImport.FindStringSubmatch(trimmed); matches != nil {
			if matches[1] != "" {
				imports[matches[1]] = i + 1
			}
			if matches[3] != "" {
				for _, imp := range strings.Split(matches[3], ",") {
					name := strings.TrimSpace(imp)
					name = strings.Split(name, " as ")[0]
					name = strings.TrimSpace(name)
					if name != "" && name != "*" {
						imports[name] = i + 1
					}
				}
			}
		}
	}

	builtins := map[string]bool{
		"print": true, "len": true, "range": true, "str": true,
		"int": true, "list": true, "dict": true, "set": true,
		"True": true, "False": true, "None": true, "self": true,
	}

	for name, line := range defined {
		if name == "__init__" || name == "__main__" || builtins[name] {
			continue
		}
		count := strings.Count(content, name)
		if count <= 1 {
			result.DeadFunctions = append(result.DeadFunctions, scanner.DeadCode{
				File:     path,
				Line:     line,
				Name:     name,
				Type:     "function",
				Language: "python",
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
				Language: "python",
			})
		}
	}

	return result
}

func (l *PythonLanguage) ScanDeps(root string) ([]scanner.GhostDep, error) {
	reqPath := filepath.Join(root, "requirements.txt")
	f, err := os.Open(reqPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var declared []scanner.GhostDep
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := regexp.MustCompile(`[>=<!~^]+`).Split(line, 2)
		name := strings.TrimSpace(parts[0])
		version := ""
		if len(parts) > 1 {
			version = strings.TrimSpace(parts[1])
		}

		declared = append(declared, scanner.GhostDep{
			Name:    name,
			Version: version,
			File:    reqPath,
		})
	}

	usedImports := make(map[string]bool)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".py") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		re := regexp.MustCompile(`^(?:import|from)\s+(\w+)`)
		for _, match := range re.FindAllSubmatch(content, -1) {
			if len(match) > 1 {
				usedImports[string(match[1])] = true
			}
		}

		return nil
	})

	var ghosts []scanner.GhostDep
	for _, dep := range declared {
		normalizedName := strings.ToLower(strings.ReplaceAll(dep.Name, "-", "_"))
		if !usedImports[normalizedName] && !usedImports[dep.Name] {
			ghosts = append(ghosts, dep)
		}
	}

	return ghosts, nil
}
