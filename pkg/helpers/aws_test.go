package helpers

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEC2Client struct {
	ec2iface.EC2API
	mock.Mock
}

func (m *MockEC2Client) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*ec2.DescribeInstancesOutput), args.Error(1)
}

func TestNewAWSClient(t *testing.T) {
	client, err := NewAWSClient("", "us-west-2")
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestGetInstanceData(t *testing.T) {
	mockEc2 := new(MockEC2Client)
	testClient := &AWSClient{Client: mockEc2}

	dnsName := "ip-10-0-0-1.ec2.internal"
	expectedOutput := &ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{
					{
						InstanceId:     aws.String("i-1234567890abcdef0"),
						PrivateDnsName: aws.String(dnsName),
						State: &ec2.InstanceState{
							Name: aws.String("running"),
						},
					},
				},
			},
		},
	}

	mockEc2.On("DescribeInstances", mock.Anything).Return(expectedOutput, nil)

	output, err := testClient.GetInstanceData(dnsName)

	mockEc2.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)
}

func TestGetInstanceDataWithError(t *testing.T) {
	mockEc2 := new(MockEC2Client)
	testClient := &AWSClient{Client: mockEc2}

	dnsName := "nonexistent.ec2.internal"
	mockEc2.On("DescribeInstances", mock.Anything).Return((*ec2.DescribeInstancesOutput)(nil), errors.New("instance not found"))

	_, err := testClient.GetInstanceData(dnsName)

	mockEc2.AssertExpectations(t)
	assert.Error(t, err)
}

func TestParseInstanceData(t *testing.T) {
	instanceID := "i-1234567890abcdef0"
	input := &ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{
					{
						InstanceId:     aws.String("i-1234567890abcdef0"),
						PrivateDnsName: aws.String("ip-10-0-0-1.ec2.internal"),
						State: &ec2.InstanceState{
							Name: aws.String("running"),
						},
					},
				},
			},
		},
	}

	id, err := ParseInstanceData(input)

	assert.NoError(t, err)
	assert.Equal(t, instanceID, id)
}

func TestParseInstanceDataWithNoInstances(t *testing.T) {
	input := &ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{},
	}

	_, err := ParseInstanceData(input)

	assert.Error(t, err)
}
