package cmd

import (
	"fmt"

	"github.com/VioletCranberry/kubectl-node-ssm/pkg/helpers"
	"github.com/spf13/cobra"
)

func newStartSessionCmd(opts *cliOptions, target *string, params *[]string) *cobra.Command {
	startSessionCmd := &cobra.Command{
		Use:   "start-session",
		Short: "Start AWS systems manager session using local AWS CLI and session-manager-plugin",
		Long:  "Start AWS systems manager session using local AWS CLI and session-manager-plugin",
		RunE: func(cmd *cobra.Command, args []string) error {

			kubeConfig, err := readKubeConfig(opts)
			if err != nil {
				return fmt.Errorf("error reading kubeconfig file: %w", err)
			}
			instanceId, err := resolveTargetToId(kubeConfig.AwsProfile, kubeConfig.AwsRegion, *target)
			if err != nil {
				return fmt.Errorf("error resolving target node name to instance ID: %w", err)
			}
			err = newSSMsession(kubeConfig.AwsProfile, kubeConfig.AwsRegion, instanceId, *params)
			if err != nil {
				return fmt.Errorf("error starting new SSM session: %w", err)
			}

			return nil
		},
	}
	return startSessionCmd
}

func readKubeConfig(opts *cliOptions) (*helpers.KubeConfig, error) {
	clientConfig, err := opts.flags.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to create REST config from kubeconfig flags: %w", err)
	}
	kubeConfig, err := helpers.NewKubeConfig(clientConfig)
	if err != nil {
		return nil, err
	}
	return kubeConfig, nil
}

func resolveTargetToId(awsProfile, awsRegion, target string) (string, error) {
	client, err := helpers.NewAWSClient(awsProfile, awsRegion)
	if err != nil {
		return "", fmt.Errorf("error setting up AWS client: %w", err)
	}
	instanceData, err := client.GetInstanceData(target)
	if err != nil {
		return "", fmt.Errorf("error getting instance data for target %s: %w", target, err)
	}
	instanceId, err := helpers.ParseInstanceData(instanceData)
	if err != nil {
		return "", fmt.Errorf("error parsing instance data for target %s: %w", target, err)
	}
	return instanceId, nil
}

func newSSMsession(awsProfile, awsRegion, instanceId string, params []string) error {
	client, err := helpers.NewSSMClient(instanceId, params, awsProfile, awsRegion)
	if err != nil {
		return fmt.Errorf("error setting up SSM client: %w", err)
	}
	err = client.RunCMD()
	if err != nil {
		return fmt.Errorf("error running SSM session command: %w", err)
	}
	return nil
}
