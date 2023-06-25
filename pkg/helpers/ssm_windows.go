//go:build windows

package helpers

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/VioletCranberry/kubectl-node-ssm/pkg/utils"
	"golang.org/x/sys/windows"
)

type SSMClient struct {
	CMD *exec.Cmd
}

func (c *SSMClient) SetCMD(targetId string, params []string) {
	cmdArgs := []string{"ssm", "start-session", "--target", targetId}

	cmdArgs = append(cmdArgs, params...)
	cmd := exec.Command("aws", cmdArgs...)

	// Probably that will work. Sadly I do not have access
	// neither to Windows machine nor to SSM-enabled environments
	// right now.

	cmd.SysProcAttr = &windows.SysProcAttr{
		CreationFlags:    windows.CREATE_NEW_CONSOLE,
		NoInheritHandles: true,
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
