package pkg

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"machine-operator/pkg/aliyun"
	"machine-operator/pkg/aws"
)

type InstanceUtil interface {
	GetInstanceIDsByTag(tags map[string]string) ([]string, error)
	GetInstanceStatusByID(instanceID string) (string, error)
	StartInstance(instanceID string) error
	StopInstance(instanceID string) (string, error)
}

func GetInstanceUtil(platform, region string) (InstanceUtil, error) {
	var instanceUtil InstanceUtil
	switch platform {
	case PlatformAWS:
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
		if err != nil {
			log.Fatal(err)
			return instanceUtil, err
		}
		client := ec2.NewFromConfig(cfg)
		ec2Util := &aws.EC2Util{
			Client: client,
		}
		instanceUtil = ec2Util
	case PlatformAliyun:
		accessKeyId := os.Getenv(AliyunAccessKeyID)
		accessKeySecret := os.Getenv(AliyunAccessKeySecret)
		client, err := ecs.NewClientWithAccessKey(region, accessKeyId, accessKeySecret)
		if err != nil {
			log.Fatal(err)
			return instanceUtil, err
		}
		ecsUtil := &aliyun.ECSUtil{
			Client: client,
		}
		instanceUtil = ecsUtil
	default:
		return instanceUtil, errors.New("unsupported platform type ")
	}

	return instanceUtil, nil
}
