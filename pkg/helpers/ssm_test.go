package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func Test_setCMD(t *testing.T) {
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

func Test_setEnv(t *testing.T) {
	type testCases struct {
		awsProfile string
		awsRegion  string
		instanceId string

		expectedEnv []string
	}

	region := "test-aws-region"
	profile := "test-aws-profile"

	for _, scenario := range []testCases{
		{
			instanceId: "test-instance-id",
			awsProfile: profile,
			awsRegion:  region,
			expectedEnv: append(os.Environ(),
				fmt.Sprintf("AWS_REGION=%s", region),
				fmt.Sprintf("AWS_PROFILE=%s", profile),
			),
		},
		{
			instanceId: "test-instance-id",
			awsProfile: "",
			awsRegion:  region,
			expectedEnv: append(os.Environ(),
				fmt.Sprintf("AWS_REGION=%s", region),
			),
		},
	} {

		ssmClient := SSMClient{}
		ssmClient.SetCMD(scenario.instanceId, []string{})
		ssmClient.SetEnv(scenario.awsProfile, scenario.awsRegion)

		if !reflect.DeepEqual(scenario.expectedEnv, ssmClient.CMD.Env) {
			t.Errorf("set invalid CMD env: expected %s / got %s",
				scenario.expectedEnv, ssmClient.CMD.Env)
		}
	}
}

func Test_runCMDPanic(t *testing.T) {
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
