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
	Client     ec2iface.EC2API
	AWSProfile string
	AWSRegion  string
}

func (c *AWSClient) SetClient() {
	sess := session.Must(
		session.NewSessionWithOptions(
			session.Options{
				Profile: c.AWSProfile,
				Config: aws.Config{
					Region: aws.String(c.AWSRegion),
				},
			}),
	)
	c.Client = ec2.New(sess)
}

func (c *AWSClient) GetInstanceData(privateDnsName string) *ec2.DescribeInstancesOutput {
	res, err := c.Client.DescribeInstances(
		&ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name: aws.String("private-dns-name"),
					Values: []*string{
						aws.String(privateDnsName),
					},
				},
				{
					Name: aws.String("instance-state-name"),
					Values: []*string{
						aws.String("running"),
					},
				},
			},
		})
	if err != nil {
		panic(err)
	}
	if res.Reservations == nil {
		errmsg := fmt.Errorf("no instance data found for %s", privateDnsName)
		panic(errmsg)
	}
	return res
}

func ParseInstanceData(input *ec2.DescribeInstancesOutput) string {
	var instanceId string

	for _, reservation := range input.Reservations {
		for _, instance := range reservation.Instances {
			instanceId = *instance.InstanceId
		}
	}
	if instanceId == "" {
		errmsg := errors.New("instance id was not found in provided data")
		panic(errmsg)
	}
	return instanceId
}
