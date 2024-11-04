package mysql

import (
	"bytes"
	"db-dumper/db/config"
	"db-dumper/db/installer"
	"db-dumper/db/utils"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type MySQLBackup struct{}

func NewMySQLBackup() *MySQLBackup {
	return &MySQLBackup{}
}

// checkDependencies verifies if required utilities are installed
func (m *MySQLBackup) checkDependencies() error {
	mysqlInstaller := installer.NewMySQLDumpInstaller()

	if !mysqlInstaller.IsInstalled() {
		fmt.Printf("%s not found. Attempting to install...\n", mysqlInstaller.GetName())
		if err := mysqlInstaller.Install(); err != nil {
			return fmt.Errorf("failed to install %s: %v", mysqlInstaller.GetName(), err)
		}
		fmt.Printf("%s installed successfully\n", mysqlInstaller.GetName())
	}

	// Check gzip
	if _, err := exec.LookPath("gzip"); err != nil {
		return fmt.Errorf("gzip is not installed. Please install it manually")
	}

	return nil
}

func (m *MySQLBackup) Backup(hostname, username, password, port, database, output string) error {
	// Load config
	config, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	// Use config defaults if not provided
	if hostname == "" {
		hostname = config.DefaultHost
	}
	if port == "" {
		port = config.DefaultPort
	}
	if output == "" {
		output = utils.GenerateBackupFilename("mysql", database, config.BackupDirectory)
	}

	// Check dependencies first
	if err := m.checkDependencies(); err != nil {
		return err
	}

	// If password is empty, prompt user
	if password == "" {
		var err error
		password, err = utils.GetPasswordFromUser("MySQL")
		if err != nil {
			return err
		}
	}

	// Ensure output directory exists
	if err := utils.EnsureOutputDirectory(filepath.Dir(output)); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Format hostname
	hostname = utils.FormatHostname(hostname)

	// Test connection first
	err = utils.TestConnection("MySQL", "mysqladmin",
		"-h", hostname,
		"-P", port,
		"-u", username,
		fmt.Sprintf("-p%s", password),
		"ping")
	if err != nil {
		return err
	}

	// Create mysqldump command with protocol flag
	cmd := exec.Command("mysqldump",
		"-h", hostname,
		"-P", port,
		"-u", username,
		"--protocol=TCP",
		fmt.Sprintf("-p%s", password),
		database)

	// Capture error output
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Create output file
	outFile, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Setup progress reporting
	progress := utils.NewProgressReporter()
	defer progress.Finish()

	// Setup pipes
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	gzipCmd := exec.Command("gzip")
	gzipCmd.Stdout = outFile

	// Pipe through progress reporter
	gzipCmd.Stdin = io.TeeReader(stdout, progress)

	if err := gzipCmd.Start(); err != nil {
		return fmt.Errorf("failed to start gzip: %v", err)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run mysqldump: %v\nError output: %s", err, stderr.String())
	}

	if err := gzipCmd.Wait(); err != nil {
		return fmt.Errorf("failed to finish gzip: %v", err)
	}

	// Print restore instructions
	fmt.Println(utils.GenerateRestoreInstructions(
		"mysql",
			hostname,
			port,
			username,
			database,
			output,
	))

	// After successful backup, rotate old backups
	if err := utils.RotateBackups(config.BackupDirectory, 5); err != nil {
		fmt.Printf("Warning: failed to rotate old backups: %v\n", err)
	}

	return nil
}
