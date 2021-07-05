package _package

import (
	"testing"

	"github.com/sealyun/cloud-kernel/pkg/vars"
)

//go test -v -timeout 30000s -test.run TestPackageDocker

func TestPackageDocker(t *testing.T) {
	vars.Testing = false
	vars.Uploading = false
	vars.LoadAKSK()
	Package("1.19.9")
}

func TestPackageContainerd(t *testing.T) {
	vars.Testing = false
	vars.Uploading = false
	vars.LoadAKSK()
	Package("1.20.0")
}
