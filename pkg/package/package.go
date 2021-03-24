package _package

import (
	"errors"
	"fmt"
	aliyunEcs "github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/rafaeljesus/retry-go"
	"github.com/sealyun/cloud-kernel/pkg/ecs"
	"github.com/sealyun/cloud-kernel/pkg/logger"
	"github.com/sealyun/cloud-kernel/pkg/sshcmd/sshutil"
	"time"
)

func Package() {
	instance := ecs.New(1, false)
	logger.Info("1. begin create ecs")
	var instanceInfo *aliyunEcs.DescribeInstanceAttributeResponse
	_ = retry.Do(func() error {
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
	}, 100, 500*time.Millisecond)
	publicIP := instanceInfo.PublicIpAddress.IpAddress[0]
	//publicIP:="8.210.82.137"
	s := sshutil.SSH{
		User:     "root",
		Password: "Fanux#123",
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
	}, 100, 500*time.Millisecond)
	k8s := "1.20.0"
	docker := "19.03.12"
	logger.Info("3. install k8s[ " + k8s + " ] : " + publicIP)
	s.CmdAsync(publicIP, fmt.Sprintf(docker_shell, k8s, docker, docker, k8s))
	logger.Info("4. wait k8s[ " + k8s + " ] pull all images: " + publicIP)
	retry.Do(func() error {
		logger.Debug("4. retry wait k8s all pod is running :" + publicIP)
		checkShell := "kubectl  get pod -n kube-system   | grep -v \"RESTARTS\" | wc -l"
		podNum := s.CmdToString(publicIP, checkShell, "")
		if podNum == "0" {
			return errors.New("retry error")
		}
		checkShell = "kubectl  get pod -n kube-system  | grep -v \"Running\" | grep -v \"RESTARTS\" | wc -l"
		notRunningNum := s.CmdToString(publicIP, checkShell, "")
		if notRunningNum != "0" {
			return errors.New("retry error")
		}
		checkShell = "kubectl  get pod -n kube-system  | grep  0/ | wc -l"
		zeroRunningNum := s.CmdToString(publicIP, checkShell, "")
		if zeroRunningNum != "0" {
			return errors.New("retry error")
		}
		return nil
	}, 100, 500*time.Millisecond)
	logger.Info("5. k8s[ " + k8s + " ] is running: " + publicIP)
	s.CmdAsync(publicIP, "kubectl get pod -n kube-system")
}
