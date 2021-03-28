package marketctl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/wonderivan/logger"
	"io/ioutil"
	"net/http"
)

func Do(domain, uri, method, accessToken string, post []byte) error {
	_, err := do(domain, uri, method, accessToken, nil, post)
	return err
}
func DoBody(domain, uri, method, accessToken string, post []byte) ([]byte, error) {
	return do(domain, uri, method, accessToken, nil, post)
}
func DoBodyAddHeader(domain, uri, method, accessToken string, headers map[string]string, post []byte) ([]byte, error) {
	return do(domain, uri, method, accessToken, headers, post)
}
func do(domain, uri, method, accessToken string, headers map[string]string, post []byte) ([]byte, error) {
	req, err := http.NewRequest(method, domain+uri, bytes.NewReader(post))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-agent", "MarketCtl")
	req.Header.Set("AccessToken", accessToken)
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		//logger.Error("response error is %s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		logger.Error(resp.Status)
		return nil, fmt.Errorf("respone status is not correct")
	} else {
		var out map[string]interface{}
		_ = json.Unmarshal(body, &out)
		if code, ok := out["code"].(float64); ok && code != 200 {
			return nil, fmt.Errorf(out["message"].(string))
		}
	}
	return body, nil
}
