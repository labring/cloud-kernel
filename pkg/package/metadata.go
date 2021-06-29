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
	"bytes"
	"text/template"

	"github.com/sealyun/cloud-kernel/pkg/logger"
)

///Metadata
const metadataTemplate = `{
  "k8sVersion": "{{.k8sVersion}}",
  "cniVersion": "{{.cniVersion}}",
  "cniName": "{{.cniName}}"
}
`

type Metadata struct {
	K8sVersion string
	CniVersion string
	CniName    string
}

func (m *Metadata) Template() string {
	return metadataTemplate
}

func (m *Metadata) TemplateConvert() string {
	if m.CniName == "" {
		m.CniName = "calico"
	}
	p := map[string]interface{}{
		"k8sVersion": m.K8sVersion,
		"cniVersion": m.CniVersion,
		"cniName":    m.CniName,
	}

	data, err := templateFromContent(m.Template(), p)
	if err != nil {
		logger.Error(err) //nolint:typecheck
		return ""
	}
	return data
}

func templateFromContent(templateContent string, param map[string]interface{}) (string, error) {
	tmpl, err := template.New("text").Parse(templateContent)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, param)
	bs := buffer.Bytes()
	if len(bs) > 0 {
		return string(bs), nil
	}
	return "", err
}
