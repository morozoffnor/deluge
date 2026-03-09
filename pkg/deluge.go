package deluge

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

const (
	defaultWorkerCount    = 1
	defaultUpdatesChanBuf = 10
)

type DelugeClient struct {
	baseURL           string
	password          string
	basicAuthLogin    string
	basicAuthPassword string
	updates           chan *DelugeUpdate
	cookie            string
}

type DelugeUpdate struct {
	UpdateType int
	UpdateID   int
	Content    any
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
		updates:           make(chan *DelugeUpdate, defaultUpdatesChanBuf),
	}
	return d, nil
}

func (d *DelugeClient) Start(ctx context.Context) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	// login
	fmt.Println("[Deluge] trying to log in")
	_, err := d.login()
	if err != nil {
		fmt.Println("[Deluge] failed to log in")
		return
	}
	fmt.Println("[Deluge] logged in")
	for i := 0; i < defaultWorkerCount; i++ {
		go d.waitUpdates(ctx, &wg)
	}
	wg.Wait()

}

func (d *DelugeClient) login() (bool, error) {
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
