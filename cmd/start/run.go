package start

import (
	"context"
	"errors"
	"log"
	"machine-operator/pkg"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"
)

var (
	region           string
	serverLabelKey   string
	serverLabelValue string
	StartCmd         = &cobra.Command{
		Use:          "start",
		Short:        "start machine",
		Example:      "machine-operator start",
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

func init() {
	StartCmd.PersistentFlags().StringVar(&region,
		"region", os.Getenv("region"),
		"the region of aws config")
	StartCmd.PersistentFlags().StringVar(&serverLabelKey,
		"serverLabelKey", os.Getenv("serverLabelKey"),
		"the key of server label")
	StartCmd.PersistentFlags().StringVar(&serverLabelValue,
		"serverLabelValue", os.Getenv("serverLabelValue"),
		"the value of server label")
}

func preRun() {
	if serverLabelKey == "" {
		serverLabelKey = "server"
	}
}

func run() error {
	if region == "" {
		log.Fatal("missing 'region' parameter")
		return errors.New("missing 'region' parameter")
	}
	if serverLabelValue == "" {
		log.Fatal("missing 'serverLabelValue' parameter")
		return errors.New("missing 'serverLabelValue' parameter")
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
		return err
	}
	client := ec2.NewFromConfig(cfg)

	ec2Util := pkg.EC2Util{
		Client: client,
	}

	instancesIDs, err := ec2Util.GetInstanceIDsStatusByLabel(serverLabelKey, serverLabelValue)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if len(instancesIDs) < 1 {
		log.Fatal("No matching instance found")
		return errors.New("No matching instance found")
	}

	startOneInstance := false
	startInstanceID := ""
	for true {
		for _, instanceID := range instancesIDs {
			instanceState, err := ec2Util.GetInstanceStatusByID(instanceID)
			if err != nil {
				log.Fatal(err)
				return err
			}
			if instanceState == pkg.InstanceStatusStopped {
				err := ec2Util.StartInstance(instanceID)
				if err != nil {
					log.Fatal(err)
					return err
				}
				startOneInstance = true
				startInstanceID = instanceID
				break
			}
		}
		if startOneInstance {
			break
		}
	}

	timeout := time.After(time.Duration(5) * time.Minute)
	isTimeout := false
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

StateWatch:
	for {
		select {
		case <-ticker.C:
			instanceState, _ := ec2Util.GetInstanceStatusByID(startInstanceID)
			if instanceState == pkg.InstanceStatusRunning {
				break StateWatch
			}
		case <-timeout:
			isTimeout = true
			break StateWatch
		}
	}

	if isTimeout {
		return errors.New("start instance timeout")
	}
	err = os.Setenv("INSTANCE_ID", startInstanceID)
	return err
}
