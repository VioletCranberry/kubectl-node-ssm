# kubectl-node-ssm

## Description

`node-ssm` is a straightforward `kubectl` plugin designed for establishing direct connections to EKS cluster nodes managed by AWS Systems Manager. It operates by utilizing the locally installed AWS CLI and session-manager-plugin. The plugin simplifies the process by automatically converting the provided EKS node name into its corresponding instance ID.

### Install with Krew plugin manager

```shell
# see https://krew.sigs.k8s.io/
kubectl krew update
kubectl krew install node-ssm
```

### Usage

```shell
❯ kubectl get nodes --no-headers | head -n 1            
ip-10-10-10-10.ec2.internal   Ready                      <none>   8d      v1.22.17-eks-48e63af
❯ kubectl node-ssm start-session --target ip-10-10-10-10.ec2.internal

Starting session with SessionId: <username>@<domain>-0480532656ed795d8
sh-4.2$ 
```

All global global command-line flags listed in `kubectl options` are supported, for example:

```shell
❯ kubectl config current-context
<my-current-context>
❯ kubectl get nodes --context <my-another-context> --no-headers | head -n 1
ip-20-20-20-20.ec2.internal   Ready   <none>   4d19h   v1.22.17-eks-48e63af
❯ kubectl node-ssm start-session --context <my-another-context> --target ip-20-20-20-20.ec2.internal 

Starting session with SessionId: <username>@<domain>-0dd10b4b84087dff4
sh-4.2$
```

SSM `start-session` [parameters](https://docs.aws.amazon.com/cli/latest/reference/ssm/start-session.html) can be set with optional `--session-params` flag:

```shell
kubectl node-ssm start-session --target ip-30-30-30-30.ec2.internal --session-params '--reason=test' --session-params '--debug'
2023-04-10 00:54:45,509 - MainThread - awscli.clidriver - DEBUG - CLI version: aws-cli/2.11.0 Python/3.11.2 Darwin/22.3.0 source/arm64
2023-04-10 00:54:45,509 - MainThread - awscli.clidriver - DEBUG - Arguments entered to CLI: ['ssm', 'start-session', '--target', 'i-057750d42936e468a', '--reason=test', '--debug']
...
```

### Build and install manually

```shell
go build -o kubectl-node_ssm \
  && sudo cp kubectl-node_ssm /usr/local/bin \
  && kubectl plugin list | grep node_ssm \
  && kubectl node-ssm --help
# rm -f /usr/local/bin/node_ssm
```

### Requirements

1. Installed [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html) and [AWS session-manager-plugin](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)
2. Configured [AWS Systems Manager Session Manager](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-getting-started.html)
3. IAM Permissions to perform `ec2:DescribeInstances`

### Logic

1. Extract `AWS_REGION` and `AWS_PROFILE` from [Config.Host](https://pkg.go.dev/k8s.io/client-go@v0.26.3/rest#Config.Host) and [[]ExecEnvVar](https://pkg.go.dev/k8s.io/client-go/tools/clientcmd/api#ExecConfig.Env) array of current kubeconfig context.
2. Create [AWS session](https://pkg.go.dev/github.com/aws/aws-sdk-go/aws/session) and resolve EKS node `private-dns-name` to instance ID using [(*EC2) DescribeInstances](https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#EC2.DescribeInstances) API operation.
3. Build `aws ssm start-session --target <instance id>` [command](https://pkg.go.dev/os/exec#Command) with specified parameters and environment and execute it.
