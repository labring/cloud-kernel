package _package

import (
	"github.com/sealyun/cloud-kernel/pkg/vars"
	"testing"
)

func TestPackage(t *testing.T) {
	Package(vars.Docker, "1.20.0", "19.03.12")
}
