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

//k8s docker docker k8s
var dockerShell = `yum install -y git conntrack && \
git clone https://github.com/sealyun/cloud-kernel && \
cd cloud-kernel && mkdir -p kube && cp -rf runtime/docker/* kube/ && \
wget https://dl.k8s.io/v%s/kubernetes-server-linux-amd64.tar.gz && \
wget https://download.docker.com/linux/static/stable/x86_64/docker-%s.tgz && \
cp  docker-%s.tgz kube/docker/docker.tgz && \
tar zxvf kubernetes-server-linux-amd64.tar.gz && \
cp  kubernetes/server/bin/kubectl kube/bin/ && \
cp  kubernetes/server/bin/kubelet kube/bin/ && \
cp  kubernetes/server/bin/kubeadm kube/bin/ && \
sed s/k8s_version/%s/g -i kube/conf/kubeadm.yaml && \
cd kube/shell && chmod a+x docker.sh && sh docker.sh && \
rm -rf /etc/docker/daemon.json && systemctl restart docker && \
sh init.sh && sh master.sh && \
docker pull fanux/lvscare &&  \
cp /usr/sbin/conntrack ../bin/`

var dockerSaveShell = `cd cloud-kernel &&  \
docker save -o images.tar  ` + "`docker images|grep ago|awk '{print $1\":\"$2}'` " + `&& \
mv images.tar kube/images/ && \
tar zcvf kube%s.tar.gz kube && mv kube%s.tar.gz /tmp/`

type dockerK8s struct {
	k8sVersion    string
	dockerVersion string
	ssh           sshutil.SSH
	publicIP      string
}

func NewDockerK8s(k8sVersion, dockerVersion, publicIP string) _package {
	return &dockerK8s{
		k8sVersion:    k8sVersion,
		dockerVersion: dockerVersion,
		ssh: sshutil.SSH{
			User:     "root",
			Password: vars.EcsPassword,
			Timeout:  nil,
		},
		publicIP: publicIP,
	}
}
func (d *dockerK8s) InitK8sServer() error {
	if d.dockerVersion == "" {
		d.dockerVersion = vars.DefaultDockerVersion
	}
	err := d.ssh.CmdAsync(d.publicIP, fmt.Sprintf(dockerShell, d.k8sVersion, d.dockerVersion, d.dockerVersion, d.k8sVersion))
	if err != nil {
		return utils.ProcessError(err)
	}
	return nil
}

func (d *dockerK8s) WaitImages() error {
	if err := d.ssh.CmdAsync(d.publicIP, "docker images"); err != nil {
		_ = utils.ProcessError(err)
		return err
	}
	err := retry.Do(func() error {
		logger.Debug(fmt.Sprintf("%d. retry wait k8s  pod is running :%s", 4, d.publicIP))
		checkShell := "docker images   | grep  \"lvscare\" | wc -l"
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

func (d *dockerK8s) SavePackage() error {
	err := d.ssh.CmdAsync(d.publicIP, fmt.Sprintf(dockerSaveShell, d.k8sVersion, d.k8sVersion))
	if err != nil {
		return utils.ProcessError(err)
	}
	return nil
}
