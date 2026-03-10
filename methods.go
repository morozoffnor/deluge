package deluge

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
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
		fmt.Println("error marshalling request", err)
		return nil, err
	}
	request, err := http.NewRequest("POST", d.baseURL+"/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("error creating http request", err)
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
		fmt.Println("error sending request", err)
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
		fmt.Println("could not read response body", err)
		return nil, err
	}
	var dResp delugeResp
	if err = json.Unmarshal(respBody, &dResp); err != nil {
		return nil, err
	}

	return &dResp, nil
}

func (d *DelugeClient) AddTorrentMagnet(magnetURI string, opts TorrentOptions) (string, error) {
	resp, err := d.call("core.add_torrent_magnet", []interface{}{magnetURI, opts.ToMap()})
	if err != nil {
		return "", err
	}

	if resp.Error != nil {
		return "", fmt.Errorf("deluge error: %v", resp.Error)
	}

	var torrentID string
	err = json.Unmarshal(resp.Result, &torrentID)
	if err != nil {
		return "", err
	}
	return torrentID, nil
}

func (d *DelugeClient) AddTorrentFile(filePath string, opts TorrentOptions) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	fileName := file.Name()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("error reading file", err)
		return "", err
	}

	var fileDump bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &fileDump)
	encoder.Write(data)
	encoder.Close()

	resp, err := d.call("core.add_torrent_file", []interface{}{fileName, fileDump.String(), opts.ToMap()})
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", fmt.Errorf("deluge error: %v", resp.Error)
	}

	var torrentID string
	err = json.Unmarshal(resp.Result, &torrentID)
	if err != nil {
		return "", err
	}
	return torrentID, nil
}

type TorrentOptions struct {
	DownloadLocation    string  // where to save files
	MoveCompleted       bool    // move on completion
	MoveCompletedPath   string  // where to move when done
	MaxDownloadSpeed    float64 // KB/s, -1 for unlimited
	MaxUploadSpeed      float64 // KB/s, -1 for unlimited
	MaxConnections      int     // -1 for unlimited
	MaxUploadSlots      int     // -1 for unlimited
	Paused              bool    // add in paused state
	PrioritizeFirstLast bool    // prioritize first/last pieces
	SeedMode            bool    // skip hash check
	SuperSeeding        bool
	AddPaused           bool // same as Paused
}

func (o TorrentOptions) ToMap() map[string]interface{} {
	m := map[string]interface{}{}

	if o.DownloadLocation != "" {
		m["download_location"] = o.DownloadLocation
	}
	if o.MoveCompleted {
		m["move_completed"] = true
		m["move_completed_path"] = o.MoveCompletedPath
	}
	if o.MaxDownloadSpeed != 0 {
		m["max_download_speed"] = o.MaxDownloadSpeed
	}
	if o.MaxUploadSpeed != 0 {
		m["max_upload_speed"] = o.MaxUploadSpeed
	}
	if o.MaxConnections != 0 {
		m["max_connections"] = o.MaxConnections
	}
	if o.MaxUploadSlots != 0 {
		m["max_upload_slots"] = o.MaxUploadSlots
	}
	if o.Paused || o.AddPaused {
		m["add_paused"] = true
	}
	if o.PrioritizeFirstLast {
		m["prioritize_first_last_pieces"] = true
	}
	if o.SeedMode {
		m["seed_mode"] = true
	}
	if o.SuperSeeding {
		m["super_seeding"] = true
	}

	return m
}
