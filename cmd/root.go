package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type cliOptions struct {
	flags *genericclioptions.ConfigFlags
	genericclioptions.IOStreams
}

func newCliOptions(streams genericclioptions.IOStreams) *cliOptions {
	return &cliOptions{
		flags:     genericclioptions.NewConfigFlags(true),
		IOStreams: streams,
	}
}

// NewRootCmd creates and returns the root Cobra command for the node-ssm CLI.
func NewRootCmd(streams genericclioptions.IOStreams) *cobra.Command {
	var target string
	var params []string
	cliOpts := newCliOptions(streams)

	rootCmd := &cobra.Command{
		Use:          "node-ssm",
		SilenceUsage: true,
	}

	// Make all ConfigFlags (--kubeconfig, --context, etc.) persistent.
	cliOpts.flags.AddFlags(rootCmd.PersistentFlags())

	// Wire up the IOStreams on the root command.
	rootCmd.SetIn(cliOpts.In)
	rootCmd.SetOut(cliOpts.Out)
	rootCmd.SetErr(cliOpts.ErrOut)

	rootCmd.AddCommand(newStartSessionCmd(cliOpts, &target, &params))
	return rootCmd
}
