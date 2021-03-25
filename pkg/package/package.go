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
	InitK8sServer() error
	WaitImages() error
	SavePackage() error
}

func Package(rt vars.Runtime, k8sVersion, runtimeVersion string) {
	instance := ecs.New(1, false, "", true)
	logger.Info("1. begin create ecs")
	var instanceInfo *aliyunEcs.DescribeInstanceAttributeResponse
	defer func() {
		_ = ecs.Delete(false, instance, "")
	}()
	var err error
	if err = retry.Do(func() error {
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
	}, 100, 1*time.Second, false); err != nil {
		_ = exit.ProcessError(err)
		return
	}
	publicIP := instanceInfo.PublicIpAddress.IpAddress[0]
	s := sshutil.SSH{
		User:     "root",
		Password: vars.EcsPassword,
		Timeout:  nil,
	}
	logger.Debug("2. connect ssh: " + publicIP)
	if err = retry.Do(func() error {
		var err error
		logger.Debug("2. retry test ecs ssh: " + publicIP)
		_, err = s.CmdAndError(publicIP, "ls /")
		if err != nil {
			return err
		} else {
			return nil
		}
	}, 100, 500*time.Millisecond, true); err != nil {
		_ = exit.ProcessError(err)
		return
	}
	var k8s _package
	switch rt {
	case vars.Docker:
		k8s = NewDockerK8s(k8sVersion, runtimeVersion, publicIP)
	}
	if k8s == nil {
		_ = exit.ProcessError(errors.New("k8s interface is nil"))
		return
	}
	logger.Info("3. install k8s[ " + k8sVersion + " ] : " + publicIP)
	if err = k8s.InitK8sServer(); err != nil {
		_ = exit.ProcessError(err)
		return
	}
	logger.Info("4. wait k8s[ " + k8sVersion + " ] pull all images: " + publicIP)
	if err = checkKubeStatus("4", publicIP, s, false); err != nil {
		_ = exit.ProcessError(err)
		return
	}
	if err = s.CmdAsync(publicIP, "docker images"); err != nil {
		_ = exit.ProcessError(err)
		return
	}
	if err = s.CmdAsync(publicIP, "kubectl get pod -n kube-system"); err != nil {
		_ = exit.ProcessError(err)
		return
	}
	if err = k8s.WaitImages(); err != nil {
		_ = exit.ProcessError(err)
		return
	}
	logger.Info("5. k8s[ " + k8sVersion + " ] image save: " + publicIP)
	if err = k8s.SavePackage(); err != nil {
		_ = exit.ProcessError(err)
		return
	}
	logger.Info("6. k8s[ " + k8sVersion + " ] testing: " + publicIP)
	test(publicIP, k8sVersion)
}
