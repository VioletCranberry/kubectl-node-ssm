package helpers

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	ec2 "github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type mockEC2Client struct {
	ec2iface.EC2API
	Res ec2.DescribeInstancesOutput
}

func (m mockEC2Client) DescribeInstances(in *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return &m.Res, nil
}

func Test_GetInstanceData(t *testing.T) {
	type testCases struct {
		res            ec2.DescribeInstancesOutput
		privateDnsName string
	}
	for _, scenario := range []testCases{
		{
			res: ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					{
						Instances: []*ec2.Instance{
							{
								InstanceId:     aws.String("test-instance-id"),
								PrivateDnsName: aws.String("test-private-dns-test"),
								State: &ec2.InstanceState{
									Name: aws.String("running"),
								},
							},
						},
					},
				},
			},
			privateDnsName: "test-private-dns-test",
		},
	} {
		testAwsClient := AWSClient{
			Client: mockEC2Client{Res: scenario.res},
		}
		instanceData := testAwsClient.GetInstanceData(scenario.privateDnsName)
		for _, reservation := range instanceData.Reservations {
			for _, instance := range reservation.Instances {
				if scenario.privateDnsName != *instance.PrivateDnsName {
					t.Errorf("got invalid instance data: expected %s / got %s",
						scenario.privateDnsName, *instance.PrivateDnsName)
				}
			}
		}
	}
}

func Test_GetInstanceDataPanic(t *testing.T) {
	testAwsClient := AWSClient{
		Client: mockEC2Client{
			Res: ec2.DescribeInstancesOutput{},
		},
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the code did not panic")
		}
	}()
	testAwsClient.GetInstanceData("test")
}

func Test_ParseInstanceData(t *testing.T) {
	type testCases struct {
		res        ec2.DescribeInstancesOutput
		instanceId string
	}
	for _, scenario := range []testCases{
		{
			res: ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{
					{
						Instances: []*ec2.Instance{
							{
								InstanceId:     aws.String("test-instance-id"),
								PrivateDnsName: aws.String("test-private-dns"),
								State: &ec2.InstanceState{
									Name: aws.String("running"),
								},
							},
						},
					},
				},
			},
			instanceId: "test-instance-id",
		},
	} {
		instanceId := ParseInstanceData(&scenario.res)
		if instanceId != scenario.instanceId {
			t.Errorf("parsed invalid instance data: expected %s / got %s",
				scenario.instanceId, instanceId)
		}
	}
}

func Test_ParseInstanceDataPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the code did not panic")
		}
	}()
	ParseInstanceData(&ec2.DescribeInstancesOutput{})
}
