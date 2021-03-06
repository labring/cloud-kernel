/*
Copyright 2021 cuisongliu@qq.com.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package _package

import (
	"github.com/labring/cloud-kernel/pkg/utils"
	"github.com/labring/cloud-kernel/pkg/vars"
)

func getCNIVersion() (string, string) {
	k8s := vars.KubeVersion
	if k8s != "" {
		if utils.For119(k8s) {
			return "31901", "v3.19.1"
		}
	}
	return "30802", "v3.8.2"
}
