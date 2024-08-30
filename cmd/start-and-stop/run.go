package start_and_stop

import (
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"machine-operator/pkg"
)

var (
	platform string
	region   string
	tags     string
	StartCmd = &cobra.Command{
		Use:          "start-and-stop",
		Short:        "start and stop machine",
		Example:      "machine-operator start-and-stop",
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
	StartCmd.PersistentFlags().StringVar(&platform,
		"platform", os.Getenv("platform"),
		"the platform")
	StartCmd.PersistentFlags().StringVar(&region,
		"region", os.Getenv("region"),
		"the region of cloud platform")
	StartCmd.PersistentFlags().StringVar(&tags,
		"tags", os.Getenv("tags"),
		"the instance tags")
}

func preRun() {}

func run() error {
	if platform == "" {
		log.Println("missing 'platform' parameter")
		return errors.New("missing platform parameter")
	}
	if region == "" {
		log.Println("missing 'region' parameter")
		return errors.New("missing 'region' parameter")
	}

	if tags == "" {
		log.Println("missing 'tags' parameter")
		return errors.New("missing 'tags' parameter")
	}

	instanceUtil, err := pkg.GetInstanceUtil(platform, region)
	if err != nil {
		log.Println("failed to get platform client")
		return errors.New(err.Error())
	}

	instanceTags := make(map[string]string)
	tagsSlice := strings.Split(tags, ",")
	for _, tag := range tagsSlice {
		tagSlice := strings.Split(tag, ":")
		instanceTags[tagSlice[0]] = tagSlice[1]
	}

	instancesIDs, err := instanceUtil.GetInstanceIDsByTag(instanceTags)
	if err != nil {
		log.Println(err)
		return err
	}
	if len(instancesIDs) < 1 {
		log.Println("No matching instance found")
		return errors.New("No matching instance found ")
	}

	instanceNeedsStop := make([]string, 0)

	for i := range instancesIDs {
		instanceState, err := instanceUtil.GetInstanceStatusByID(instancesIDs[i])
		if err != nil {
			log.Println(err.Error())
			continue
		}
		if instanceState == pkg.InstanceStatusStopped {
			err := instanceUtil.StartInstance(instancesIDs[i])
			if err != nil {
				log.Println(err.Error())
				continue
			}
			instanceNeedsStop = append(instanceNeedsStop, instancesIDs[i])
		}

	}

	time.Sleep(2 * time.Minute)
	for i := range instanceNeedsStop {
		instanceState, _ := instanceUtil.GetInstanceStatusByID(instanceNeedsStop[i])
		if instanceState == pkg.InstanceStatusRunning {
			_, err := instanceUtil.StopInstance(instanceNeedsStop[i])
			if err != nil {
				log.Println(err.Error())
				continue
			}
		}
	}

	return nil
}
