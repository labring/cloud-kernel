package _package

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sealyun/cloud-kernel/pkg/ecs"
	"github.com/sealyun/cloud-kernel/pkg/logger"
	"github.com/sealyun/cloud-kernel/pkg/retry"
	"github.com/sealyun/cloud-kernel/pkg/sshcmd/sshutil"
	"github.com/sealyun/cloud-kernel/pkg/vars"
)

func test(publicIP, k8sVersion string) error {
	master0 := ecs.NewCloud().New(1, false, true)
	others := ecs.NewCloud().New(3, false, false)
	instance := append(master0, others...)
	instanceInfos := make([]*ecs.CloudInstanceResponse, len(instance))
	logger.Info("test1. begin create ecs")
	defer func() {
		ecs.NewCloud().Delete(instance, 10)
	}()
	var err error
	if err = retry.Do(func() error {
		var err error
		logger.Debug("test1. retry fetch ecs info " + strings.Join(instance, ","))
		for i, v := range instance {
			instanceInfos[i], err = ecs.NewCloud().Describe(v)
			if err != nil {
				return err
			}
			//master0 publicIP
			if i == 0 {
				if instanceInfos[i].PublicIP == "" {
					return errors.New("retry error")
				}
			}
			if instanceInfos[i].PrivateIP == "" {
				return errors.New("retry error")
			}
			if !instanceInfos[i].IsOk {
				return errors.New("retry error")
			}
		}
		return nil
	}, 100, 1*time.Second, false); err != nil {
		return err
	}
	privateIPs := make([]string, len(instanceInfos))
	for i, v := range instanceInfos {
		privateIPs[i] = v.PrivateIP
	}
	s := sshutil.SSH{
		User:     "root",
		Password: vars.EcsPassword,
		Timeout:  nil,
	}
	if err = downloadBin(s, publicIP, vars.SSHCmdDownload, "sshcmd"); err != nil {
		return err
	}
	if err = downloadBin(s, publicIP, vars.SealosDownload, "sealos"); err != nil {
		return err
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
		return err
	}
	logger.Debug("test2. install  k8s ( 3 master 1 node ) ")
	cmd = "sealos init --master %s --master %s --master %s --node %s --passwd %s --version v%s --pkg-url /tmp/kube%s.tar.gz"
	installCmd := fmt.Sprintf(cmd, privateIPs[0], privateIPs[1], privateIPs[2], privateIPs[3], vars.EcsPassword, k8sVersion, k8sVersion)
	if err = s.CmdAsync(publicIP, installCmd); err != nil {
		return err
	}
	logger.Debug("test3. check  k8s ( 3 master 1 node ) status")
	master0IP := instanceInfos[0].PublicIP
	if err = checkKubeStatus("test3", master0IP, s, true); err != nil {
		return err
	}
	checkShell := "kubectl  get nodes  | grep master | wc -l"
	masterNum := strings.TrimSpace(s.CmdToString(master0IP, checkShell, ""))
	if masterNum != "3" {
		return errors.New("当前集群master节点不为3")
	}
	checkShell = "kubectl  get nodes  | grep \"<none>\" | wc -l"
	nodeNum := strings.TrimSpace(s.CmdToString(master0IP, checkShell, ""))
	if nodeNum != "1" {
		return errors.New("当前集群node节点不为1")
	}
	_ = s.CmdAsync(master0IP, "kubectl get pod -n kube-system")
	_ = s.CmdAsync(master0IP, "kubectl get nodes -owide")
	return nil
}

func downloadBin(s sshutil.SSH, publicIP, url, name string) error {
	if err := s.CmdAsync(publicIP, fmt.Sprintf("wget %s -O %s && chmod +x %s && mv %s /usr/bin", url, name, name, name)); err != nil {
		return err
	}
	return nil
}
