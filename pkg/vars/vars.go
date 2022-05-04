package vars

import (
	"fmt"
	"os"
)

var (
	DingDing          string
	AkID              string
	AkSK              string
	OSSAkID           string
	OSSAkSK           string
	MarketCtlToken    string
	IsArm64           bool
	Uploading         bool
	Testing           bool
	SSHCmdDownload    string
	SealosDownload    string
	MarketCtlDownload string
	OSSUtilDownload   string
	DockerShell       string
	ContainerdShell   string
	CrictlShell       string
	KubeShell         string
	NerdctlShell      string
	LibraryURL        string
	KubeVersion       string
	DefaultPrice      float64
	DefaultZeroPrice  float64
	DefaultClass      = "cloud_kernel" //cloud_kernel
	DefaultProduct    = "kubernetes"   //kubernetes

	defaultSealosVersion        = "3.3.9-rc.1"
	defaultMarketCtlVersion     = "1.0.5" //v1.0.5
	defaultSSHCmdVersion        = "1.5.5"
	defaultNerdctlVersion       = "0.19.0"
	defaultCriCtlVersion        = "1.22.0"
	defaultDockerVersion        = "19.03.12"
	defaultContainerdVersion    = "1.6.4"
	defaultContainerdArmVersion = "1.6.4"
)

const (
	EcsPassword          = "Fanux#123"
	FmtOSSUtilDownload   = "https://gosspublic.alicdn.com/ossutil/1.7.3/ossutil%s"                               //https://gosspublic.alicdn.com/ossutil/1.7.3/ossutilarm64
	FmtSSHCmdDownload    = "https://github.com/cuisongliu/sshcmd/releases/download/v%s/sshcmd%s"                 //sshcmd-arm64
	FmtSealosDownload    = "https://sealyun-home.oss-accelerate.aliyuncs.com/sealos/v%s/sealos%s"                //sealos-arm64
	FmtMarketCtlDownload = "https://sealyun-market.oss-accelerate.aliyuncs.com/marketctl/v%s/linux_%s/marketctl" //linux_arm64
	FmtDockerShell       = "wget https://download.docker.com/linux/static/stable/%s/docker-%s.tgz  " +           //aarch64
		"-O docker.tgz && cp docker.tgz kube/docker/docker.tgz"
	FmtContainerdShell = "wget https://github.com/containerd/containerd/releases/download/v%s/cri-containerd-cni-%s-linux-%s.tar.gz " +
		"-O cri-containerd-cni-linux.tar.gz && " +
		"cp cri-containerd-cni-linux.tar.gz kube/containerd/cri-containerd-cni-linux.tar.gz"
	FmtCrictlShell = "wget https://github.com/kubernetes-sigs/cri-tools/releases/download/v%s/crictl-v%s-linux-%s.tar.gz " +
		"-O  crictl.tar.gz &&  tar xf crictl.tar.gz && cp crictl kube/bin/"
	FmtKubeShell = "wget https://dl.k8s.io/v%s/kubernetes-server-linux-%s.tar.gz -O  kubernetes-server.tar.gz && " + //arm64
		"tar zxvf kubernetes-server.tar.gz"
	FmtNerdctlShell = "wget https://github.com/containerd/nerdctl/releases/download/v%s/nerdctl-%s-linux-%s.tar.gz " +
		"-O  nerdctl.tar.gz &&  tar xf nerdctl.tar.gz && cp nerdctl kube/bin/ && cp containerd-rootless* kube/bin/"
	FmtLibraryURL = "https://sealyun-home.oss-accelerate.aliyuncs.com/images/library-2.5-%s-%s.tar.gz"
)

func platform() map[string]map[bool]string {
	return map[string]map[bool]string{
		"sshcmd": {
			false: "",
			true:  "-arm64",
		},
		"oss": {
			false: "64",
			true:  "arm64",
		},
		"sealos": {
			false: "",
			true:  "-arm64",
		},
		"nerdctl": {
			false: "amd64",
			true:  "arm64",
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
		"library": {
			false: "amd64",
			true:  "arm64",
		},
	}
}

func loadEnv() {
	if v := os.Getenv("SEALOS_VERSION"); v != "" {
		defaultSealosVersion = v
	}
	if v := os.Getenv("MARKET_CTL_VERSION"); v != "" {
		defaultMarketCtlVersion = v
	}
	if v := os.Getenv("SSH_CMD_VERSION"); v != "" {
		defaultSSHCmdVersion = v
	}
	if v := os.Getenv("NERD_CTL_VERSION"); v != "" {
		defaultNerdctlVersion = v
	}
	if v := os.Getenv("CRI_CTL_VERSION"); v != "" {
		defaultCriCtlVersion = v
	}
	if v := os.Getenv("DOCKER_VERSION"); v != "" {
		defaultDockerVersion = v
	}
	if v := os.Getenv("CONTAINERD_VERSION"); v != "" {
		defaultContainerdVersion = v
	}
}
func LoadAKSK() {
	if v := os.Getenv("OSS_AKID"); v != "" {
		OSSAkID = v
	}
	if v := os.Getenv("OSS_AKSK"); v != "" {
		OSSAkSK = v
	}
	if v := os.Getenv("CLOUD_KERNEL_AKID"); v != "" {
		AkID = v
	}
	if v := os.Getenv("CLOUD_KERNEL_AKSK"); v != "" {
		AkSK = v
	}
}
func LoadVars() error {
	loadEnv()
	KubeShell = fmt.Sprintf(FmtKubeShell, KubeVersion, platform()["kube"][IsArm64])
	//sshcmd
	SSHCmdDownload = fmt.Sprintf(FmtSSHCmdDownload, defaultSSHCmdVersion, platform()["sshcmd"][IsArm64])
	//oss
	OSSUtilDownload = fmt.Sprintf(FmtOSSUtilDownload, platform()["oss"][IsArm64])
	//sealos
	SealosDownload = fmt.Sprintf(FmtSealosDownload, defaultSealosVersion, platform()["sealos"][IsArm64])
	//marketctl
	MarketCtlDownload = fmt.Sprintf(FmtMarketCtlDownload, defaultMarketCtlVersion, platform()["marketctl"][IsArm64])
	//containerd
	var ContainerdVersion string
	if IsArm64 {
		ContainerdVersion = defaultContainerdArmVersion
	} else {
		ContainerdVersion = defaultContainerdVersion
	}
	ContainerdShell = fmt.Sprintf(FmtContainerdShell, ContainerdVersion, ContainerdVersion, platform()["containerd"][IsArm64])
	DockerShell = fmt.Sprintf(FmtDockerShell, platform()["docker"][IsArm64], defaultDockerVersion)
	CrictlShell = fmt.Sprintf(FmtCrictlShell, defaultCriCtlVersion, defaultCriCtlVersion, platform()["crictl"][IsArm64])
	NerdctlShell = fmt.Sprintf(FmtNerdctlShell, defaultNerdctlVersion, defaultNerdctlVersion, platform()["nerdctl"][IsArm64])
	LibraryURL = fmt.Sprintf(FmtLibraryURL, "linux", platform()["library"][IsArm64])

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
