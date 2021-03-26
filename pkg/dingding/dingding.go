package dingding

import (
	"bytes"
	"github.com/sealyun/cloud-kernel/pkg/httplib"
	"github.com/sealyun/cloud-kernel/pkg/logger"
	"github.com/sealyun/cloud-kernel/pkg/vars"
	"io/ioutil"
	"log"
	"text/template"
)

const linkBody = `{
    "msgtype": "link", 
    "link": {
        "text": "{{.text}}", 
        "title": "{{.title}}", 
        "picUrl": "", 
        "messageUrl": "{{.url}}"
    },
    "at": {
        "isAtAll": {{.at_all}}
    }
}`

const textBody = `{
    "msgtype": "text", 
    "text": {
        "content": "{{.text}}"
    },
    "at": {
        "isAtAll": {{.at_all}}
    }
}`

func DingdingText(text string, atAll bool) {
	var envMap = make(map[string]interface{})
	envMap["text"] = text
	envMap["at_all"] = atAll
	dingdingFromMap(envMap, true, vars.DingDing)
}

func DingdingLink(title, text, url string, atAll bool) {
	var envMap = make(map[string]interface{})
	envMap["title"] = title
	envMap["text"] = text
	envMap["url"] = url
	envMap["at_all"] = atAll
	dingdingFromMap(envMap, false, vars.DingDing)
}

func dingdingFromMap(data map[string]interface{}, text bool, token string) {
	if token == "" {
		logger.Warn("钉钉Token未配置,跳过通知")
		return
	}
	r := httplib.Post("https://oapi.dingtalk.com/robot/send?access_token=" + token)
	var templateContent string
	if text {
		templateContent = textBody
	} else {
		templateContent = linkBody
	}
	tmpl, _ := template.New("text").Parse(templateContent)
	if tmpl != nil {
		var buffer bytes.Buffer
		_ = tmpl.Execute(&buffer, data)
		dd := buffer.String()
		resp, _ := r.Body(dd).Header("Content-Type", "application/json").DoRequest()
		result, _ := ioutil.ReadAll(resp.Body)
		log.Println("DingDing: " + string(result))
	}

}
