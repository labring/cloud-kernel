package _package

import (
	"errors"
	"time"

	"github.com/sealyun/cloud-kernel/pkg/ecs"
	"github.com/sealyun/cloud-kernel/pkg/logger"
	"github.com/sealyun/cloud-kernel/pkg/retry"
	"github.com/sealyun/cloud-kernel/pkg/sshcmd/sshutil"
	"github.com/sealyun/cloud-kernel/pkg/utils"
	"github.com/sealyun/cloud-kernel/pkg/vars"
)

type _package interface {
	InitK8sServer() error
	WaitImages() error
	SavePackage() error
}

func Package(k8sVersion string) error {
	vars.KubeVersion = k8sVersion
	err := vars.LoadVars()
	if err != nil {
		return err
	}
	instance := ecs.NewCloud().New(1, false, true)
	if instance == nil {
		return errors.New("create ecs is error")
	}
	logger.Info("1. begin create ecs")
	var instanceInfo *ecs.CloudInstanceResponse
	defer func() {
		ecs.NewCloud().Delete(instance, 10)
	}()
	if err = retry.Do(func() error {
		var err error
		logger.Debug("1. retry fetch ecs info " + instance[0])
		instanceInfo, err = ecs.NewCloud().Describe(instance[0])
		if err != nil {
			return err
		}
		if instanceInfo.PublicIP == "" {
			return errors.New("retry error")
		}
		if !instanceInfo.IsOk {
			return errors.New("retry error")
		}
		return nil
	}, 100, 1*time.Second, false); err != nil {
		return utils.ProcessError(err)
	}
	publicIP := instanceInfo.PublicIP
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
		return utils.ProcessError(err)
	}
	var k8s _package
	if utils.For120(k8sVersion) {
		k8s = NewContainerdK8s(publicIP)
	} else {
		k8s = NewDockerK8s(publicIP)
	}
	if k8s == nil {
		return utils.ProcessError(errors.New("k8s interface is nil"))
	}
	logger.Info("3. install k8s[ " + k8sVersion + " ] : " + publicIP)
	if err = k8s.InitK8sServer(); err != nil {
		return utils.ProcessError(err)
	}
	logger.Info("4. wait k8s[ " + k8sVersion + " ] pull all images: " + publicIP)
	if err = checkKubeStatus("4", publicIP, s, false); err != nil {
		return utils.ProcessError(err)
	}
	if err = s.CmdAsync(publicIP, "kubectl get pod -n kube-system"); err != nil {
		return utils.ProcessError(err)
	}
	if err = k8s.WaitImages(); err != nil {
		return utils.ProcessError(err)
	}
	logger.Info("5. k8s[ " + k8sVersion + " ] image save: " + publicIP)
	if err = k8s.SavePackage(); err != nil {
		return utils.ProcessError(err)
	}
	if vars.Testing {
		logger.Info("6. k8s[ " + k8sVersion + " ] testing: " + publicIP)
		if err = test(publicIP, k8sVersion); err != nil {
			return utils.ProcessError(err)
		}
	} else {
		logger.Info("6. k8s[ " + k8sVersion + " ] skip testing: " + publicIP)
	}
	if vars.Uploading {
		logger.Info("7. k8s[ " + k8sVersion + " ] uploading: " + publicIP)
		upload(publicIP, k8sVersion)
	} else {
		logger.Info("7. k8s[ " + k8sVersion + " ] uploading test oss: " + publicIP)
		uploadOSS(publicIP, k8sVersion)
	}
	logger.Info("8. k8s[ " + k8sVersion + " ] finished. " + publicIP)
	return nil
}
