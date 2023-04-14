package helpers

import (
	"os"
	"os/exec"
	"testing"
)

func Test_SetCmdArgs(t *testing.T) {
	type testCases struct {
		args       []string
		params     []string
		instanceId string
	}
	for _, scenario := range []testCases{
		{
			args:       []string{"aws", "ssm", "start-session", "--target"},
			params:     []string{},
			instanceId: "test-instance-id",
		},
		{
			args:       []string{"aws", "ssm", "start-session", "--target"},
			params:     []string{"--debug"},
			instanceId: "test-instance-id",
		},
	} {
		expectedArgs := append(scenario.args, scenario.instanceId)
		expectedArgs = append(expectedArgs, scenario.params...)

		ssmClient := SSMClient{}
		ssmClient.SetCMD(scenario.instanceId, scenario.params)

		if len(expectedArgs) != len(ssmClient.CMD.Args) {
			t.Errorf("set invalid CMD args: expected %s / got %s",
				expectedArgs, ssmClient.CMD.Args)
		}
	}
}

func Test_SetCmdEnv(t *testing.T) {
	type testCases struct {
		awsProfile string
		awsRegion  string
		env        []string
		instanceId string
	}
	for _, scenario := range []testCases{
		{
			awsProfile: "test-aws-profile",
			awsRegion:  "us-east-1",
			env:        os.Environ(),
			instanceId: "test-instance-id",
		},
		{
			env:        os.Environ(),
			instanceId: "test-instance-id",
		},
	} {
		expectedEnv := append(scenario.env, scenario.awsProfile, scenario.awsRegion)

		ssmClient := SSMClient{
			AWSProfile: scenario.awsProfile,
			AWSRegion:  scenario.awsRegion,
		}
		ssmClient.SetCMD(scenario.instanceId, []string{})

		if len(expectedEnv) != len(ssmClient.CMD.Env) {
			t.Errorf("set invalid CMD env: expected %s / got %s",
				expectedEnv, ssmClient.CMD.Env)
		}
	}
}

func TestRunCMDPanic(t *testing.T) {
	testSSMClient := SSMClient{
		CMD: &exec.Cmd{
			Args: []string{"false"},
		},
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the code did not panic")
		}
	}()
	testSSMClient.RunCMD()
}
