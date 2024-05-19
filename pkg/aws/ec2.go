package aws

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2Util struct {
	Client *ec2.Client
}

func (u *EC2Util) GetInstanceIDsByTag(tags map[string]string) ([]string, error) {
	instanceIDs := []string{}
	describeInstancesInput := &ec2.DescribeInstancesInput{}
	filters := []types.Filter{}
	for key, value := range tags {
		filters = append(filters, types.Filter{
			Name:   aws.String("tag:" + key),
			Values: []string{value},
		})
	}
	describeInstancesInput.Filters = filters

	describeInstancesOutput, err := u.Client.DescribeInstances(context.TODO(), describeInstancesInput)
	if err != nil {
		log.Fatal(err)
		return instanceIDs, err
	}
	for _, reservation := range describeInstancesOutput.Reservations {
		for _, instance := range reservation.Instances {
			instanceIDs = append(instanceIDs, *instance.InstanceId)
		}
	}
	return instanceIDs, nil
}

func (u *EC2Util) GetInstanceStatusByID(instanceID string) (string, error) {
	describeInstancesInput := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}
	describeInstancesOutput, err := u.Client.DescribeInstances(context.TODO(), describeInstancesInput)
	if err != nil {
		return "", err
	}
	if len(describeInstancesOutput.Reservations[0].Instances) < 1 {
		return "", errors.New("instance not found")
	}
	instanceStatus := describeInstancesOutput.Reservations[0].Instances[0].State.Name

	return string(instanceStatus), nil
}

func (u *EC2Util) StartInstance(instanceID string) error {
	startInstanceInput := &ec2.StartInstancesInput{
		InstanceIds: []string{instanceID},
	}
	_, err := u.Client.StartInstances(context.TODO(), startInstanceInput)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (u *EC2Util) StopInstance(instanceID string) (string, error) {
	stopInstancesInput := &ec2.StopInstancesInput{
		InstanceIds: []string{instanceID},
	}
	stopInstancesOutput, err := u.Client.StopInstances(context.TODO(), stopInstancesInput)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return string(stopInstancesOutput.StoppingInstances[0].CurrentState.Name), nil
}
