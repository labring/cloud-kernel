package _package

import (
	"fmt"
	"github.com/sealyun/cloud-kernel/pkg/sshcmd/sshutil"
	"github.com/sealyun/cloud-kernel/pkg/utils"
	"github.com/sealyun/cloud-kernel/pkg/vars"
)

func upload(publicIP, k8sVersion string) {
	s := sshutil.SSH{
		User:     "root",
		Password: vars.EcsPassword,
		Timeout:  nil,
	}
	if err := downloadBin(s, publicIP, vars.MarketCtlDownload, "marketctl"); err != nil {
		_ = utils.ProcessError(err)
		return
	}
	_, v := utils.GetMajorMinorInt(k8sVersion)
	price := 50
	if v == 0 {
		price = 0.01
	}
	yaml := fmt.Sprintf(vars.MarketYaml, k8sVersion, price, k8sVersion)
	_ = s.CmdAsync(publicIP, "echo \""+yaml+"\" > /tmp/marketctl_"+k8sVersion+".yaml")
	_ = s.CmdAsync(publicIP, "cat /tmp/marketctl_"+k8sVersion+".yaml")
	//marketctl apply -f /tmp/marketctl_%s.yaml --domain https://www.sealyun.com --token %s --dd-token %s
}
