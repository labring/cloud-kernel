package _package

import (
	"errors"
	"fmt"
	"github.com/sealyun/cloud-kernel/pkg/logger"
	"github.com/sealyun/cloud-kernel/pkg/retry"
	"github.com/sealyun/cloud-kernel/pkg/sshcmd/sshutil"
	"strings"
	"time"
)

func checkKubeStatus(step string, publicIP string, s sshutil.SSH, allRunning bool) error {
	defaultTimeout := time.Duration(10) * time.Second
	s.Timeout = &defaultTimeout
	return retry.Do(func() error {
		logger.Debug(fmt.Sprintf("%s. retry wait k8s  pod is running :%s", step, publicIP))
		checkShell := "kubectl  get pod -n kube-system   | grep -v \"RESTARTS\" | wc -l"
		podNum := strings.TrimSpace(s.CmdToString(publicIP, checkShell, ""))
		logger.Debug("当前pod的数量为: " + podNum)
		if podNum == "0" {
			return errors.New("retry error")
		}
		checkShell = "kubectl  get pod -n kube-system  | grep -v \"Running\" | grep -v \"RESTARTS\" | wc -l"
		notRunningNum := strings.TrimSpace(s.CmdToString(publicIP, checkShell, ""))
		logger.Debug("当前pod未Running状态的的数量为: " + notRunningNum)
		if notRunningNum != "" && notRunningNum != "0" {
			return errors.New("retry error")
		}
		checkShell = "kubectl  get pod -n kube-system  | grep  \"Running\" | grep calico-node | wc -l"
		calicoNodeNum := strings.TrimSpace(s.CmdToString(publicIP, checkShell, ""))
		logger.Debug("当前calico的Running状态的的数量为: " + calicoNodeNum)
		if calicoNodeNum == "0" {
			return errors.New("retry error")
		}
		if allRunning {
			checkShell = "kubectl  get pod -n kube-system  | grep  0/ | wc -l"
			zeroRunningNum := strings.TrimSpace(s.CmdToString(publicIP, checkShell, ""))
			if zeroRunningNum != "0" {
				return errors.New("retry error")
			}
		}
		return nil
	}, 50, 1*time.Second, false)
}
