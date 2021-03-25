package vars

import "os"

var (
	DingDing string
	AkId     string
	AkSK     string
)

type Runtime string

const (
	EcsPassword    = "Fanux#123"
	SSHCmdDownload = "https://github.com/cuisongliu/sshcmd/releases/download/v1.5.3/sshcmd"
	SealosDownload = "https://sealyun.oss-accelerate.aliyuncs.com/latest/sealos"
)

const (
	Docker     Runtime = "docker"
	Containerd Runtime = "containerd"
)

func LoadVars() {
	AkId = os.Getenv("CLOUD_KERNEL_AKID")
	AkSK = os.Getenv("CLOUD_KERNEL_AKSK")
	DingDing = os.Getenv("CLOUD_KERNEL_DINGDING")
}
