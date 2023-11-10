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

func NewKubeConfig(config *rest.Config) (*KubeConfig, error) {
	profile := extractAwsProfile(config)
	region, err := extractAwsRegion(config)
	if err != nil {
		return nil, err
	}

	return &KubeConfig{
		Config:     config,
		AwsProfile: profile,
		AwsRegion:  region,
	}, nil
}

func extractAwsProfile(config *rest.Config) string {
	for _, envvar := range config.ExecProvider.Env {
		if envvar.Name == "AWS_PROFILE" {
			return envvar.Value
		}
	}
	return ""
}

func extractAwsRegion(config *rest.Config) (string, error) {
	re := regexp.MustCompile(`https?:\/\/[a-zA-Z0-9.-]+\.([a-z]+-(?:gov-|cn-)?[a-z]+-\d+)\.eks\.amazonaws\.com`)
	matches := re.FindStringSubmatch(config.Host)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("can't parse AWS_REGION from %s", config.Host)
}
