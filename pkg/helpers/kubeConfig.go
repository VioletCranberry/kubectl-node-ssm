package helpers

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"k8s.io/client-go/rest"
)

// KubeConfig holds a Kubernetes REST configuration along with AWS context
// (profile and region) extracted from the kubeconfig.
type KubeConfig struct {
	Config     *rest.Config
	AwsProfile string
	AwsRegion  string
}

// NewKubeConfig creates a KubeConfig instance from the given rest.Config.
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
	regionRe := regexp.MustCompile(`^[a-z]{2}(?:-[a-z]+)+-\d+$`)
	u, err := url.Parse(config.Host)
	if err != nil {
		return "", fmt.Errorf("invalid AWS host URL %q: %w", config.Host, err)
	}
	host := u.Hostname()
	for label := range strings.SplitSeq(host, ".") {
		if regionRe.MatchString(label) {
			return label, nil
		}
	}
	return "", fmt.Errorf("could not parse AWS region from host %q", host)
}
