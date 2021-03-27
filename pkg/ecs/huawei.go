package ecs

import (
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

func (a *HuaweiEcs) New(amount int, dryRun bool, bandwidthOut bool) []string {
	client := a.getClient()
	ids, err := client.RunInstances(amount, dryRun, bandwidthOut)
	if err != nil {
		_ = cutils.ProcessCloudError(err)
		return nil
	}
	return ids
}

func (a *HuaweiEcs) Delete(dryRun bool, instanceId []string) error {
	client := a.getClient()
	response, err := client.DeleteInstances(instanceId, true)
	if err != nil {
		_ = cutils.ProcessCloudError(err)
		logger.Error("递归删除ecs")
		return a.Delete(dryRun, instanceId)
	}
	logger.Info("删除成功: %s", *response.JobId)
	return nil
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
