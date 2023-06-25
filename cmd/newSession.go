package cmd

import (
	"github.com/VioletCranberry/kubectl-node-ssm/pkg/helpers"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type Opts struct {
	configFlags *genericclioptions.ConfigFlags
	genericclioptions.IOStreams
}

func NewCmdOpts(streams genericclioptions.IOStreams) *Opts {
	return &Opts{
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
			opts.ssmConnect()
		},
	}

	cmd.PersistentFlags().StringVar(&target, "target", "", "node name (required)")
	// nolint: errcheck
	cmd.MarkPersistentFlagRequired("target")
	cmd.PersistentFlags().StringSliceVar(&params, "session-params", []string{}, "ssm session parameters")

	opts.configFlags.AddFlags(cmd.Flags())
	return cmd
}

func (opts *Opts) ssmConnect() {
	ssmClient := helpers.SSMClient{}

	config, _ := opts.configFlags.ToRESTConfig()
	kubeConfig := helpers.NewKubeConfig(config)

	awsClient := helpers.NewAWSClient(
		kubeConfig.AwsProfile,
		kubeConfig.AwsRegion,
	)

	instanceData := awsClient.GetInstanceData(target)
	instanceId := helpers.ParseInstanceData(instanceData)

	ssmClient.SetCMD(instanceId, params)
	ssmClient.SetEnv(kubeConfig.AwsProfile, kubeConfig.AwsRegion)
	ssmClient.RunCMD()
}
