package helpers

import (
	"fmt"
	"os"
	"os/exec"
)

type SSMClient struct {
	AWSProfile string
	AWSRegion  string
	CMD        *exec.Cmd
}

func (c *SSMClient) SetCMD(targetId string, params []string) {
	cmdArgs := []string{"ssm", "start-session", "--target", targetId}

	cmdArgs = append(cmdArgs, params...)
	cmd := exec.Command("aws", cmdArgs...)

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env,
		fmt.Sprintf("AWS_PROFILE=%s", c.AWSProfile),
		fmt.Sprintf("AWS_REGION=%s", c.AWSRegion),
	)
	c.CMD = cmd
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
