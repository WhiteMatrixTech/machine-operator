package pkg

type InstanceUtil interface {
	GetInstanceIDsByTag(tags map[string]string) ([]string, error)
	GetInstanceStatusByID(instanceID string) (string, error)
	StartInstance(instanceID string) error
	StopInstance(instanceID string) (string, error)
}
