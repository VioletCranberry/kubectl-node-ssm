package helpers

import (
	"fmt"
	"regexp"

	"k8s.io/client-go/rest"
)

type KubeConfig struct {
	Config     *rest.Config
	AwsProfile string
	AwsRegion  string
}

func NewKubeConfig(config *rest.Config) *KubeConfig {
	kubeConfig := KubeConfig{
		Config: config,
	}
	kubeConfig.setAwsRegion()
	kubeConfig.setAwsProfile()
	return &kubeConfig
}

func (k *KubeConfig) setAwsProfile() {
	// AWS_PROFILE can be empty if user is not providing AWS_PROFILE var (e.g. using default profile)
	var awsProfile string
	for _, envvar := range k.Config.ExecProvider.Env {
		if envvar.Name == "AWS_PROFILE" {
			awsProfile = envvar.Value
		}
	}
	k.AwsProfile = awsProfile
}

func (k *KubeConfig) setAwsRegion() {
	re := regexp.MustCompile(`\w+-(gov-)?\w+-\d+`)
	errmsg := fmt.Errorf("can't parse AWS_REGION from %s", k.Config.Host)
	AWSregion := re.FindString(k.Config.Host)
	if AWSregion != "" {
		k.AwsRegion = AWSregion
	} else {
		panic(errmsg)
	}
}
