package vars

import (
	"errors"
	"fmt"
	"github.com/sealyun/cloud-kernel/pkg/sshcmd/cmd"
	"os"
	"strings"
)

var (
	DingDing       string
	AkId           string
	AkSK           string
	MarketCtlToken string
	KubeVersion    string

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

func IsAmd64() bool {
	arch := strings.ReplaceAll(cmd.CmdToString("arch"), "\n", "")
	amd64 := true
	switch arch {
	case "x86_64":
	case "i386":
		amd64 = true
	case "arm64":
	case "aarch64":
		amd64 = false
	default:
		amd64 = true
	}
	return amd64
}

func platform() map[string]map[bool]string {
	return map[string]map[bool]string{
		"sshcmd": {
			true:  "",
			false: "-arm64",
		},
		"sealos": {
			true:  "",
			false: "-arm64",
		},
		"marketctl": {
			true:  "amd64",
			false: "arm64",
		},
		"docker": {
			true:  "x86_64",
			false: "aarch64",
		},
		"containerd": {
			true:  "amd64",
			false: "arm64",
		},
		"crictl": {
			true:  "amd64",
			false: "arm64",
		},
		"kube": {
			true:  "amd64",
			false: "arm64",
		},
	}
}

func LoadVars() error {
	AkId = os.Getenv("CLOUD_KERNEL_AKID")
	if AkId == "" {
		return errors.New("环境变量CLOUD_KERNEL_AKID未设置，请设置后重试")
	}
	AkSK = os.Getenv("CLOUD_KERNEL_AKSK")
	if AkSK == "" {
		return errors.New("环境变量CLOUD_KERNEL_AKSK未设置，请设置后重试")
	}
	if KubeVersion == "" {
		KubeVersion = os.Getenv("CLOUD_KUBE_VERSION")
		if KubeVersion == "" {
			return errors.New("环境变量CLOUD_KUBE_VERSION未设置，请设置后重试")
		}
	}
	isAmd64 := IsAmd64()
	KubeShell = fmt.Sprintf(KubeShell, KubeVersion, platform()["kube"][isAmd64])
	DingDing = os.Getenv("CLOUD_KERNEL_DINGDING")
	MarketCtlToken = os.Getenv("CLOUD_KERNEL_MARKETCTL")
	//sshcmd
	sshCmdVersion := defaultSSHCmdVersion
	if v := os.Getenv("CLOUD_KERNEL_SSHCMD_VERSION"); v != "" {
		sshCmdVersion = v
	}
	SSHCmdDownload = fmt.Sprintf(SSHCmdDownload, sshCmdVersion, platform()["sshcmd"][isAmd64])
	//sealos
	sealosVersion := defaultSealosVersion
	if v := os.Getenv("CLOUD_KERNEL_SEALOS_VERSION"); v != "" {
		sealosVersion = v
	}
	SealosDownload = fmt.Sprintf(SealosDownload, sealosVersion, platform()["sealos"][isAmd64])
	//marketctl
	marketCtlVersion := defaultMarketCtlVersion
	if v := os.Getenv("CLOUD_KERNEL_MARKET_CTL_VERSION"); v != "" {
		marketCtlVersion = v
	}
	MarketCtlDownload = fmt.Sprintf(MarketCtlDownload, marketCtlVersion, platform()["marketctl"][isAmd64])

	//containerd
	containerdVersion := defaultContainerdVersion
	if v := os.Getenv("CLOUD_KERNEL_CONTAINERD_VERSION"); v != "" {
		containerdVersion = v
	}
	ContainerdShell = fmt.Sprintf(ContainerdShell, containerdVersion, platform()["containerd"][isAmd64])

	dockerVersion := defaultDockerVersion
	if v := os.Getenv("CLOUD_KERNEL_DOCKER_VERSION"); v != "" {
		dockerVersion = v
	}
	DockerShell = fmt.Sprintf(DockerShell, platform()["docker"][isAmd64], dockerVersion)

	crictlVersion := defaultCriCtlVersion
	if v := os.Getenv("CLOUD_KERNEL_CRICTL_VERSION"); v != "" {
		crictlVersion = v
	}
	CrictlShell = fmt.Sprintf(CrictlShell, crictlVersion, platform()["crictl"][isAmd64])
	return nil
}

const MarketYaml = `
market:
  body:
    spec:
      name: v%s
      price: %.2f
      product:
        class: cloud_kernel
        productName: kubernetes
      url: /tmp/kube%s.tar.gz
    status:
      productVersionStatus: ONLINE
  kind: productVersion`

const MarketArmYaml = `
market:
  body:
    spec:
      name: v%s
      price: %.2f
      product:
        class: cloud_kernel
        productName: kubernetes-arm
      url: /tmp/kube%s-arm64.tar.gz
    status:
      productVersionStatus: ONLINE
  kind: productVersion`
