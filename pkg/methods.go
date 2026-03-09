package deluge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
)

type delugeReq struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
	ID     int           `json:"id"`
}

type delugeResp struct {
	ID     int             `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  interface{}     `json:"error"`
}

func (d *DelugeClient) call(method string, params []interface{}) (*delugeResp, error) {
	req := delugeReq{
		Method: method,
		Params: params,
		ID:     rand.Int(),
	}

	body, err := json.Marshal(req)
	if err != nil {
		fmt.Println("[DelugeWorker] error marshalling request", err)
		return nil, err
	}
	request, err := http.NewRequest("POST", d.baseURL+"/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("[DelugeWorker] error creating http request", err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	// auth logic
	if d.basicAuthLogin != "" && d.basicAuthPassword != "" {
		request.SetBasicAuth(d.basicAuthLogin, d.basicAuthPassword)
	}

	if d.cookie != "" {
		request.Header.Set("Cookie", d.cookie)
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("[DelugeWorker] error sending request", err)
		return nil, err
	}

	defer resp.Body.Close()

	// store session cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "_session_id" {
			d.cookie = "_session_id=" + cookie.Value
		}
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[DelugeWorker] could not read response body", err)
		return nil, err
	}
	var dResp delugeResp
	if err = json.Unmarshal(respBody, &dResp); err != nil {
		return nil, err
	}

	return &dResp, nil
}

func (d *DelugeClient) waitUpdates(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case upd := <-d.updates:
			fmt.Println("[Deluge] recieved an update", upd)
			// process update
		}
	}
}
