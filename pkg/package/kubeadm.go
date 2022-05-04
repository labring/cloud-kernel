/*
Copyright 2022 cuisongliu@qq.com.

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
	"github.com/labring/cloud-kernel/pkg/logger"
	"strconv"
	"strings"
)

const (
	KubeadmV1beta1 = "kubeadm.k8s.io/v1beta1"
	KubeadmV1beta2 = "kubeadm.k8s.io/v1beta2"
	KubeadmV1beta3 = "kubeadm.k8s.io/v1beta3"
)

func getKubeadmAPI(version string) string {
	var KubeadmAPI string
	major, _ := GetMajorMinorInt(version)
	switch {
	//
	case major < 120:
		KubeadmAPI = KubeadmV1beta1
	case major < 123 && major >= 120:
		KubeadmAPI = KubeadmV1beta2
	case major >= 123:
		KubeadmAPI = KubeadmV1beta3
	default:
		KubeadmAPI = KubeadmV1beta3
	}
	logger.Debug("KubeadmApi: %s", KubeadmAPI)
	return KubeadmAPI
}

// GetMajorMinorInt
func GetMajorMinorInt(version string) (major, minor int) {
	// alpha beta rc version
	if strings.Contains(version, "-") {
		v := strings.Split(version, "-")[0]
		version = v
	}
	version = strings.Replace(version, "v", "", -1)
	versionArr := strings.Split(version, ".")
	if len(versionArr) >= 2 {
		majorStr := versionArr[0] + versionArr[1]
		minorStr := versionArr[2]
		if major, err := strconv.Atoi(majorStr); err == nil {
			if minor, err := strconv.Atoi(minorStr); err == nil {
				return major, minor
			}
		}
	}
	return 0, 0
}
