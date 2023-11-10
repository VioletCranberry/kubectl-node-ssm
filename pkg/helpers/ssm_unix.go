//go:build linux || darwin

package helpers

import (
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

type SSMClient struct {
	Cmd *exec.Cmd
}

func NewSSMClient(targetId string, params []string,
	awsProfile, awsRegion string) (*SSMClient, error) {
	client := &SSMClient{}
	cmd, err := client.buildCMD(targetId, params)
	if err != nil {
		return nil, fmt.Errorf("error building command: %w", err)
	}

	client.Cmd = cmd
	client.setEnv(awsProfile, awsRegion)

	return client, nil
}

func (c *SSMClient) buildCMD(targetId string, params []string) (*exec.Cmd, error) {
	cmdArgs := append([]string{"ssm", "start-session", "--target", targetId}, params...)
	cmd := exec.Command("aws", cmdArgs...)

	cmd.SysProcAttr = &unix.SysProcAttr{Foreground: true}
	return cmd, nil
}

func (c *SSMClient) setEnv(awsProfile, awsRegion string) {
	env := os.Environ()
	env = append(env, fmt.Sprintf("AWS_REGION=%s", awsRegion))
	if awsProfile != "" {
		env = append(env, fmt.Sprintf("AWS_PROFILE=%s", awsProfile))
	}
	c.Cmd.Env = env
}

func (c *SSMClient) RunCMD() error {
	c.Cmd.Stdin = os.Stdin
	c.Cmd.Stdout = os.Stdout
	c.Cmd.Stderr = os.Stderr

	if err := c.Cmd.Run(); err != nil {
		return fmt.Errorf("error running command: %w", err)
	}
	return nil
}
