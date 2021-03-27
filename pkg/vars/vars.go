package vars

import (
	"fmt"
)

var (
	DingDing       string
	AkId           string
	AkSK           string
	MarketCtlToken string
	IsArm64        bool

	SSHCmdDownload    = "https://github.com/cuisongliu/sshcmd/releases/download/v%s/sshcmd%s"                 //sshcmd-arm64
	SealosDownload    = "https://sealyun.oss-accelerate.aliyuncs.com/v%s/sealos%s"                            //sealos-arm64
	MarketCtlDownload = "https://sealyun-market.oss-accelerate.aliyuncs.com/marketctl/v%s/linux_%s/marketctl" //linux_arm64

	DockerShell = "wget https://download.docker.com/linux/static/stable/%s/docker-%s.tgz  " + //aarch64
		"-O docker.tgz && cp docker.tgz kube/docker/docker.tgz"

	ContainerdShell = "wget https://sealyun.oss-accelerate.aliyuncs.com/tools/cri-containerd-cni-%s-linux-%s.tar.gz " +
		"-O cri-containerd-cni-linux.tar.gz && " +
		"cp cri-containerd-cni-linux.tar.gz kube/containerd/cri-containerd-cni-linux.tar.gz"

	CrictlShell = "wget https://sealyun.oss-accelerate.aliyuncs.com/tools/crictl-v%s-linux-%s.tar.gz " +
		"-O  crictl.tar.gz &&  tar xf crictl.tar.gz && cp crictl kube/bin/"

	KubeShell = "wget https://dl.k8s.io/v%s/kubernetes-server-linux-%s.tar.gz -O  kubernetes-server.tar.gz && " + //arm64
		"tar zxvf kubernetes-server.tar.gz"
	KubeVersion      string
	DefaultPrice     float64
	DefaultZeroPrice float64
	DefaultClass     = "cloud_kernel" //cloud_kernel
	DefaultProduct   = "kubernetes"   //kubernetes
)

const (
	EcsPassword = "Fanux#123"

	defaultCriCtlVersion     = "1.20.0"
	defaultDockerVersion     = "19.03.12"
	defaultContainerdVersion = "1.4.3"

	defaultSealosVersion    = "3.3.9-rc.1"
	defaultMarketCtlVersion = "1.0.5" //v1.0.5
	defaultSSHCmdVersion    = "1.5.5"
)

func platform() map[string]map[bool]string {
	return map[string]map[bool]string{
		"sshcmd": {
			false: "",
			true:  "-arm64",
		},
		"sealos": {
			false: "",
			true:  "-arm64",
		},
		"marketctl": {
			false: "amd64",
			true:  "arm64",
		},
		"docker": {
			false: "x86_64",
			true:  "aarch64",
		},
		"containerd": {
			false: "amd64",
			true:  "arm64",
		},
		"crictl": {
			false: "amd64",
			true:  "arm64",
		},
		"kube": {
			false: "amd64",
			true:  "arm64",
		},
	}
}

func LoadVars() error {
	KubeShell = fmt.Sprintf(KubeShell, KubeVersion, platform()["kube"][IsArm64])
	//sshcmd
	SSHCmdDownload = fmt.Sprintf(SSHCmdDownload, defaultSSHCmdVersion, platform()["sshcmd"][IsArm64])
	//sealos
	SealosDownload = fmt.Sprintf(SealosDownload, defaultSealosVersion, platform()["sealos"][IsArm64])
	//marketctl
	MarketCtlDownload = fmt.Sprintf(MarketCtlDownload, defaultMarketCtlVersion, platform()["marketctl"][IsArm64])
	//containerd
	ContainerdShell = fmt.Sprintf(ContainerdShell, defaultContainerdVersion, platform()["containerd"][IsArm64])
	DockerShell = fmt.Sprintf(DockerShell, platform()["docker"][IsArm64], defaultDockerVersion)
	CrictlShell = fmt.Sprintf(CrictlShell, defaultCriCtlVersion, platform()["crictl"][IsArm64])
	if IsArm64 {
		DefaultProduct = "kubernetes-arm64"
	}
	return nil
}

const MarketYaml = `
market:
  body:
    spec:
      name: v%s
      price: %.2f
      product:
        class: %s
        productName: %s
      url: /tmp/kube%s.tar.gz
    status:
      productVersionStatus: ONLINE
  kind: productVersion`
