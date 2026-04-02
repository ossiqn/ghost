package reporter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/ossiqn/ghost/internal/scanner"
)

type Reporter struct {
	cfg scanner.Config
}

func New(cfg scanner.Config) *Reporter {
	return &Reporter{cfg: cfg}
}

func (r *Reporter) Print(result *scanner.Result) {
	if r.cfg.JSON {
		r.printJSON(result)
		return
	}
	r.printPretty(result)
}

func (r *Reporter) printPretty(result *scanner.Result) {
	cyan := color.New(color.FgCyan, color.Bold)
	red := color.New(color.FgRed, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen, color.Bold)
	white := color.New(color.FgWhite)
	dim := color.New(color.Faint)

	fmt.Println()
	cyan.Println("👻 ghost scan complete")
	dim.Printf("  %d files scanned\n\n", result.ScannedFiles)

	if len(result.DeadFunctions) > 0 {
		red.Println("  Dead Code / Functions:")
		fmt.Println(strings.Repeat("─", 60))
		for _, d := range result.DeadFunctions {
			white.Printf("  %s", d.File)
			dim.Printf(":%d\n", d.Line)
			yellow.Printf("  └─ [%s] %s", d.Language, d.Name)
			dim.Println(" → never called")
		}
		fmt.Println()
	}

	if len(result.UnusedImports) > 0 {
		red.Println("  Unused Imports:")
		fmt.Println(strings.Repeat("─", 60))
		for _, d := range result.UnusedImports {
			white.Printf("  %s", d.File)
			dim.Printf(":%d\n", d.Line)
			yellow.Printf("  └─ [%s] import \"%s\"", d.Language, d.Name)
			dim.Println(" → never used")
		}
		fmt.Println()
	}

	if len(result.UnusedVariables) > 0 {
		red.Println("  Unused Variables:")
		fmt.Println(strings.Repeat("─", 60))
		for _, d := range result.UnusedVariables {
			white.Printf("  %s", d.File)
			dim.Printf(":%d\n", d.Line)
			yellow.Printf("  └─ [%s] var %s", d.Language, d.Name)
			dim.Println(" → never used")
		}
		fmt.Println()
	}

	if len(result.GhostDeps) > 0 {
		red.Println("  Ghost Dependencies:")
		fmt.Println(strings.Repeat("─", 60))
		for _, d := range result.GhostDeps {
			white.Printf("  %s\n", d.File)
			yellow.Printf("  └─ %s", d.Name)
			dim.Printf(" %s", d.Version)
			dim.Println(" → not imported anywhere")
		}
		fmt.Println()
	}

	fmt.Println(strings.Repeat("─", 60))
	total := len(result.DeadFunctions) + len(result.UnusedImports) + len(result.UnusedVariables)

	if total == 0 && len(result.GhostDeps) == 0 {
		green.Println("  ✓ No ghosts found. Clean codebase!")
		return
	}

	red.Printf("  👻 %d dead functions\n", len(result.DeadFunctions))
	red.Printf("  👻 %d unused imports\n", len(result.UnusedImports))
	red.Printf("  👻 %d unused variables\n", len(result.UnusedVariables))
	red.Printf("  👻 %d ghost dependencies\n", len(result.GhostDeps))

	saved := formatBytes(result.EstSavedBytes)
	green.Printf("\n  💾 Est. size saved: ~%s\n", saved)
}

func (r *Reporter) printJSON(result *scanner.Result) {
	data, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(data))
}

func formatBytes(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%dB", bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.1fKB", float64(bytes)/1024)
	}
	return fmt.Sprintf("%.1fMB", float64(bytes)/(1024*1024))
}
