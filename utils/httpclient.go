package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func Request(method string, url string, header http.Header, body interface{}, token string) (int, []byte, http.Header, error) {
	var req *http.Request
	var err error
	if body != nil {
		data, _ := json.Marshal(body)
		payload := strings.NewReader(string(data))
		req, err = http.NewRequest(method, url, payload)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return 0, nil, nil, err
	}
	if body != nil {
		req.Header.Add("Content-Type", "application/json; charset=utf-8")
	}
	if len(token) != 0 {
		req.Header.Add("Authorization", token)
	}
	if header != nil {
		for k, v := range header {
			req.Header.Add(k, v[0])
		}
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	response, _ := ioutil.ReadAll(res.Body)
	return res.StatusCode, response, res.Header, err
}
