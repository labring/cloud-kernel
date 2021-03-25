package _package

import (
	"errors"
	aliyunEcs "github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/sealyun/cloud-kernel/pkg/ecs"
	"github.com/sealyun/cloud-kernel/pkg/exit"
	"github.com/sealyun/cloud-kernel/pkg/logger"
	"github.com/sealyun/cloud-kernel/pkg/retry"
	"github.com/sealyun/cloud-kernel/pkg/sshcmd/sshutil"
	"github.com/sealyun/cloud-kernel/pkg/vars"
	"time"
)

type _package interface {
	InitK8sServer()
	WaitImages()
}

func Package(rt vars.Runtime, k8sVersion, runtimeVersion string) {
	instance := ecs.New(1, false, "")
	logger.Info("1. begin create ecs")
	var instanceInfo *aliyunEcs.DescribeInstanceAttributeResponse
	_ = retry.Do(func() error {
		var err error
		logger.Debug("1. retry fetch ecs info " + instance[0])
		instanceInfo, err = ecs.Describe(instance[0], "")
		if err != nil {
			return err
		}
		if len(instanceInfo.PublicIpAddress.IpAddress) == 0 {
			return errors.New("retry error")
		}
		if instanceInfo.Status != "Running" {
			return errors.New("retry error")
		}
		return nil
	}, 200, 500*time.Millisecond, false)
	publicIP := instanceInfo.PublicIpAddress.IpAddress[0]
	s := sshutil.SSH{
		User:     "root",
		Password: vars.EcsPassword,
		Timeout:  nil,
	}
	logger.Debug("2. connect ssh: " + publicIP)
	_ = retry.Do(func() error {
		var err error
		logger.Debug("2. retry test ecs ssh: " + publicIP)
		_, err = s.CmdAndError(publicIP, "ls /")
		if err != nil {
			return err
		} else {
			return nil
		}
	}, 200, 500*time.Millisecond, false)
	var k8s _package
	switch rt {
	case vars.Docker:
		k8s = NewDockerK8s(k8sVersion, runtimeVersion, publicIP)
	}
	if k8s == nil {
		exit.ProcessError(errors.New("k8s interface is nil"))
		return
	}
	logger.Info("3. install k8s[ " + k8sVersion + " ] : " + publicIP)
	k8s.InitK8sServer()
	logger.Info("4. wait k8s[ " + k8sVersion + " ] pull all images: " + publicIP)
	checkKubeStatus(4, publicIP, s, false)
	logger.Info("5. k8s[ " + k8sVersion + " ] is running: " + publicIP)
	s.CmdAsync(publicIP, "docker images")
	s.CmdAsync(publicIP, "kubectl get pod -n kube-system")
}
