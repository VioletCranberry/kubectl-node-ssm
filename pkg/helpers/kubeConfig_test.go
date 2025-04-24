package helpers

import (
	"testing"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"
)

var mockEnv = []api.ExecEnvVar{
	{Name: "SOME_ENV", Value: "some_value"},
	{Name: "AWS_PROFILE", Value: "my_aws_profile"},
	{Name: "ANOTHER_ENV", Value: "another_value"},
}

func TestExtractAwsProfile(t *testing.T) {
	tests := []struct {
		name       string
		execEnvVar []api.ExecEnvVar
		want       string
	}{
		{
			name:       "AWS_PROFILE present",
			execEnvVar: mockEnv,
			want:       "my_aws_profile",
		},
		{
			name:       "AWS_PROFILE absent",
			execEnvVar: []api.ExecEnvVar{{Name: "ANOTHER_ENV", Value: "another_value"}},
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &rest.Config{
				ExecProvider: &api.ExecConfig{
					Env: tt.execEnvVar,
				},
			}
			if got := extractAwsProfile(config); got != tt.want {
				t.Errorf("getAwsProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractAwsRegion(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		want    string
		wantErr bool
	}{
		{
			name:    "Valid AWS region in EKS cluster hostname",
			host:    "https://123456789asdfghjk.gr1.eu-west-2.eks.amazonaws.com",
			want:    "eu-west-2",
			wantErr: false,
		},
		{
			name:    "Valid AWS GOV region in EKS cluster hostname",
			host:    "https://123456789asdfghjk.gr1.us-gov-west-1.eks.amazonaws.com",
			want:    "us-gov-west-1",
			wantErr: false,
		},
		{
			name:    "Valid AWS CN region in EKS cluster hostname",
			host:    "https://123456789asdfghjk.gr1.cn-north-1.eks.amazonaws.com",
			want:    "cn-north-1",
			wantErr: false,
		},
		{
			name:    "No AWS region in ARN",
			host:    "localhost:8080",
			want:    "",
			wantErr: true,
		},
		{
			name:    "Valid AWS region in non-eks cluster",
			host:    "https://api-demo-non-eks.elb.eu-west-1.amazonaws.com?",
			want:    "eu-west-1",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &rest.Config{
				Host: tt.host,
			}
			got, err := extractAwsRegion(config)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAwsRegion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getAwsRegion() = %v, want %v", got, tt.want)
			}
		})
	}
}
