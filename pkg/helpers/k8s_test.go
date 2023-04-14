package helpers

import (
	"testing"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"
)

func Test_SetRegion(t *testing.T) {
	type testCases struct {
		config         *KubeConfig
		expectedRegion string
	}
	for _, scenario := range []testCases{
		{
			config: &KubeConfig{
				Config: &rest.Config{
					Host: "https://us-east-1.eks.amazon.com",
				},
			},
			expectedRegion: "us-east-1",
		},
		{
			config: &KubeConfig{Config: &rest.Config{
				Host: "https://us-gov-east-1.eks.amazon.com",
			}},
			expectedRegion: "us-gov-east-1",
		},
	} {
		scenario.config.SetRegion()
		if scenario.expectedRegion != scenario.config.AWSRegion {
			t.Errorf("set invalid region: expected %s / got %s",
				scenario.expectedRegion, scenario.config.AWSRegion)
		}
	}
}

func Test_SetRegionPanic(t *testing.T) {
	testConfig := KubeConfig{
		Config: &rest.Config{
			Host: "",
		},
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the code did not panic")
		}
	}()
	testConfig.SetRegion()
}

func Test_SetProfile(t *testing.T) {
	type testCases struct {
		config          *KubeConfig
		expectedProfile string
	}

	for _, scenario := range []testCases{
		{
			config: &KubeConfig{Config: &rest.Config{
				ExecProvider: &api.ExecConfig{
					Env: []api.ExecEnvVar{
						{
							Name:  "AWS_PROFILE",
							Value: "test-profile-1",
						},
					},
				},
			}},
			expectedProfile: "test-profile-1",
		},
		{
			config: &KubeConfig{Config: &rest.Config{
				ExecProvider: &api.ExecConfig{
					Env: []api.ExecEnvVar{
						{
							Name:  "AWS_PROFILE",
							Value: "test-profile-2",
						},
					},
				},
			}},
			expectedProfile: "test-profile-2",
		},
	} {
		scenario.config.SetProfile()
		if scenario.expectedProfile != scenario.config.AWSProfile {
			t.Errorf("set invalid profile: expected %s / got %s",
				scenario.expectedProfile, scenario.config.AWSProfile)
		}
	}
}

func Test_SetProfilePanic(t *testing.T) {
	testConfig := KubeConfig{
		Config: &rest.Config{},
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the code did not panic")
		}
	}()
	testConfig.SetProfile()
}
