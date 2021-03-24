package ecs

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/utils"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/sealyun/cloud-kernel/pkg/exit"
	"github.com/sealyun/cloud-kernel/pkg/logger"
	"github.com/sealyun/cloud-kernel/pkg/vars"
	"strconv"
	"sync"
)

var once sync.Once
var aliyun *ecs.Client

func getClient() *ecs.Client {
	once.Do(func() {
		var err error
		vars.LoadVars()
		aliyun, err = ecs.NewClientWithAccessKey("cn-hongkong", vars.AkId, vars.AkSK)
		if err != nil {
			exit.ProcessError(err)
		}
	})
	return aliyun
}

func New(amount int, dryRun bool) []string {
	client := getClient()
	// 创建请求并设置参数
	request := ecs.CreateRunInstancesRequest()
	request.Amount = requests.Integer(strconv.Itoa(amount))
	request.ImageId = "centos_7_04_64_20G_alibase_201701015.vhd"
	request.InstanceType = "ecs.c5.xlarge"
	request.InternetChargeType = "PayByTraffic"
	request.InternetMaxBandwidthIn = "50"
	request.InternetMaxBandwidthOut = "50"
	//request.KeyPairName = "release"
	request.InstanceChargeType = "PostPaid"
	request.SpotStrategy = "SpotAsPriceGo"
	request.RegionId = "cn-hongkong"
	request.SecurityGroupId = "sg-j6cb45dolegxcb32b47w"
	request.VSwitchId = "vsw-j6cvaap9o5a7et8uumqyx"
	request.ZoneId = "cn-hongkong-c"
	request.Password = "Fanux#123"
	request.ClientToken = utils.GetUUID()
	request.DryRun = requests.Boolean(strconv.FormatBool(dryRun))
	response, err := client.RunInstances(request)
	if err != nil {
		exit.ProcessError(err)
	}
	return response.InstanceIdSets.InstanceIdSet
}

func Delete(dryRun bool, instanceId []string) {
	client := getClient()
	// 创建请求并设置参数
	request := ecs.CreateDeleteInstancesRequest()
	request.DryRun = requests.Boolean(strconv.FormatBool(dryRun))
	request.Force = "true"
	request.RegionId = "cn-hongkong"
	request.InstanceId = &instanceId
	response, err := client.DeleteInstances(request)
	if err != nil {
		exit.ProcessError(err)
	}
	logger.Info("删除成功: %s", response.RequestId)
}

func Describe(instanceId string) (*ecs.DescribeInstanceAttributeResponse, error) {
	client := getClient()
	request := ecs.CreateDescribeInstanceAttributeRequest()
	request.InstanceId = instanceId
	return client.DescribeInstanceAttribute(request)
}
