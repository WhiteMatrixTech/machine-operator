package stop

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"machine-operator/pkg"
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

	machineUtil, err := pkg.GetInstanceUtil(platform, region)
	if err != nil {
		log.Println("failed to get platform client")
		return errors.New(err.Error())
	}

	instanceState, err := machineUtil.StopInstance(instanceID)
	if err != nil {
		return err
	}
	fmt.Println("------------------stop the machine--------------------")
	fmt.Println(instanceState)
	return nil
}
