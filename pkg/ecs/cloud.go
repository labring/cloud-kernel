package ecs

import "github.com/sealyun/cloud-kernel/pkg/vars"

type cloud interface {
	New(amount int, dryRun bool, bandwidthOut bool) []string
	Delete(dryRun bool, instanceId []string) error
	Describe(instanceId string) (*CloudInstanceResponse, error)
}
type CloudInstanceResponse struct {
	IsOk      bool
	PrivateIP string
	PublicIP  string
}

func NewCloud() cloud {
	var c cloud
	if vars.IsAmd64 {
		c = &AliyunEcs{}
	} else {
		c = &HuaweiEcs{}
	}
	return c
}
