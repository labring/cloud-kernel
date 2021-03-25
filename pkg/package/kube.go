package _package

import (
	"errors"
	"fmt"
	"github.com/sealyun/cloud-kernel/pkg/logger"
	"github.com/sealyun/cloud-kernel/pkg/retry"
	"github.com/sealyun/cloud-kernel/pkg/sshcmd/sshutil"
	"time"
)

func checkKubeStatus(step string, publicIP string, s sshutil.SSH, allRunning bool) error {
	return retry.Do(func() error {
		logger.Debug(fmt.Sprintf("%s. retry wait k8s  pod is running :%s", step, publicIP))
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
		if allRunning {
			checkShell = "kubectl  get pod -n kube-system  | grep  0/ | wc -l"
			zeroRunningNum := s.CmdToString(publicIP, checkShell, "")
			if zeroRunningNum != "0" {
				return errors.New("retry error")
			}
		}
		return nil
	}, 200, 500*time.Millisecond, true)
}
