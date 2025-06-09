// Example program demonstrating custom color themes in dyff
package main

import (
	"fmt"
	"os"

	"github.com/gonvenience/ytbx"
	"github.com/homeport/dyff/pkg/dyff"
	"github.com/lucasb-eyer/go-colorful"
)

func main() {
	// Create sample YAML content
	fromYAML := `
name: example
version: 1.0.0
features:
  - authentication
  - logging
config:
  timeout: 30
  retries: 3
`

	toYAML := `
name: example
version: 2.0.0
features:
  - authentication
  - monitoring
  - caching
config:
  timeout: 60
  retries: 5
  max_connections: 100
`

	// Load YAML content
	from, err := ytbx.LoadYAMLDocuments([]byte(fromYAML))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading 'from' YAML: %v\n", err)
		os.Exit(1)
	}

	to, err := ytbx.LoadYAMLDocuments([]byte(toYAML))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading 'to' YAML: %v\n", err)
		os.Exit(1)
	}

	// Create input files
	fromFile := ytbx.InputFile{
		Location:  "from.yaml",
		Documents: from,
	}

	toFile := ytbx.InputFile{
		Location:  "to.yaml",
		Documents: to,
	}

	// Generate comparison report
	report, err := dyff.CompareInputFiles(fromFile, toFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error comparing files: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Default Colors ===")
	// Create human report with default colors
	defaultReport := &dyff.HumanReport{
		Report: report,
	}
	if err := defaultReport.WriteReport(os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing default report: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n=== Custom Colors (High Contrast) ===")
	// Create custom high-contrast color theme
	highContrastTheme := &dyff.ColorTheme{
		Addition:     colorful.Color{R: 0.0, G: 1.0, B: 0.0},   // Pure green
		Modification: colorful.Color{R: 1.0, G: 0.65, B: 0.0},  // Orange
		Removal:      colorful.Color{R: 1.0, G: 0.0, B: 0.0},   // Pure red
	}

	// Create human report with custom colors
	customReport := &dyff.HumanReport{
		Report:     report,
		ColorTheme: highContrastTheme,
	}
	if err := customReport.WriteReport(os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing custom report: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n=== Custom Colors (Colorblind Friendly) ===")
	// Create colorblind-friendly theme (using blue/orange instead of red/green)
	colorblindTheme := &dyff.ColorTheme{
		Addition:     colorful.Color{R: 0.0, G: 0.45, B: 0.70},  // Blue
		Modification: colorful.Color{R: 0.90, G: 0.60, B: 0.0},  // Orange
		Removal:      colorful.Color{R: 0.80, G: 0.40, B: 0.0},  // Dark orange
	}

	// Create brief report with colorblind-friendly theme
	briefReport := &dyff.BriefReport{
		Report:     report,
		ColorTheme: colorblindTheme,
	}
	fmt.Println("Brief report:")
	if err := briefReport.WriteReport(os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing brief report: %v\n", err)
		os.Exit(1)
	}
}