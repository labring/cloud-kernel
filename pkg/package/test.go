package _package

import (
	"errors"
	"fmt"
	aliyunEcs "github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/sealyun/cloud-kernel/pkg/ecs"
	"github.com/sealyun/cloud-kernel/pkg/logger"
	"github.com/sealyun/cloud-kernel/pkg/retry"
	"github.com/sealyun/cloud-kernel/pkg/sshcmd/sshutil"
	"github.com/sealyun/cloud-kernel/pkg/utils"
	"github.com/sealyun/cloud-kernel/pkg/vars"
	"strings"
	"time"
)

func test(publicIP, k8sVersion string) {
	master0 := ecs.New(1, false, "", true)
	others := ecs.New(3, false, "", false)
	instance := append(master0, others...)
	instanceInfos := make([]*aliyunEcs.DescribeInstanceAttributeResponse, len(instance))
	logger.Info("test1. begin create ecs")
	defer func() {
		_ = ecs.Delete(false, instance, "")
	}()
	var err error
	if err = retry.Do(func() error {
		var err error
		logger.Debug("test1. retry fetch ecs info " + strings.Join(instance, ","))
		for i, v := range instance {
			instanceInfos[i], err = ecs.Describe(v, "")
			if err != nil {
				return err
			}
			//master0 publicIP
			if i == 0 {
				if len(instanceInfos[i].PublicIpAddress.IpAddress) == 0 {
					return errors.New("retry error")
				}
			}
			if len(instanceInfos[i].VpcAttributes.PrivateIpAddress.IpAddress) == 0 {
				return errors.New("retry error")
			}
			if instanceInfos[i].Status != "Running" {
				return errors.New("retry error")
			}
		}
		return nil
	}, 100, 1*time.Second, false); err != nil {
		_ = utils.ProcessError(err)
		return
	}
	privateIPs := make([]string, len(instanceInfos))
	for i, v := range instanceInfos {
		privateIPs[i] = v.VpcAttributes.PrivateIpAddress.IpAddress[0]
	}
	s := sshutil.SSH{
		User:     "root",
		Password: vars.EcsPassword,
		Timeout:  nil,
	}
	if err = downloadBin(s, publicIP, vars.SSHCmdDownload, "sshcmd"); err != nil {
		_ = utils.ProcessError(err)
		return
	}
	if err = downloadBin(s, publicIP, vars.SealosDownload, "sealos"); err != nil {
		_ = utils.ProcessError(err)
		return
	}
	cmd := "sshcmd --user root --passwd %s --host %s --host %s --host %s --host %s --cmd \"ls -l\""
	connShell := fmt.Sprintf(cmd, vars.EcsPassword, privateIPs[0], privateIPs[1], privateIPs[2], privateIPs[3])
	logger.Debug("test2. connect ssh ")
	if err = retry.Do(func() error {
		var err error
		logger.Debug("test2. retry test ecs ssh ")
		_, err = s.CmdAndError(publicIP, connShell)
		if err != nil {
			return err
		} else {
			return nil
		}
	}, 100, 500*time.Millisecond, true); err != nil {
		_ = utils.ProcessError(err)
		return
	}
	logger.Debug("test2. install  k8s ( 3 master 1 node ) ")
	cmd = "sealos init --master %s --master %s --master %s --node %s --passwd %s --version v%s --pkg-url /tmp/kube%s.tar.gz"
	installCmd := fmt.Sprintf(cmd, privateIPs[0], privateIPs[1], privateIPs[2], privateIPs[3], vars.EcsPassword, k8sVersion, k8sVersion)
	if err = s.CmdAsync(publicIP, installCmd); err != nil {
		_ = utils.ProcessError(err)
		return
	}
	logger.Debug("test3. check  k8s ( 3 master 1 node ) status")
	master0IP := instanceInfos[0].PublicIpAddress.IpAddress[0]
	if err = checkKubeStatus("test3", master0IP, s, true); err != nil {
		_ = utils.ProcessError(err)
		return
	}
	defer func() {
		_ = s.CmdAsync(master0IP, "kubectl get pod -n kube-system")
	}()
}

func downloadBin(s sshutil.SSH, publicIP, url, name string) error {
	if err := s.CmdAsync(publicIP, fmt.Sprintf("wget %s && chmod +x %s && mv %s /usr/bin", url, name, name)); err != nil {
		return err
	}
	return nil
}
