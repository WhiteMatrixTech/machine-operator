package aliyun

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"log"
	"strings"
)

type ECSUtil struct {
	Client *ecs.Client
}

func (u *ECSUtil) GetInstanceIDsByTag(tags map[string]string) ([]string, error) {
	instanceIDs := []string{}
	request := ecs.CreateDescribeInstancesRequest()
	tag := []ecs.DescribeInstancesTag{}
	for key, value := range tags {
		tag = append(tag, ecs.DescribeInstancesTag{
			Key:   key,
			Value: value,
		})
	}
	request.Tag = &tag
	response, err := u.Client.DescribeInstances(request)
	if err != nil {
		fmt.Println("Error getting ECS instance :", err)
		return instanceIDs, err
	}
	for _, instance := range response.Instances.Instance {
		instanceIDs = append(instanceIDs, instance.InstanceId)

	}
	return instanceIDs, nil

}

func (u *ECSUtil) GetInstanceStatusByID(instanceID string) (string, error) {
	request := ecs.CreateDescribeInstancesRequest()
	instanceIDs := []string{instanceID}
	instanceIDJson, err := json.Marshal(instanceIDs)
	if err != nil {
		log.Fatal("convert instance id slice to json failed")
		return "", err
	}
	request.InstanceIds = string(instanceIDJson)
	response, err := u.Client.DescribeInstances(request)
	if err != nil {
		log.Fatal("get instance status failed")
		return "", err
	}
	if len(response.Instances.Instance) < 1 {
		log.Fatal("instance not found")
		return "", errors.New("instance not found")
	}

	instanceStatus := strings.ToLower(response.Instances.Instance[0].Status)
	return instanceStatus, nil
}

func (u *ECSUtil) StartInstance(instanceID string) error {
	request := ecs.CreateStartInstanceRequest()
	request.InstanceId = instanceID
	_, err := u.Client.StartInstance(request)
	if err != nil {
		log.Fatal("start instance failed")
		return err
	}
	return nil
}

func (u *ECSUtil) StopInstance(instanceID string) (string, error) {
	request := ecs.CreateStopInstanceRequest()
	request.InstanceId = instanceID
	request.StoppedMode = "StopCharging"
	_, err := u.Client.StopInstance(request)
	if err != nil {
		log.Fatal("failed to stop instance")
		return "", err
	}
	return "stopping", nil
}
