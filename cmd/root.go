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

func NewRootCmd(streams genericclioptions.IOStreams) *cobra.Command {
	var target string
	var params []string
	cliOptions := newCliOptions(streams)

	rootCmd := &cobra.Command{Use: "node-ssm", SilenceUsage: true}
	rootCmd.PersistentFlags().StringVar(&target, "target", "", "EKS node name (private-dns-name)")
	_ = rootCmd.MarkPersistentFlagRequired("target")

	rootCmd.PersistentFlags().StringSliceVar(&params, "session-params",
		[]string{}, "SSM session parameters")

	cliOptions.flags.AddFlags(rootCmd.Flags())
	rootCmd.AddCommand(newStartSessionCmd(cliOptions, &target, &params))

	return rootCmd
}
