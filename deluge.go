package deluge

import (
	"encoding/json"
	"fmt"
)

type DelugeClient struct {
	baseURL           string
	password          string
	basicAuthLogin    string
	basicAuthPassword string
	cookie            string
}

func New(baseUrl string, pass string, basicAuthLogin string, basicAuthPass string) (*DelugeClient, error) {
	if baseUrl == "" {
		return nil, fmt.Errorf("empty deluge base url")
	}

	d := &DelugeClient{
		baseURL:           baseUrl,
		password:          pass,
		basicAuthLogin:    basicAuthLogin,
		basicAuthPassword: basicAuthPass,
	}
	return d, nil
}

func (d *DelugeClient) Login() (bool, error) {
	resp, err := d.call("auth.login", []interface{}{d.password})
	if err != nil {
		return false, err
	}
	var success bool
	err = json.Unmarshal(resp.Result, &success)
	if err != nil {
		return false, err
	}
	return success, nil
}
