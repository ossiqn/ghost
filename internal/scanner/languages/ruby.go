package languages

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/ossiqn/ghost/internal/scanner"
)

type RubyLanguage struct{}

var (
	rubyMethodDecl = regexp.MustCompile(`^\s*def\s+(\w+)`)
	rubyClassDecl  = regexp.MustCompile(`^\s*class\s+(\w+)`)
	rubyRequire    = regexp.MustCompile(`^\s*require\s+['"]([^'"]+)['"]`)
	rubyRequireRel = regexp.MustCompile(`^\s*require_relative\s+['"]([^'"]+)['"]`)
)

func (r *RubyLanguage) Name() string { return "ruby" }
func (r *RubyLanguage) Extensions() []string { return []string{".rb", ".rake"} }

func (r *RubyLanguage) Scan(path string) *scanner.Result {
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

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if matches := rubyMethodDecl.FindStringSubmatch(trimmed); matches != nil {
			defined[matches[1]] = i + 1
		}
		if matches := rubyClassDecl.FindStringSubmatch(trimmed); matches != nil {
			defined[matches[1]] = i + 1
		}
	}

	skip := map[string]bool{
		"initialize": true,
		"to_s":       true,
		"to_i":       true,
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
				Language: "ruby",
			})
		}
	}

	return result
}

func (r *RubyLanguage) ScanDeps(root string) ([]scanner.GhostDep, error) {
	// Ruby Gemfile scan edilebilir ileride
	return nil, nil
}
