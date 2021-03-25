package vars

import "os"

var (
	DingDing       string
	AkId           string
	AkSK           string
	MarketCtlToken string
)

type Runtime string

const (
	EcsPassword              = "Fanux#123"
	SSHCmdDownload           = "https://github.com/cuisongliu/sshcmd/releases/download/v1.5.3/sshcmd"
	SealosDownload           = "https://sealyun.oss-accelerate.aliyuncs.com/latest/sealos"
	MarketCtlDownload        = "https://sealyun-market.oss-accelerate.aliyuncs.com/marketctl/v1.0.5/linux_amd64/marketctl"
	CriCtlVersion            = "1.20.0"
	DefaultDockerVersion     = "19.03.12"
	DefaultContainerdVersion = "1.3.9"
)

const (
	Docker     Runtime = "docker"
	Containerd Runtime = "containerd"
)

func LoadVars() {
	AkId = os.Getenv("CLOUD_KERNEL_AKID")
	AkSK = os.Getenv("CLOUD_KERNEL_AKSK")
	DingDing = os.Getenv("CLOUD_KERNEL_DINGDING")
	MarketCtlToken = os.Getenv("CLOUD_KERNEL_MARKETCTL")
}

const MarketYaml = `
market:
  body:
    spec:
      name: v%s
      price: %d
      product:
        class: cloud_kernel
        productName: kubernetes
      url: /tmp/kube%s.tar.gz
    status:
      productVersionStatus: ONLINE
  kind: productVersion`
