// cmd/list.go
package cmd

import (
	"fmt"
	"ssh-tool/internal/config"
	"strings"

	"github.com/spf13/cobra"
)

const (
	colorReset   = "\033[0m"
	colorBlue    = "\033[34m"
	colorCyan    = "\033[36m"
	colorGreen   = "\033[32m"
	colorMagenta = "\033[35m"
	colorYellow  = "\033[33m"
	colorBold    = "\033[1m"
)

type columnConfig struct {
	name      string
	width     int
	formatter func(string) string
}

func formatKeyPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return path
}

func truncateString(str string, length int) string {
	if length <= 3 || len(str) <= length {
		return str
	}
	return str[:length-3] + "..."
}

func colorize(text, color string) string {
	return color + text + colorReset
}

var (
	showAll bool
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all available servers",
		Run:   runList,
	}
)

func runList(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	var columns []columnConfig

	if showAll {
		// Detailed view columns
		columns = []columnConfig{
			{name: "ID", width: 4,
				formatter: func(s string) string { return colorize(s, colorCyan) }},
			{name: "SERVER NAME", width: 24,
				formatter: func(s string) string { return colorize(s, colorGreen) }},
			{name: "IP", width: 14,
				formatter: func(s string) string { return colorize(s, colorMagenta) }},
			{name: "USER", width: 6,
				formatter: func(s string) string { return colorize(s, colorYellow) }},
			{name: "KEY FILE", width: 24,
				formatter: func(s string) string { return colorize(s, colorBlue) }},
		}
	} else {
		// Default minimal view columns
		columns = []columnConfig{
			{name: "ID", width: 4,
				formatter: func(s string) string { return colorize(s, colorCyan) }},
			{name: "SERVER NAME", width: 24,
				formatter: func(s string) string { return colorize(s, colorGreen) }},
			{name: "IP", width: 14,
				formatter: func(s string) string { return colorize(s, colorMagenta) }},
		}
	}

	fmt.Printf("\nAvailable Servers:\n")
	printBorder(columns)

	// Print header row
	fmt.Printf("|")
	for _, col := range columns {
		fmt.Printf(" %-*s |", col.width, col.name)
	}
	fmt.Printf(" %s\n", "CUSTOMERS")

	printBorder(columns)

	// Print data rows
	servers := cfg.GetServersList()
	for i, server := range servers {
		fmt.Printf("|")

		// Print each column based on view type
		if showAll {
			fmt.Printf(" %-*s |", columns[0].width,
				columns[0].formatter(fmt.Sprintf("%d", i+1)))

			fmt.Printf(" %-*s |", columns[1].width,
				columns[1].formatter(truncateString(server.Name, columns[1].width)))

			fmt.Printf(" %-*s |", columns[2].width,
				columns[2].formatter(server.Hostname))

			fmt.Printf(" %-*s |", columns[3].width,
				columns[3].formatter(server.User))

			fmt.Printf(" %-*s |", columns[4].width,
				columns[4].formatter(truncateString(formatKeyPath(server.PemFile), columns[4].width)))
		} else {
			fmt.Printf(" %-*s |", columns[0].width,
				columns[0].formatter(fmt.Sprintf("%d", i+1)))

			fmt.Printf(" %-*s |", columns[1].width,
				columns[1].formatter(truncateString(server.Name, columns[1].width)))

			fmt.Printf(" %-*s |", columns[2].width,
				columns[2].formatter(server.Hostname))
		}

		// Always show customers
		fmt.Printf(" %s", colorize(server.Description, colorCyan))
		fmt.Printf("\n")
	}

	printBorder(columns)

	// Print usage based on view type
	if showAll {
		fmt.Printf("\nUsage: ssh-tool connect <id>\n")
		fmt.Printf("Use 'ssh-tool list' for minimal view\n\n")
	} else {
		fmt.Printf("\nUsage: ssh-tool connect <id>\n")
		fmt.Printf("Use 'ssh-tool list -a' for detailed view\n\n")
	}
}

func printBorder(columns []columnConfig) {
	fmt.Printf("|")
	for _, col := range columns {
		fmt.Printf("-%s-|", strings.Repeat("-", col.width))
	}
	fmt.Printf("---\n")
}

func init() {
	listCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all fields")
	rootCmd.AddCommand(listCmd)
}
