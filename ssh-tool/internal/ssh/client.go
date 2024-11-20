package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"ssh-tool/internal/config"
	"strings"
)

type Client struct {
	Server config.Server
}

func NewClient(server config.Server) *Client {
	return &Client{
		Server: server,
	}
}

func (c *Client) Connect() error {
	// Expand the ~ in the pem file path
	pemFile := c.Server.PemFile
	if strings.HasPrefix(pemFile, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting home directory: %v", err)
		}
		pemFile = filepath.Join(home, pemFile[1:])
	}

	// Check if the pem file exists
	if _, err := os.Stat(pemFile); err != nil {
		return fmt.Errorf("pem file not found: %v", err)
	}

	// Prepare the SSH command
	cmd := exec.Command("ssh",
		"-i", pemFile,
		fmt.Sprintf("%s@%s", c.Server.User, c.Server.Hostname))

	// Set up the command to use the current terminal
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the SSH command
	return cmd.Run()
}
