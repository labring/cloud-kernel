package marketctl

import (
	"fmt"
	"github.com/labring/cloud-kernel/pkg/vars"
)

const defaultDomain = "https://www.sealyun.com"

func Healthy() error {
	uri := fmt.Sprintf("/api/v2/healthy")
	return Do(defaultDomain, uri, "GET", vars.MarketCtlToken, nil)
}
