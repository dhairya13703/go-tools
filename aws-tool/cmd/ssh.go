package cmd

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Start an SSM session with an EC2 instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		instanceID, _ := cmd.Flags().GetString("instance-id")
		if instanceID == "" {
			return fmt.Errorf("error: instance-id is required")
		}

		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")
		debug, _ := cmd.Flags().GetBool("debug")

		return ssmConnect(instanceID, region, profile, debug)
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)
	sshCmd.Flags().StringP("instance-id", "i", "", "EC2 instance ID to connect to")
	sshCmd.Flags().BoolP("debug", "d", false, "Enable debug output")
	sshCmd.MarkFlagRequired("instance-id")
}

func ssmConnect(instanceID, region, profile string, debug bool) error {
	ssmCmd := fmt.Sprintf("aws ssm start-session --target %s --region %s --profile %s", instanceID, region, profile)
	if debug {
		fmt.Println("Executing SSM command:", ssmCmd)
	}

	cmd := exec.Command("sh", "-c", ssmCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start SSM command: %v", err)
	}

	fmt.Println("Attempting to establish SSM session...")

	errChan := make(chan error, 1)
	go func() {
		errChan <- cmd.Wait()
	}()

	select {
	case err := <-errChan:
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					switch status.ExitStatus() {
					case 254:
						return fmt.Errorf("instance is not connected to SSM. Please ensure the instance is running and has the SSM agent installed and configured correctly")
					case 255:
						return fmt.Errorf("failed to start SSM session. Please check your AWS credentials and permissions")
					default:
						return fmt.Errorf("SSM session ended with exit code: %d", status.ExitStatus())
					}
				}
			}
			return fmt.Errorf("SSM command failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		fmt.Println("SSM session established. Waiting for SSH connection...")
		fmt.Println("Once connected, you can start typing commands.")
		fmt.Println("To exit the session, type 'exit' twice or use Ctrl+C.")

		_, err = io.WriteString(stdinPipe, "sudo ssh ec2-user@localhost\n")
		if err != nil {
			return fmt.Errorf("failed to write SSH command: %v", err)
		}

		go func() {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				io.WriteString(stdinPipe, scanner.Text()+"\n")
			}
		}()

		go func() {
			<-sigChan
			cmd.Process.Signal(syscall.SIGTERM)
		}()

		err = <-errChan
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					return fmt.Errorf("Session ended with exit code: %d", status.ExitStatus())
				}
			}
			return fmt.Errorf("Command failed: %v", err)
		}
	}

	fmt.Println("Session ended.")
	return nil
}

func sshConnect(instanceID string, debug bool) error {
	// You might want to implement a way to get the instance's public IP or hostname
	// For now, we'll just use the instance ID as the hostname
	hostname := instanceID

	if !isPortOpen(hostname, 22) {
		return fmt.Errorf("SSH port is not open")
	}

	args := []string{"ssh", fmt.Sprintf("ec2-user@%s", hostname)}

	fmt.Printf("Connecting to %s via SSH...\n", hostname)

	cmd := exec.Command(os.Getenv("SHELL"), []string{"-c", strings.Join(args, " ")}...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if debug {
		fmt.Println(cmd.String())
	}

	return cmd.Run()
}

func isPortOpen(host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, time.Duration(1)*time.Second)
	if conn != nil {
		defer conn.Close()
	}

	return err == nil
}
