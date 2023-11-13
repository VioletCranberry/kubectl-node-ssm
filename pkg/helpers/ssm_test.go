package helpers

import (
	"os/exec"
	"testing"
)

var mockExecCommand func(name string, arg ...string) *exec.Cmd

func TestNewSsmClient(t *testing.T) {
	mockExecCommand = func(name string, arg ...string) *exec.Cmd {
		return &exec.Cmd{}
	}
	defer func() { mockExecCommand = nil }()

	targetID := "test-target-id"
	params := []string{"param1", "param2"}
	awsProfile := "test-profile"
	awsRegion := "test-region"

	client, err := NewSsmClient(targetID, params, awsProfile, awsRegion)
	if err != nil {
		t.Fatalf("NewSSMClient() error = %v, wantErr %v", err, nil)
	}

	if client.Cmd == nil {
		t.Errorf("NewSSMClient() Cmd is nil, want *exec.Cmd")
	}
}
