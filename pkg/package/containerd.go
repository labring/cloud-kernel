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
wget https://dl.k8s.io/v%s/kubernetes-server-linux-amd64.tar.gz && \
wget https://sealyun.oss-accelerate.aliyuncs.com/tools/cri-containerd-cni-%s-linux-amd64.tar.gz && \
wget https://sealyun.oss-accelerate.aliyuncs.com/tools/crictl-v%s-linux-amd64.tar.gz && \
cp  cri-containerd-cni-%s-linux-amd64.tar.gz kube/containerd/cri-containerd-cni-linux-amd64.tar.gz && \
tar zxvf kubernetes-server-linux-amd64.tar.gz && tar xf crictl-v%s-linux-amd64.tar.gz && \
cp  crictl kube/bin/ && \
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
	k8sVersion        string
	containerdVersion string
	ssh               sshutil.SSH
	publicIP          string
}

func NewContainerdK8s(k8sVersion, containerdVersion, publicIP string) _package {
	return &containerdK8s{
		k8sVersion:        k8sVersion,
		containerdVersion: containerdVersion,
		ssh: sshutil.SSH{
			User:     "root",
			Password: vars.EcsPassword,
			Timeout:  nil,
		},
		publicIP: publicIP,
	}
}
func (d *containerdK8s) InitK8sServer() error {
	if d.containerdVersion == "" {
		d.containerdVersion = vars.DefaultContainerdVersion
	}
	err := d.ssh.CmdAsync(d.publicIP,
		fmt.Sprintf(containerdShell, d.k8sVersion,
			d.containerdVersion, vars.CriCtlVersion, d.containerdVersion, vars.CriCtlVersion, d.k8sVersion))
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
	err := d.ssh.CmdAsync(d.publicIP, fmt.Sprintf(containerdSaveShell, d.k8sVersion, d.k8sVersion))
	if err != nil {
		return utils.ProcessError(err)
	}
	return nil
}