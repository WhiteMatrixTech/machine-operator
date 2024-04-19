package stop

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"
	"log"
	"machine-operator/pkg"
	"os"
)

var (
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
	StartCmd.PersistentFlags().StringVar(&instanceID,
		"instanceID", os.Getenv("instanceID"),
		"id of machine")
	StartCmd.PersistentFlags().StringVar(&region,
		"region", os.Getenv("region"),
		"the region of aws config")
}

func run() error {
	if region == "" {
		log.Fatal("missing 'region' parameter")
		return errors.New("missing 'region' parameter")
	}
	if instanceID == "" {
		log.Fatal("missing 'instanceID' parameter")
		return errors.New("missing 'instanceID' parameter")
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

	instanceState, err := ec2Util.StopInstance(instanceID)
	return err
	fmt.Println("------------------stop the machine--------------------")
	fmt.Println(instanceState)
	return err
}
