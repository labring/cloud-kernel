package ecs

import (
	"github.com/cuisongliu/cloud-kernel/pkg/vars"
	"testing"
)

func TestNew(t *testing.T) {
	vars.AkId = ""
	vars.AkSK = ""
	New()
}
