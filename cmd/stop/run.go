package stop

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"
	"log"
	"machine-operator/pkg"
	"machine-operator/pkg/aliyun"
	"machine-operator/pkg/aws"
	"os"
)

var (
	platform   string
	instanceID string
	region     string
	StartCmd   = &cobra.Command{
		Use:          "stop",
		Short:        "stop machine",
		Example:      "machine-operator stop",
		SilenceUsage: true,
		PreRun: func(_ *cobra.Command, _ []string) {
			log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
			preRun()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func preRun() {

}

func init() {
	StartCmd.PersistentFlags().StringVar(&platform,
		"platform", os.Getenv("platform"),
		"the platform")
	StartCmd.PersistentFlags().StringVar(&instanceID,
		"instanceID", os.Getenv("instanceID"),
		"id of machine")
	StartCmd.PersistentFlags().StringVar(&region,
		"region", os.Getenv("region"),
		"the region of aws config")
}

func run() error {
	if platform == "" {
		log.Println("missing 'platform' parameter")
		return errors.New("missing platform parameter")
	}
	if region == "" {
		log.Println("missing 'region' parameter")
		return errors.New("missing region parameter")
	}
	if instanceID == "" {
		log.Println("missing 'instanceID' parameter")
		return errors.New("missing instanceID parameter")
	}

	var machineUtil pkg.InstanceUtil

	switch platform {
	case pkg.PlatformAWS:
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
		if err != nil {
			log.Fatal(err)
			return err
		}
		client := ec2.NewFromConfig(cfg)
		ec2Util := &aws.EC2Util{
			Client: client,
		}
		machineUtil = ec2Util
	case pkg.PlatformAliyun:
		accessKeyId := os.Getenv(pkg.AliyunAccessKeyID)
		accessKeySecret := os.Getenv(pkg.AliyunAccessKeySecret)
		client, err := ecs.NewClientWithAccessKey(region, accessKeyId, accessKeySecret)
		if err != nil {
			log.Fatal(err)
			return err
		}
		ecsUtil := &aliyun.ECSUtil{
			Client: client,
		}
		machineUtil = ecsUtil
	default:
		return errors.New("unsupported platform type ")
	}

	instanceState, err := machineUtil.StopInstance(instanceID)
	if err != nil {
		return err
	}
	fmt.Println("------------------stop the machine--------------------")
	fmt.Println(instanceState)
	return nil
}
