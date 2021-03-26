package ecs

type cloud interface {
	New(amount int, dryRun bool, bandwidthOut bool) []string
	Delete(dryRun bool, instanceId []string) error
	Describe(instanceId string) (*CloudInstanceResponse, error)
}
type CloudInstanceResponse struct {
	Status    string
	PrivateIP string
	PublicIP  string
}
