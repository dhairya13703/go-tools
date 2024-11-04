package installer

import (
	"fmt"
	"os/exec"
)

type MySQLDumpInstaller struct {
	BaseInstaller
}

func NewMySQLDumpInstaller() *MySQLDumpInstaller {
	return &MySQLDumpInstaller{
		BaseInstaller: BaseInstaller{
			Name: "mysqldump",
			PackageName: map[string]string{
				"darwin":  "mysql",
				"linux":   "mysql-client",
				"windows": "mysql",
			},
			Instructions: map[string]string{
				"darwin": "brew install mysql",
				"linux":  "sudo apt update && sudo apt install -y mysql-client",
				"windows": `Please download MySQL installer from https://dev.mysql.com/downloads/mysql/
1. During installation, select "MySQL Command Line Utilities"
2. Add MySQL bin directory to your PATH`,
			},
		},
	}
}

func (m *MySQLDumpInstaller) Install() error {
	os := GetOS()
	
	switch os {
	case "darwin":
		return m.brewInstall()
	case "linux":
		return m.aptInstall()
	case "windows":
		return fmt.Errorf("automatic installation not supported on Windows. %s", m.Instructions["windows"])
	default:
		return fmt.Errorf("unsupported operating system")
	}
}

func (m *MySQLDumpInstaller) brewInstall() error {
	cmd := exec.Command("brew", "install", m.PackageName["darwin"])
	return cmd.Run()
}

func (m *MySQLDumpInstaller) aptInstall() error {
	// Try to install without updating first
	installCmd := exec.Command("sudo", "apt", "install", "-y", m.PackageName["linux"])
	output, err := installCmd.CombinedOutput()
	if err == nil {
		return nil
	}

	// If direct install failed, try updating first
	updateCmd := exec.Command("sudo", "apt-get", "update", "--allow-releaseinfo-change")
	updateOutput, updateErr := updateCmd.CombinedOutput()
	if updateErr != nil {
		// Continue anyway, as some repository errors are non-fatal
		fmt.Printf("Warning: apt update encountered issues: %s\n", string(updateOutput))
	}

	// Try installation again
	installCmd = exec.Command("sudo", "apt", "install", "-y", m.PackageName["linux"])
	output, err = installCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install mysql-client: %v\nOutput: %s", err, string(output))
	}

	return nil
}

func (m *MySQLDumpInstaller) IsInstalled() bool {
	// First check if binary exists
	path, err := exec.LookPath(m.Name)
	if err != nil {
		return false
	}

	// Verify mysqldump works by running version check
	cmd := exec.Command(path, "--version")
	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}
