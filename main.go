// Package main implements the kubectl-node-ssm plugin entrypoint.
// It ensures the binary is invoked as a kubectl plugin, sets up
// standard IO streams, and executes the root CLI command.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/VioletCranberry/kubectl-node-ssm/cmd"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func main() {
	if !strings.HasPrefix(filepath.Base(os.Args[0]), "kubectl-") {
		fmt.Fprintln(os.Stderr, "This program was not invoked as a kubectl plugin.")
		os.Exit(1)
	}
	newCmd := cmd.NewRootCmd(genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	})
	_ = newCmd.Execute()
}
