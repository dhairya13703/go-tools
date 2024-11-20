// cmd/root.go
package cmd

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

type PortResult struct {
	Port    int
	State   string
	Service string
}

// Common ports and their services
var commonPorts = map[int]string{
	20:    "FTP-DATA",
	21:    "FTP",
	22:    "SSH",
	23:    "Telnet",
	25:    "SMTP",
	53:    "DNS",
	80:    "HTTP",
	110:   "POP3",
	143:   "IMAP",
	443:   "HTTPS",
	465:   "SMTPS",
	587:   "SMTP",
	993:   "IMAPS",
	995:   "POP3S",
	1433:  "MSSQL",
	3306:  "MySQL",
	3389:  "RDP",
	5432:  "PostgreSQL",
	27017: "MongoDB",
}

var (
	ports    string
	allPorts bool
	timeout  int
	workers  int
	serverIP string
	rootCmd  = &cobra.Command{
		Use:   "portscanner",
		Short: "A fast port scanner written in Go",
		Long: `A port scanner that can scan specific ports or all ports on a given server.
Complete documentation is available at https://github.com/yourusername/portscanner`,
		Run: runScan,
	}
)

func init() {
	rootCmd.Flags().StringVarP(&ports, "ports", "p", "", "Ports to scan (comma-separated, ranges allowed e.g., 80,443,8000-8010)")
	rootCmd.Flags().BoolVarP(&allPorts, "all", "a", false, "Scan all ports (0-65535)")
	rootCmd.Flags().IntVarP(&timeout, "timeout", "t", 2, "Timeout in seconds for each port scan")
	rootCmd.Flags().IntVarP(&workers, "workers", "w", 1000, "Number of concurrent workers")
	rootCmd.Flags().StringVarP(&serverIP, "server", "s", "", "Server IP address to scan")
	rootCmd.MarkFlagRequired("server")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func parsePorts(portsFlag string) ([]int, error) {
	var ports []int
	if portsFlag == "" {
		return ports, nil
	}

	ranges := strings.Split(portsFlag, ",")
	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if strings.Contains(r, "-") {
			parts := strings.Split(r, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid port range: %s", r)
			}

			start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
			if err != nil {
				return nil, fmt.Errorf("invalid start port: %s", parts[0])
			}

			end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid end port: %s", parts[1])
			}

			for i := start; i <= end; i++ {
				ports = append(ports, i)
			}
		} else {
			port, err := strconv.Atoi(r)
			if err != nil {
				return nil, fmt.Errorf("invalid port: %s", r)
			}
			ports = append(ports, port)
		}
	}

	return ports, nil
}

func scanPort(host string, port int, timeout time.Duration) *PortResult {
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		return nil
	}
	defer conn.Close()

	service, exists := commonPorts[port]
	if !exists {
		service = "Unknown"
	}

	return &PortResult{
		Port:    port,
		State:   "Open",
		Service: service,
	}
}

func worker(host string, ports <-chan int, results chan<- *PortResult, timeout time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	for port := range ports {
		if result := scanPort(host, port, timeout); result != nil {
			results <- result
		}
	}
}

func runScan(cmd *cobra.Command, args []string) {
	var portsToScan []int
	var err error

	if allPorts {
		portsToScan = make([]int, 65536)
		for i := range portsToScan {
			portsToScan[i] = i
		}
	} else if ports != "" {
		portsToScan, err = parsePorts(ports)
		if err != nil {
			fmt.Printf("Error parsing ports: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Error: either --ports or --all flag must be specified")
		os.Exit(1)
	}

	// Create buffered channels for ports and results
	portsChan := make(chan int, workers)
	resultsChan := make(chan *PortResult, len(portsToScan))

	// Create worker pool
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(serverIP, portsChan, resultsChan, time.Duration(timeout)*time.Second, &wg)
	}

	// Send ports to workers
	go func() {
		for _, port := range portsToScan {
			portsChan <- port
		}
		close(portsChan)
	}()

	// Wait for all workers to complete in a separate goroutine
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect and sort results
	var results []PortResult
	for result := range resultsChan {
		results = append(results, *result)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Port < results[j].Port
	})

	// Print results
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "\nPort\tState\tService\t\n")
	fmt.Fprintf(w, "----\t-----\t-------\t\n")

	for _, result := range results {
		fmt.Fprintf(w, "%d\t%s\t%s\t\n", result.Port, result.State, result.Service)
	}
	w.Flush()

	if len(results) == 0 {
		fmt.Println("\nNo open ports found.")
	} else {
		fmt.Printf("\nFound %d open ports\n", len(results))
	}
}
