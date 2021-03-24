package vars

import "os"

var (
	DingDing string
	AkId     string
	AkSK     string
)

func LoadVars() {
	AkId = os.Getenv("CLOUD_KERNEL_AKID")
	AkSK = os.Getenv("CLOUD_KERNEL_AKSK")
	DingDing = os.Getenv("CLOUD_KERNEL_DINGDING")
}
