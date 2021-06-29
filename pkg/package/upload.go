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
	var price = vars.DefaultPrice
	if v == 0 {
		price = vars.DefaultZeroPrice
	}
	yaml := fmt.Sprintf(vars.MarketYaml, k8sVersion, price, vars.DefaultClass, vars.DefaultProduct, k8sVersion)
	_ = s.CmdAsync(publicIP, "echo \""+yaml+"\" > /tmp/marketctl_"+k8sVersion+".yaml")
	_ = s.CmdAsync(publicIP, "cat /tmp/marketctl_"+k8sVersion+".yaml")
	//marketctl apply -f /tmp/marketctl_%s.yaml --domain https://www.sealyun.com --token %s --dd-token %s
	marketCMD := fmt.Sprintf("marketctl apply -f /tmp/marketctl_%s.yaml --ci --token %s",
		k8sVersion, vars.MarketCtlToken)
	if vars.DingDing != "" {
		marketCMD = marketCMD + " --dd-token " + vars.DingDing
	}
	_ = s.CmdAsync(publicIP, marketCMD)
}

func uploadOSS(publicIP, k8sVersion string) {
	s := sshutil.SSH{
		User:     "root",
		Password: vars.EcsPassword,
		Timeout:  nil,
	}
	if err := downloadBin(s, publicIP, vars.OSSUtilDownload, "ossutil"); err != nil {
		_ = utils.ProcessError(err)
		return
	}
	writeOSSConfig := `cd cloud-kernel  && echo '%s' > oss-config && ossutil -c oss-config cp -f  /tmp/kube%s.tar.gz oss://sealyun-test/cloud_kernel/kube%s.tar.gz`
	ossConfig := &OSSConfig{
		KeyID:     vars.AkID,
		KeySecret: vars.AkSK,
	}
	writeShell := fmt.Sprintf(writeOSSConfig, ossConfig.TemplateConvert(), k8sVersion, k8sVersion)
	err := s.CmdAsync(publicIP, writeShell)
	if err != nil {
		_ = utils.ProcessError(err)
		return
	}
}
