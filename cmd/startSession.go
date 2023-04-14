package cmd

import (
	"github.com/VioletCranberry/kubectl-node-ssm/pkg/helpers"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type CmdOpts struct {
	configFlags *genericclioptions.ConfigFlags
	genericclioptions.IOStreams
}

func NewCmdOpts(streams genericclioptions.IOStreams) *CmdOpts {
	return &CmdOpts{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
	}
}

var (
	target string
	params []string
)

func NewSessionCmd(streams genericclioptions.IOStreams) *cobra.Command {
	opts := NewCmdOpts(streams)

	cmd := &cobra.Command{
		Use:   "start-session",
		Short: "start AWS systems manager session using local AWS CLI and session-manager-plugin",
		Long:  "start AWS systems manager session using local AWS CLI and session-manager-plugin",
		Run: func(cmd *cobra.Command, args []string) {
			opts.ssmConnect(cmd, args)
		}}

	cmd.PersistentFlags().StringVar(&target, "target", "", "node name (required)")
	// nolint: errcheck
	cmd.MarkPersistentFlagRequired("target")
	cmd.PersistentFlags().StringSliceVar(&params, "session-params", []string{}, "ssm session parameters")

	opts.configFlags.AddFlags(cmd.Flags())
	return cmd
}

func (opts *CmdOpts) ssmConnect(cmd *cobra.Command, args []string) {
	config, _ := opts.configFlags.ToRESTConfig()

	kubeConfig := helpers.KubeConfig{Config: config}
	kubeConfig.SetProfile()
	kubeConfig.SetRegion()

	awsClient := helpers.AWSClient{
		AWSProfile: kubeConfig.AWSProfile,
		AWSRegion:  kubeConfig.AWSRegion,
	}
	awsClient.SetClient()
	instanceData := awsClient.GetInstanceData(target)
	instanceId := helpers.ParseInstanceData(instanceData)

	ssmClient := helpers.SSMClient{
		AWSProfile: kubeConfig.AWSProfile,
		AWSRegion:  kubeConfig.AWSRegion,
	}
	ssmClient.SetCMD(instanceId, params)
	ssmClient.RunCMD()

}
