package github

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sealyun/cloud-kernel/pkg/logger"
	"github.com/sealyun/cloud-kernel/pkg/vars"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	githubURL     = "https://api.github.com/repos/kubernetes/kubernetes/releases"
	sealyunAmdURL = "https://www.sealyun.com/api/v2/classes/cloud_kernel/products/kubernetes/versions"
	sealyunArmURL = "https://www.sealyun.com/api/v2/classes/cloud_kernel/products/kubernetes-arm64/versions"
)

type release struct {
	CreatedAt time.Time `json:"created_at"`
	TagName   string    `json:"tag_name"`
	Draft     bool      `json:"draft"`
}
type sealyunVersion struct {
	Code int `json:"code"`
	Data []struct {
		Name string `json:"name"`
	} `json:"data"`
}

func Fetch() ([]string, error) {
	packageOffline := make([]string, 0)
	tags := fetchTags()
	sealyunV := fetchSealyunTags()
	if len(sealyunV) == 0 {
		return packageOffline, errors.New("获取sealyun的tag失败")
	}
	for _, tag := range tags {
		logger.Debug("当月github发布有效版本:" + tag)
		if v, ok := sealyunV[tag]; ok && v != "" {
			logger.Debug(fmt.Sprintf("当月github发布有效版本: %s, 已经上传到sealyun。", tag))
			continue
		} else {
			packageOffline = append(packageOffline, tag)
		}
	}
	return packageOffline, nil
}

func fetchTags() []string {
	tags := make([]string, 0)
	for i := 1; ; i++ {
		u := fmt.Sprintf("%s?page=%d", githubURL, i)
		var releases []release
		data, _ := getUrl(u)
		if data != nil {
			_ = json.Unmarshal(data, &releases)
			if len(releases) > 0 {
				for _, v := range releases {
					if v.CreatedAt.AddDate(0, 1, 0).After(time.Now()) {
						//一个月内的数据
						if v.Draft {
							continue
						}
						if strings.ContainsAny(v.TagName, "-") {
							continue
						}
						tags = append(tags, v.TagName)
					} else {
						return tags
					}
				}
			} else {
				break
			}
		}
	}

	return tags
}
func fetchSealyunTags() map[string]string {
	tags := make(map[string]string, 0)
	var u string
	if vars.IsArm64 {
		u = sealyunArmURL
	} else {
		u = sealyunAmdURL
	}
	var versions sealyunVersion
	data, _ := getUrl(u)
	if data != nil {
		_ = json.Unmarshal(data, &versions)
		if versions.Code == 200 {
			for _, v := range versions.Data {
				tags[v.Name] = v.Name
			}
		}
	}
	return tags
}

func getUrl(rawurl string) ([]byte, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, fmt.Errorf("解析url为空")
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, resp.Body)
	_ = resp.Body.Close()
	return ioutil.ReadAll(buf)
}
