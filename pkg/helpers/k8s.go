package helpers

import (
	"fmt"
	"regexp"

	"k8s.io/client-go/rest"
)

type KubeConfig struct {
	Config     *rest.Config
	AWSProfile string
	AWSRegion  string
}

func (k *KubeConfig) SetProfile() {
	errmsg := fmt.Errorf("can't parse AWS_PROFILE from %v",
		k.Config.ExecProvider.Env)
	for _, envvar := range k.Config.ExecProvider.Env {
		if envvar.Name == "AWS_PROFILE" {
			k.AWSProfile = envvar.Value
		}
	}
	if k.AWSProfile == "" {
		panic(errmsg)
	}
}

func (k *KubeConfig) SetRegion() {
	re := regexp.MustCompile(`\w+-(gov-)?\w+-\d+`)
	errmsg := fmt.Errorf("can't parse AWS_REGION from %s",
		k.Config.Host)
	AWSregion := re.FindString(k.Config.Host)
	if AWSregion != "" {
		k.AWSRegion = AWSregion
	} else {
		panic(errmsg)
	}
}
