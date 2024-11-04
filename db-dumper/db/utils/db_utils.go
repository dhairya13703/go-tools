package utils

import (
	"fmt"
	"syscall"
	"golang.org/x/term"
	"os/exec"
	"time"
	"path/filepath"
	"os"
)

// GetPasswordFromUser prompts user for password securely
func GetPasswordFromUser(dbType string) (string, error) {
	fmt.Printf("Enter %s password: ", dbType)
	password, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // Add newline after password input
	if err != nil {
		return "", fmt.Errorf("failed to read password: %v", err)
	}
	return string(password), nil
}

// TestConnection attempts to connect to a database
func TestConnection(dbType string, args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to connect to %s: %v\nOutput: %s", dbType, err, string(output))
	}
	return nil
}

// FormatHostname converts localhost to IP if needed
func FormatHostname(hostname string) string {
	if hostname == "localhost" {
		return "127.0.0.1"
	}
	return hostname
}

// GenerateBackupFilename creates a filename with timestamp and database name
func GenerateBackupFilename(dbType, dbName, outputDir string) string {
	timestamp := time.Now().Format("2006-01-02_150405") // YYYY-MM-DD_HHMMSS
	filename := fmt.Sprintf("%s_%s_%s.sql.gz", dbType, dbName, timestamp)
	
	if outputDir != "" {
		return filepath.Join(outputDir, filename)
	}
	return filename
}

// EnsureOutputDirectory creates the output directory if it doesn't exist
func EnsureOutputDirectory(path string) error {
	if path == "" {
		return nil
	}
	return os.MkdirAll(path, 0755)
}

func GenerateRestoreInstructions(dbType, hostname, port, username, database, backupFile string) string {
	switch dbType {
	case "mysql":
		return fmt.Sprintf(`
Backup completed successfully: %s

To restore this backup, use the following command:
    gunzip < %s | mysql -h %s -P %s -u %s %s

Note: You will be prompted for the password during restore.
`, backupFile, backupFile, hostname, port, username, database)
	
	// Add cases for other database types here
	default:
		return fmt.Sprintf("Backup completed successfully: %s", backupFile)
	}
}
