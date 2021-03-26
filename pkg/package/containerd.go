package _package

import (
	"errors"
	"fmt"
	"github.com/sealyun/cloud-kernel/pkg/logger"
	"github.com/sealyun/cloud-kernel/pkg/retry"
	"github.com/sealyun/cloud-kernel/pkg/sshcmd/sshutil"
	"github.com/sealyun/cloud-kernel/pkg/utils"
	"github.com/sealyun/cloud-kernel/pkg/vars"
	"time"
)

//k8s containerd containerd k8s
var containerdShell = `yum install -y git conntrack && \
git clone https://github.com/sealyun/cloud-kernel && \
cd cloud-kernel && mkdir -p kube && cp -rf runtime/containerd/* kube/ && \
%s && \
%s && \
%s && \
cp  kubernetes/server/bin/kubectl kube/bin/ && \
cp  kubernetes/server/bin/kubelet kube/bin/ && \
cp  kubernetes/server/bin/kubeadm kube/bin/ && \
sed s/k8s_version/%s/g -i kube/conf/kubeadm.yaml && \
cd kube/shell && chmod a+x containerd.sh && sh containerd.sh && \
systemctl restart containerd && \
sh init.sh && sh master.sh && \
ctr -n=k8s.io images pull docker.io/fanux/lvscare:latest && \
cp /usr/sbin/conntrack ../bin/ && \
cp /usr/lib64/libseccomp* ../lib64/`

var containerdSaveShell = `cd cloud-kernel &&  \
ctr -n=k8s.io  images export images.tar $(ctr -n=k8s.io image ls  | awk '{print $1}' | grep -v sha256  | grep -v REF) && \
mv images.tar kube/images/ && \
tar zcvf kube%s.tar.gz kube && mv kube%s.tar.gz /tmp/`

type containerdK8s struct {
	ssh      sshutil.SSH
	publicIP string
}

func NewContainerdK8s(publicIP string) _package {
	return &containerdK8s{
		ssh: sshutil.SSH{
			User:     "root",
			Password: vars.EcsPassword,
			Timeout:  nil,
		},
		publicIP: publicIP,
	}
}
func (d *containerdK8s) InitK8sServer() error {
	err := d.ssh.CmdAsync(d.publicIP,
		fmt.Sprintf(containerdShell, vars.KubeShell,
			vars.ContainerdShell, vars.CrictlShell, vars.KubeVersion))
	if err != nil {
		return utils.ProcessError(err)
	}
	return nil
}

func (d *containerdK8s) WaitImages() error {
	if err := d.ssh.CmdAsync(d.publicIP, "crictl images"); err != nil {
		_ = utils.ProcessError(err)
		return err
	}
	err := retry.Do(func() error {
		logger.Debug(fmt.Sprintf("%d. retry wait k8s  pod is running :%s", 4, d.publicIP))
		checkShell := "crictl images   | grep  \"lvscare\" | wc -l"
		podNum := d.ssh.CmdToString(d.publicIP, checkShell, "")
		if podNum == "0" {
			return errors.New("retry error")
		}
		return nil
	}, 100, 500*time.Millisecond, false)
	if err != nil {
		return utils.ProcessError(err)
	}
	return nil
}

func (d *containerdK8s) SavePackage() error {
	err := d.ssh.CmdAsync(d.publicIP, fmt.Sprintf(containerdSaveShell, vars.KubeVersion, vars.KubeVersion))
	if err != nil {
		return utils.ProcessError(err)
	}
	return nil
}
