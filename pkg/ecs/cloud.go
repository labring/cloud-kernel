package ecs

import "github.com/labring/cloud-kernel/pkg/vars"

type cloud interface {
	New(amount int, dryRun bool, bandwidthOut bool) []string
	Delete(instanceId []string, maxCount int)
	Describe(instanceId string) (*CloudInstanceResponse, error)
	Healthy() error
}
type CloudInstanceResponse struct {
	IsOk      bool
	PrivateIP string
	PublicIP  string
}

func NewCloud() cloud {
	var c cloud
	if !vars.IsArm64 {
		c = &AliyunEcs{}
	} else {
		c = &HuaweiEcs{}
	}
	return c
}
