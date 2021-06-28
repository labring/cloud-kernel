package _package

import (
	"github.com/sealyun/cloud-kernel/pkg/vars"
	"testing"
)

func TestPackageDocker(t *testing.T) {
	vars.Testing = false
	vars.LoadAKSK()
	Package("1.19.9")
}

func TestPackageContainerd(t *testing.T) {
	vars.Testing = false
	vars.LoadAKSK()
	Package("1.20.0")
}

func TestSaveImage(t *testing.T) {
	//Package(vars.Docker, "1.20.0", "19.03.12")
	k8s := NewDockerK8s("8.210.233.63")
	k8s.SavePackage()
}
