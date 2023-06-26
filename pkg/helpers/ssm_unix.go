//go:build linux || darwin

package helpers

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/VioletCranberry/kubectl-node-ssm/pkg/utils"
	"golang.org/x/sys/unix"
)

type SSMClient struct {
	CMD *exec.Cmd
}

func (c *SSMClient) SetCMD(targetId string, params []string) {
	cmdArgs := []string{"ssm", "start-session", "--target", targetId}

	cmdArgs = append(cmdArgs, params...)
	cmd := exec.Command("aws", cmdArgs...)

	// Put the child processes in the foreground and their own process group to
	// allow the child process group to capture the Ctrl-C (or SIGINT) signal,
	// which otherwise would have killed the node-ssm process and its child
	// processes when they are all in the same process group.

	cmd.SysProcAttr = &unix.SysProcAttr{
		Foreground: true,
	}
	c.CMD = cmd
}

func (c *SSMClient) SetEnv(awsProfile, awsRegion string) {
	c.CMD.Env = os.Environ()
	// aws region is always defined at this stage
	c.CMD.Env = append(c.CMD.Env, fmt.Sprintf("AWS_REGION=%s", awsRegion))

	if !utils.ContainsEmpty(awsProfile) {
		c.CMD.Env = append(c.CMD.Env,
			fmt.Sprintf("AWS_PROFILE=%s", awsProfile),
		)
	}
}

func (c *SSMClient) RunCMD() {
	c.CMD.Stdin = os.Stdin
	c.CMD.Stdout = os.Stdout
	c.CMD.Stderr = os.Stderr

	err := c.CMD.Run()
	if err != nil {
		errmsg := fmt.Errorf("can't run local command: %s ", err)
		panic(errmsg)
	}
}
