package languages

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/ossiqn/ghost/internal/scanner"
)

type JavaLanguage struct{}

var (
	javaMethodDecl  = regexp.MustCompile(`^\s*(?:public|private|protected|static|\s)+[\w<>\[\]]+\s+(\w+)\s*\(`)
	javaImportDecl  = regexp.MustCompile(`^\s*import\s+(?:static\s+)?([^;]+);`)
)

func (j *JavaLanguage) Name() string { return "java" }
func (j *JavaLanguage) Extensions() []string { return []string{".java"} }

func (j *JavaLanguage) Scan(path string) *scanner.Result {
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

		if matches := javaMethodDecl.FindStringSubmatch(trimmed); matches != nil {
			defined[matches[1]] = i + 1
		}

		if matches := javaImportDecl.FindStringSubmatch(trimmed); matches != nil {
			parts := strings.Split(strings.TrimSpace(matches[1]), ".")
			name := parts[len(parts)-1]
			if name != "*" {
				imports[name] = i + 1
			}
		}
	}

	skip := map[string]bool{
		"main":        true,
		"toString":    true,
		"hashCode":    true,
		"equals":      true,
		"getInstance": true,
	}

	for name, line := range defined {
		if skip[name] {
			continue
		}
		count := strings.Count(content, name)
		if count <= 1 {
			result.DeadFunctions = append(result.DeadFunctions, scanner.DeadCode{
				File:     path,
				Line:     line,
				Name:     name,
				Type:     "method",
				Language: "java",
			})
		}
	}

	for name, line := range imports {
		if strings.Count(content, name) <= 1 {
			result.UnusedImports = append(result.UnusedImports, scanner.DeadCode{
				File:     path,
				Line:     line,
				Name:     name,
				Type:     "import",
				Language: "java",
			})
		}
	}

	return result
}

func (j *JavaLanguage) ScanDeps(root string) ([]scanner.GhostDep, error) {
	// Maven pom.xml scan edilebilir ileride
	return nil, nil
}
