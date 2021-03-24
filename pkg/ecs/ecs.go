package ecs

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/utils"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/cuisongliu/cloud-kernel/pkg/exit"
	"github.com/cuisongliu/cloud-kernel/pkg/vars"
)

func New() {
	client, err := ecs.NewClientWithAccessKey("cn-hongkong", vars.AkId, vars.AkSK)
	if err != nil {
		exit.ProcessError(err)
	} else {
		// 创建请求并设置参数
		request := ecs.CreateRunInstancesRequest()
		request.ImageId = "centos_7_04_64_20G_alibase_201701015.vhd"
		request.InstanceType = "ecs.c5.xlarge"
		request.InternetChargeType = "PayByTraffic"
		request.InternetMaxBandwidthIn = "50"
		request.InternetMaxBandwidthOut = "50"
		request.KeyPairName = "KeyPairName"
		request.InstanceChargeType = "PostPaid"
		request.SpotStrategy = "SpotAsPriceGo"
		request.RegionId = "cn-hongkong"
		request.SecurityGroupId = "sg-j6cb45dolegxcb32b47w"
		request.VSwitchId = "vsw-j6cvaap9o5a7et8uumqyx"
		request.ZoneId = "cn-hongkong-c"
		request.InstanceName = "MyInstance"
		request.ClientToken = utils.GetUUID()
		request.DryRun = "true"
		response, err := client.RunInstances(request)
		if err != nil {
			exit.ProcessError(err)
		}
		print(response, err)
	}
}
