package scanner

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

type Config struct {
	Path      string
	JSON      bool
	Verbose   bool
	AutoClean bool
	Lang      string
}

type Result struct {
	DeadFunctions   []DeadCode
	UnusedVariables []DeadCode
	UnusedImports   []DeadCode
	GhostDeps       []GhostDep
	TotalFiles      int
	ScannedFiles    int
	EstSavedBytes   int64
}

type DeadCode struct {
	File     string
	Line     int
	Name     string
	Type     string
	Language string
}

type GhostDep struct {
	Name    string
	File    string
	Version string
}

type Scanner struct {
	cfg Config
	reg *Registry
	mu  sync.Mutex
}

func New(cfg Config, reg *Registry) *Scanner {
	return &Scanner{cfg: cfg, reg: reg}
}

func (s *Scanner) Run() (*Result, error) {
	files, err := s.collectFiles()
	if err != nil {
		return nil, err
	}

	result := &Result{
		TotalFiles: len(files),
	}

	if len(files) == 0 {
		return result, nil
	}

	bar := progressbar.NewOptions(len(files),
		progressbar.OptionSetDescription("Scanning"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "█",
			SaucerPadding: "░",
			BarStart:      "",
			BarEnd:        "",
		}),
		progressbar.OptionShowCount(),
		progressbar.OptionClearOnFinish(),
	)

	var wg sync.WaitGroup
	resultChan := make(chan *Result, len(files))

	for _, file := range files {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			defer bar.Add(1)

			r := s.scanFile(f)
			if r != nil {
				resultChan <- r
			}
		}(file)
	}

	wg.Wait()
	close(resultChan)

	for r := range resultChan {
		result.DeadFunctions = append(result.DeadFunctions, r.DeadFunctions...)
		result.UnusedVariables = append(result.UnusedVariables, r.UnusedVariables...)
		result.UnusedImports = append(result.UnusedImports, r.UnusedImports...)
		result.ScannedFiles++
	}

	result.GhostDeps, err = s.scanDependencies()
	if err != nil {
		return nil, err
	}

	result.EstSavedBytes = s.estimateSavings(result)

	return result, nil
}

func (s *Scanner) scanFile(path string) *Result {
	ext := strings.ToLower(filepath.Ext(path))
	if lang, ok := s.reg.Get(ext); ok {
		return lang.Scan(path)
	}
	return nil
}

func (s *Scanner) collectFiles() ([]string, error) {
	var files []string

	err := filepath.Walk(s.cfg.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && s.shouldSkipDir(info.Name()) {
			return filepath.SkipDir
		}

		if !info.IsDir() && s.shouldScanFile(path) {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func (s *Scanner) shouldSkipDir(name string) bool {
	skip := map[string]bool{
		"node_modules": true,
		".git":         true,
		"vendor":       true,
		"dist":         true,
		"build":        true,
		".cache":       true,
		"__pycache__":  true,
		".venv":        true,
		"venv":         true,
	}
	return skip[name]
}

func (s *Scanner) shouldScanFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	if s.cfg.Lang != "auto" {
		lang, ok := s.reg.Get(ext)
		if ok && lang.Name() == s.cfg.Lang {
			return true
		}
		return false
	}
	_, ok := s.reg.Get(ext)
	return ok
}

func (s *Scanner) scanDependencies() ([]GhostDep, error) {
	var deps []GhostDep
	for _, lang := range s.reg.All() {
		d, err := lang.ScanDeps(s.cfg.Path)
		if err == nil && len(d) > 0 {
			deps = append(deps, d...)
		}
	}
	return deps, nil
}

func (s *Scanner) estimateSavings(r *Result) int64 {
	var total int64
	total += int64(len(r.DeadFunctions)) * 512
	total += int64(len(r.UnusedVariables)) * 64
	total += int64(len(r.GhostDeps)) * 1024 * 1024
	return total
}

func (s *Scanner) Clean(r *Result) {
}
