
### Description:

`node-ssm` is a dead simple `kubectl` plugin that allows direct connections to AWS EKS cluster Systems Manager managed nodes relying on local AWS CLI and session-manager-plugin installed. Provided EKS node name is automatically resolved to instance ID.

### Usage:
Using current `kubectl` context:
```
❯ kubectl get nodes --no-headers | head -n 1            
ip-10-10-10-10.ec2.internal   Ready                      <none>   8d      v1.22.17-eks-48e63af
❯ kubectl node-ssm start-session --target ip-10-10-10-10.ec2.internal

Starting session with SessionId: <username>@<domain>-0480532656ed795d8
sh-4.2$ 
```

All global `kubectl` options are supported, for example:
```
❯ kubectl config current-context
<my-current-context>
❯ kubectl get nodes --context <my-another-context> --no-headers | head -n 1
ip-20-20-20-20.ec2.internal   Ready   <none>   4d19h   v1.22.17-eks-48e63af
❯ k node-ssm start-session --context <my-another-context> --target ip-20-20-20-20.ec2.internal 

Starting session with SessionId: <username>@<domain>-0dd10b4b84087dff4
sh-4.2$
```

SSM `start-session` [parameters](https://docs.aws.amazon.com/cli/latest/reference/ssm/start-session.html) can be set with optional `--session-params` flag:
```
kubectl node-ssm start-session --target ip-30-30-30-30.ec2.internal --session-params '--reason=test' --session-params '--debug'
2023-04-10 00:54:45,509 - MainThread - awscli.clidriver - DEBUG - CLI version: aws-cli/2.11.0 Python/3.11.2 Darwin/22.3.0 source/arm64
2023-04-10 00:54:45,509 - MainThread - awscli.clidriver - DEBUG - Arguments entered to CLI: ['ssm', 'start-session', '--target', 'i-057750d42936e468a', '--reason=test', '--debug']
...
```

### Build and install manually:
```
❯ go build -o kubectl-node_ssm
❯ cp kubectl-node_ssm /usr/local/bin
❯ kubectl plugin list | grep node_ssm
/usr/local/bin/kubectl-node_ssm
❯ kubectl node-ssm --help
start AWS systems manager session using local AWS CLI and session-manager-plugin

Usage:
  start-session [flags]

Flags:
...
```

### Requirements:

1. [`AWS CLI installed`](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html)
2. [`AWS session-manager-plugin installed`](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)
3. [`AWS Systems Manager Session Manager configured`](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-getting-started.html)
4. IAM Permissions to perform `ec2:DescribeInstances`

### Logic:

1. Return [`kubeconfig`](https://docs.aws.amazon.com/eks/latest/userguide/create-kubeconfig.html) REST client from current context
2. Parse [`Config.Host`](https://pkg.go.dev/k8s.io/client-go@v0.26.3/rest#Config.Host) for `AWS_REGION`
3. Parse [`[]ExecEnvVar`](https://pkg.go.dev/k8s.io/client-go/tools/clientcmd/api#ExecConfig.Env) for `AWS_PROFILE`
4. Resolve EKS node `private-dns-name` to instance ID using [`AWS describe-instances`](https://pkg.go.dev/github.com/aws/aws-sdk-go@v1.44.239/service/ec2#EC2.DescribeInstances)
5. Build `aws ssm start-session --target <instance id>` [command](https://pkg.go.dev/os/exec#Command) with specified parameters
6. Specify environment of the command with `AWS_REGION` and `AWS_PROFILE` variables
7. Start the command and wait for it to complete

### Notes:

Most likely `aws eks update-kubeconfig --region region-code --name my-cluster` should include `--profile`
to have `AWS_PROFILE` variable exposed in cluster environment variables of `kubeconfig` file. I did not test it out since `--profile` was always present when I was creating or updating `kubeconfig` file for our clusters.

