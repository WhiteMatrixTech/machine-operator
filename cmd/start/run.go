package start

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"machine-operator/pkg"
)

var (
	platform       string
	region         string
	tags           string
	instanceIDPath string
	StartCmd       = &cobra.Command{
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
	StartCmd.PersistentFlags().StringVar(&platform,
		"platform", os.Getenv("platform"),
		"the platform")
	StartCmd.PersistentFlags().StringVar(&region,
		"region", os.Getenv("region"),
		"the region of cloud platform")
	StartCmd.PersistentFlags().StringVar(&tags,
		"tags", os.Getenv("tags"),
		"the instance tags")
	StartCmd.PersistentFlags().StringVar(&instanceIDPath,
		"instanceIDPath", os.Getenv("instanceIDPath"),
		"file path to write instanceID")
}

func preRun() {
	if instanceIDPath == "" {
		instanceIDPath = "instanceID.txt"
	}
}

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

	startOneInstance := false
	startedInstanceID := ""
	for true {
		for _, instanceID := range instancesIDs {
			instanceState, err := instanceUtil.GetInstanceStatusByID(instanceID)
			if err != nil {
				log.Fatal(err)
				return err
			}
			if instanceState == pkg.InstanceStatusStopped {
				err := instanceUtil.StartInstance(instanceID)
				if err != nil {
					log.Println(err)
					return err
				}
				startOneInstance = true
				startedInstanceID = instanceID
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
			instanceState, _ := instanceUtil.GetInstanceStatusByID(startedInstanceID)
			fmt.Println("the instance state: " + instanceState)
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
	fmt.Println("The instance started is " + startedInstanceID)
	err = pkg.WriteToFile(instanceIDPath, startedInstanceID)
	return err
}
