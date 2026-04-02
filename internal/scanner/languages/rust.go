package languages

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ossiqn/ghost/internal/scanner"
)

type RustLanguage struct{}

var (
	rustFnDecl  = regexp.MustCompile(`^\s*(?:pub\s+)?(?:async\s+)?fn\s+(\w+)`)
	rustUseDecl = regexp.MustCompile(`^\s*use\s+([^;]+);`)
)

func (r *RustLanguage) Name() string { return "rust" }
func (r *RustLanguage) Extensions() []string { return []string{".rs"} }

func (r *RustLanguage) Scan(path string) *scanner.Result {
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

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if matches := rustFnDecl.FindStringSubmatch(trimmed); matches != nil {
			defined[matches[1]] = i + 1
		}

		if matches := rustUseDecl.FindStringSubmatch(trimmed); matches != nil {
			lastPart := matches[1]
			parts := strings.Split(lastPart, "::")
			name := strings.TrimSpace(parts[len(parts)-1])
			name = strings.Trim(name, "{}")
			imports[name] = i + 1
		}
	}

	skip := map[string]bool{
		"main": true,
		"new":  true,
	}

	for name, line := range defined {
		if skip[name] {
			continue
		}
		if strings.HasPrefix(name, "test_") {
			continue
		}
		count := strings.Count(content, name)
		if count <= 1 {
			result.DeadFunctions = append(result.DeadFunctions, scanner.DeadCode{
				File:     path,
				Line:     line,
				Name:     name,
				Type:     "function",
				Language: "rust",
			})
		}
	}

	for name, line := range imports {
		if !strings.Contains(content, name) {
			result.UnusedImports = append(result.UnusedImports, scanner.DeadCode{
				File:     path,
				Line:     line,
				Name:     name,
				Type:     "use",
				Language: "rust",
			})
		}
	}

	return result
}

func (r *RustLanguage) ScanDeps(root string) ([]scanner.GhostDep, error) {
	cargoPath := filepath.Join(root, "Cargo.toml")
	f, err := os.Open(cargoPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var deps []scanner.GhostDep
	sc := bufio.NewScanner(f)
	inDeps := false

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())

		if line == "[dependencies]" || line == "[dev-dependencies]" {
			inDeps = true
			continue
		}
		if strings.HasPrefix(line, "[") {
			inDeps = false
			continue
		}

		if inDeps && strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			name := strings.TrimSpace(parts[0])
			version := strings.Trim(strings.TrimSpace(parts[1]), `"`)
			deps = append(deps, scanner.GhostDep{
				Name:    name,
				Version: version,
				File:    cargoPath,
			})
		}
	}

	return deps, nil
}
