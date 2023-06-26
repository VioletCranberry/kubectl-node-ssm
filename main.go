package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/VioletCranberry/kubectl-node-ssm/cmd"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func main() {
	if strings.HasPrefix(filepath.Base(os.Args[0]), "kubectl-") {
		// nolint: errcheck
		flags := pflag.NewFlagSet("kubectl-node-ssm", pflag.ExitOnError)
		pflag.CommandLine = flags

		sessionCmd := cmd.NewSessionCmd(genericclioptions.IOStreams{
			In:     os.Stdin,
			Out:    os.Stdout,
			ErrOut: os.Stderr,
		})
		err := sessionCmd.Execute()
		if err != nil {
			errmsg := errors.New("was unable to execute cobra session cmd")
			panic(errmsg)
		}
	} else {
		errmsg := errors.New("was not invoked as kubectl plugin")
		panic(errmsg)
	}
}
