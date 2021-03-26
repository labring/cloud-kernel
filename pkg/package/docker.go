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

//k8s dockerShell k8s
var dockerShell = `yum install -y git conntrack && \
git clone https://github.com/sealyun/cloud-kernel && \
cd cloud-kernel && mkdir -p kube && cp -rf runtime/docker/* kube/ && \
%s && \
%s && \
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
	ssh      sshutil.SSH
	publicIP string
}

func NewDockerK8s(publicIP string) _package {

	return &dockerK8s{
		ssh: sshutil.SSH{
			User:     "root",
			Password: vars.EcsPassword,
			Timeout:  nil,
		},
		publicIP: publicIP,
	}
}
func (d *dockerK8s) InitK8sServer() error {
	err := d.ssh.CmdAsync(d.publicIP, fmt.Sprintf(dockerShell, vars.KubeShell, vars.DockerShell, vars.KubeVersion))
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
	err := d.ssh.CmdAsync(d.publicIP, fmt.Sprintf(dockerSaveShell, vars.KubeVersion, vars.KubeVersion))
	if err != nil {
		return utils.ProcessError(err)
	}
	return nil
}
