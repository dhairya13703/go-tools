package installer

import (
	// "fmt"
	"os/exec"
	"runtime"
)

type InstallerType string

const (
	MySQL    InstallerType = "mysql"
	Postgres InstallerType = "postgres"
	MongoDB  InstallerType = "mongodb"
)

type BaseInstaller struct {
	Name         string
	PackageName  map[string]string // OS specific package names
	DownloadURL  map[string]string // OS specific download URLs
	Instructions map[string]string // OS specific manual instructions
}

func (b *BaseInstaller) IsInstalled() bool {
	_, err := exec.LookPath(b.Name)
	return err == nil
}

func (b *BaseInstaller) GetName() string {
	return b.Name
}

func GetOS() string {
	switch runtime.GOOS {
	case "darwin":
		return "darwin"
	case "linux":
		return "linux"
	case "windows":
		return "windows"
	default:
		return "unknown"
	}
}
