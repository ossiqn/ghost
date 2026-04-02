package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/ossiqn/ghost/internal/reporter"
	"github.com/ossiqn/ghost/internal/scanner"
	"github.com/ossiqn/ghost/internal/scanner/languages"
	"github.com/spf13/cobra"
)

var version = "2.0.0-plugin-arch"

func main() {
	root := &cobra.Command{
		Use:   "ghost",
		Short: "👻 Dead code & unused dependency hunter",
		Long:  banner(),
	}

	scanCmd := &cobra.Command{
		Use:   "scan [path]",
		Short: "Scan project for dead code",
		Args:  cobra.MaximumNArgs(1),
		Run:   runScan,
	}

	scanCmd.Flags().BoolP("json", "j", false, "Output as JSON")
	scanCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	scanCmd.Flags().BoolP("clean", "c", false, "Auto clean after scan")
	scanCmd.Flags().StringP("lang", "l", "auto", "Language (go, javascript, python, ruby, rust, java, auto)")

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("ghost v%s by ossiqn\n", version)
		},
	}

	root.AddCommand(scanCmd, versionCmd)

	if err := root.Execute(); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}
}

func runScan(cmd *cobra.Command, args []string) {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	jsonOut, _ := cmd.Flags().GetBool("json")
	verbose, _ := cmd.Flags().GetBool("verbose")
	autoClean, _ := cmd.Flags().GetBool("clean")
	lang, _ := cmd.Flags().GetString("lang")

	cfg := scanner.Config{
		Path:      path,
		JSON:      jsonOut,
		Verbose:   verbose,
		AutoClean: autoClean,
		Lang:      lang,
	}

	// Dependency Injection ile Plugin Mimarisi
	registry := scanner.NewRegistry()
	registry.Register(&languages.GoLanguage{})
	registry.Register(&languages.JavaScriptLanguage{})
	registry.Register(&languages.PythonLanguage{})
	registry.Register(&languages.RubyLanguage{})
	registry.Register(&languages.RustLanguage{})
	registry.Register(&languages.JavaLanguage{})

	s := scanner.New(cfg, registry)
	result, err := s.Run()
	if err != nil {
		color.Red("Scan failed: %v", err)
		os.Exit(1)
	}

	r := reporter.New(cfg)
	r.Print(result)

	if autoClean || promptClean() {
		s.Clean(result)
	}
}

func promptClean() bool {
	fmt.Print("\nRun 'ghost clean' to remove all? (y/n): ")
	var input string
	fmt.Scanln(&input)
	return input == "y" || input == "Y"
}

func banner() string {
	return `
  ██████  ██░ ██  ▒█████    ██████ ▄▄▄█████▓
▒██    ▒ ▓██░ ██▒▒██▒  ██▒▒██    ▒ ▓  ██▒ ▓▒
░ ▓██▄   ▒██▀▀██░▒██░  ██▒░ ▓██▄   ▒ ▓██░ ▒░
  ▒   ██▒░▓█ ░██ ▒██   ██░  ▒   ██▒░ ▓██▓ ░ 
▒██████▒▒░▓█▒░██▓░ ████▓▒░▒██████▒▒  ▒██▒ ░ 
▒ ▒▓▒ ▒ ░ ▒ ░░▒░▒░ ▒░▒░▒░ ▒ ▒▓▒ ▒ ░  ▒ ░░   
░ ░▒  ░ ░ ▒ ░▒░ ░  ░ ▒ ▒░ ░ ░▒  ░ ░    ░    
░  ░  ░   ░  ░░ ░░ ░ ░ ▒  ░  ░  ░    ░      
      ░   ░  ░  ░    ░ ░        ░           
                                    by ossiqn
  👻 Dead Code & Unused Dependency Hunter
`
}
