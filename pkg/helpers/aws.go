package helpers

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type AWSClient struct {
	Client ec2iface.EC2API
}

func NewAWSClient(awsProfile, awsRegion string) (*AWSClient, error) {
	sess, err := createAWSSession(awsProfile, awsRegion)
	if err != nil {
		return nil, err
	}

	return &AWSClient{
		Client: ec2.New(sess),
	}, nil
}

func createAWSSession(awsProfile, awsRegion string) (*session.Session, error) {
	var opts session.Options

	if awsProfile == "" {
		opts = session.Options{
			Config: aws.Config{
				Region: aws.String(awsRegion),
			},
			SharedConfigState: session.SharedConfigEnable,
		}
	} else {
		opts = session.Options{
			Profile: awsProfile,
			Config: aws.Config{
				Region: aws.String(awsRegion),
			},
		}
	}

	return session.NewSessionWithOptions(opts)
}

func (c *AWSClient) GetInstanceData(privateDnsName string) (*ec2.DescribeInstancesOutput, error) {
	res, err := c.Client.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("private-dns-name"),
				Values: []*string{aws.String(privateDnsName)},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error describing instances: %w", err)
	}
	if res.Reservations == nil || len(res.Reservations) == 0 {
		return nil, fmt.Errorf("no instance data found for %s", privateDnsName)
	}
	return res, nil
}

func ParseInstanceData(input *ec2.DescribeInstancesOutput) (string, error) {
	if input == nil || len(input.Reservations) == 0 {
		return "", errors.New("no reservations found in the instance data")
	}

	for _, reservation := range input.Reservations {
		for _, instance := range reservation.Instances {
			if instance.InstanceId != nil && *instance.InstanceId != "" {
				return *instance.InstanceId, nil
			}
		}
	}

	return "", errors.New("instance ID was not found in the provided data")
}
