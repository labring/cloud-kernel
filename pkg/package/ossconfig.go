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

import "github.com/labring/cloud-kernel/pkg/logger"

///oss
const ossTemplate = `[Credentials]
language=CH
endpoint=oss-accelerate.aliyuncs.com
accessKeyID={{.KEY_ID}}
accessKeySecret={{.KEY_SECRET}}
`

type OSSConfig struct {
	KeyID     string
	KeySecret string
}

func (m *OSSConfig) Template() string {
	return ossTemplate
}

func (m *OSSConfig) TemplateConvert() string {
	p := map[string]interface{}{
		"KEY_ID":     m.KeyID,
		"KEY_SECRET": m.KeySecret,
	}

	data, err := templateFromContent(m.Template(), p)
	if err != nil {
		logger.Error(err) //nolint:typecheck
		return ""
	}
	return data
}
