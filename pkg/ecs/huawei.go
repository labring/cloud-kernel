package ecs

import (
	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	"github.com/sealyun/cloud-kernel/pkg/ecs/huawei"
	"github.com/sealyun/cloud-kernel/pkg/logger"
	cutils "github.com/sealyun/cloud-kernel/pkg/utils"
	"github.com/sealyun/cloud-kernel/pkg/vars"
	"sync"
)

func (a *HuaweiEcs) getClient() *huawei.HClient {
	a.ecsOnce.Do(func() {
		a.ecsHKCli = huawei.NewClientWithAccessKey(vars.AkId, vars.AkSK)
	})
	return a.ecsHKCli
}

type HuaweiEcs struct {
	ecsOnce  sync.Once
	ecsHKCli *huawei.HClient
}

func (a *HuaweiEcs) Healthy() error {
	client := a.getClient()
	_, err := client.EcsClient.NovaListAvailabilityZones(&model.NovaListAvailabilityZonesRequest{})
	if err != nil {
		return errors.New("华为云 " + err.Error())
	}
	return nil
}
func (a *HuaweiEcs) New(amount int, dryRun bool, bandwidthOut bool) []string {
	client := a.getClient()
	ids, err := client.RunInstances(amount, dryRun, bandwidthOut)
	if err != nil {
		_ = cutils.ProcessCloudError(err)
		return nil
	}
	return ids
}

func (a *HuaweiEcs) Delete(instanceId []string, maxCount int) {
	client := a.getClient()
	var response *model.DeleteServersResponse
	var err error
	for i := 0; i < maxCount; i++ {
		logger.Error("循环删除ecs")
		response, err = client.DeleteInstances(instanceId, true)
		if err != nil {
			_ = cutils.ProcessCloudError(err)
		} else {
			break
		}
	}
	if err == nil {
		logger.Info("删除ecs成功: %s", *response.JobId)
	} else {
		logger.Error("删除ecs失败: %v", instanceId)
	}
}

func (a *HuaweiEcs) Describe(instanceId string) (*CloudInstanceResponse, error) {
	client := a.getClient()
	attr, err := client.Describe(instanceId)
	if err != nil {
		return nil, err
	}
	iResponse := &CloudInstanceResponse{
		IsOk: attr.Server.Status == "ACTIVE",
	}
	if attr.Server.Addresses != nil {
		for _, v := range attr.Server.Addresses {
			if len(v) != 0 {
				for _, vv := range v {
					if vv.OSEXTIPStype != nil {
						if *vv.OSEXTIPStype == model.GetServerAddressOSEXTIPStypeEnum().FIXED {
							iResponse.PrivateIP = vv.Addr
						}
						if *vv.OSEXTIPStype == model.GetServerAddressOSEXTIPStypeEnum().FLOATING {
							iResponse.PublicIP = vv.Addr
						}
					}
				}
			}
		}
	}
	return iResponse, nil
}
