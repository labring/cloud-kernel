package _package

import (
	"errors"
	aliyunEcs "github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/sealyun/cloud-kernel/pkg/ecs"
	"github.com/sealyun/cloud-kernel/pkg/logger"
	"github.com/sealyun/cloud-kernel/pkg/retry"
	"github.com/sealyun/cloud-kernel/pkg/sshcmd/sshutil"
	"github.com/sealyun/cloud-kernel/pkg/utils"
	"github.com/sealyun/cloud-kernel/pkg/vars"
	"time"
)

type _package interface {
	InitK8sServer() error
	WaitImages() error
	SavePackage() error
}

func Package(k8sVersion string) {
	vars.KubeVersion = k8sVersion
	err := vars.LoadVars()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	instance := ecs.New(1, false, true)
	logger.Info("1. begin create ecs")
	var instanceInfo *aliyunEcs.DescribeInstanceAttributeResponse
	defer func() {
		_ = ecs.Delete(false, instance)
	}()
	if err = retry.Do(func() error {
		var err error
		logger.Debug("1. retry fetch ecs info " + instance[0])
		instanceInfo, err = ecs.Describe(instance[0])
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
		_ = utils.ProcessError(err)
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
	}, 20, 500*time.Millisecond, true); err != nil {
		_ = utils.ProcessError(err)
		return
	}
	var k8s _package
	if utils.For120(k8sVersion) {
		k8s = NewContainerdK8s(publicIP)
	} else {
		k8s = NewDockerK8s(publicIP)
	}
	if k8s == nil {
		_ = utils.ProcessError(errors.New("k8s interface is nil"))
		return
	}
	logger.Info("3. install k8s[ " + k8sVersion + " ] : " + publicIP)
	if err = k8s.InitK8sServer(); err != nil {
		_ = utils.ProcessError(err)
		return
	}
	logger.Info("4. wait k8s[ " + k8sVersion + " ] pull all images: " + publicIP)
	if err = checkKubeStatus("4", publicIP, s, false); err != nil {
		_ = utils.ProcessError(err)
		return
	}
	if err = s.CmdAsync(publicIP, "kubectl get pod -n kube-system"); err != nil {
		_ = utils.ProcessError(err)
		return
	}
	if err = k8s.WaitImages(); err != nil {
		_ = utils.ProcessError(err)
		return
	}
	logger.Info("5. k8s[ " + k8sVersion + " ] image save: " + publicIP)
	if err = k8s.SavePackage(); err != nil {
		_ = utils.ProcessError(err)
		return
	}
	logger.Info("6. k8s[ " + k8sVersion + " ] testing: " + publicIP)
	if err = test(publicIP, k8sVersion); err != nil {
		_ = utils.ProcessError(err)
		return
	} else {
		logger.Info("6. k8s[ " + k8sVersion + " ] uploading: " + publicIP)
		upload(publicIP, k8sVersion)
	}
	logger.Info("7. k8s[ " + k8sVersion + " ] finished. " + publicIP)
}
