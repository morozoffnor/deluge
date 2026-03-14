package deluge

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newMockServerWithHandler(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func TestAddTorrentMagnet(t *testing.T) {
	tests := []struct {
		name            string
		magnetURI       string
		opts            TorrentOptions
		serverResponse  delugeResp
		expectError     bool
		expectTorrentID string
		skipServer      bool
	}{
		{
			name:      "successful add",
			magnetURI: "magnet:?xt=urn:btih:abc123",
			serverResponse: delugeResp{
				Result: []byte(`"abc123def456"`),
			},
			expectTorrentID: "abc123def456",
		},
		{
			name:      "successful add with options",
			magnetURI: "magnet:?xt=urn:btih:abc123",
			opts: TorrentOptions{
				DownloadLocation: "/downloads",
				Paused:           true,
				MaxDownloadSpeed: 1024,
			},
			serverResponse: delugeResp{
				Result: []byte(`"abc123def456"`),
			},
			expectTorrentID: "abc123def456",
		},
		{
			name:      "server returns error",
			magnetURI: "magnet:?xt=urn:btih:abc123",
			serverResponse: delugeResp{
				Result: []byte(`null`),
				Error:  "torrent already exists",
			},
			expectError: true,
		},
		{
			name:        "empty magnet URI",
			magnetURI:   "",
			expectError: true,
			skipServer:  true,
		},
		{
			name:        "invalid magnet format",
			magnetURI:   "http://example.com/torrent",
			expectError: true,
			skipServer:  true,
		},
		{
			name:      "server returns null result",
			magnetURI: "magnet:?xt=urn:btih:abc123",
			serverResponse: delugeResp{
				Result: []byte(`null`),
			},
			expectTorrentID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipServer {
				client, _ := New("http://localhost:1", "", "", "")
				_, err := client.AddTorrentMagnet(tt.magnetURI, tt.opts)
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			srv := newMockServerWithHandler(func(w http.ResponseWriter, r *http.Request) {
				var req delugeReq
				json.NewDecoder(r.Body).Decode(&req)

				resp := tt.serverResponse
				resp.ID = req.ID
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)
			})
			defer srv.Close()

			client, _ := New(srv.URL, "", "", "")
			torrentID, err := client.AddTorrentMagnet(tt.magnetURI, tt.opts)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if torrentID != tt.expectTorrentID {
				t.Errorf("expected torrentID %q, got %q", tt.expectTorrentID, torrentID)
			}
		})
	}
}

func TestAddTorrentMagnet_RequestParams(t *testing.T) {
	var capturedReq delugeReq

	srv := newMockServerWithHandler(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&capturedReq)

		resp := delugeResp{
			ID:     capturedReq.ID,
			Result: []byte(`"torrent123"`),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	defer srv.Close()

	client, _ := New(srv.URL, "", "", "")
	magnet := "magnet:?xt=urn:btih:abc123"
	opts := TorrentOptions{
		DownloadLocation: "/downloads",
		MaxDownloadSpeed: 500,
		Paused:           true,
	}

	_, err := client.AddTorrentMagnet(magnet, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedReq.Method != "core.add_torrent_magnet" {
		t.Errorf("expected method %q, got %q", "core.add_torrent_magnet", capturedReq.Method)
	}

	if len(capturedReq.Params) != 2 {
		t.Fatalf("expected 2 params, got %d", len(capturedReq.Params))
	}

	if capturedReq.Params[0] != magnet {
		t.Errorf("expected first param %q, got %v", magnet, capturedReq.Params[0])
	}

	optsMap, ok := capturedReq.Params[1].(map[string]interface{})
	if !ok {
		t.Fatal("expected second param to be a map")
	}

	if optsMap["download_location"] != "/downloads" {
		t.Errorf("expected download_location %q, got %v", "/downloads", optsMap["download_location"])
	}
	if optsMap["max_download_speed"] != float64(500) {
		t.Errorf("expected max_download_speed 500, got %v", optsMap["max_download_speed"])
	}
	if optsMap["add_paused"] != true {
		t.Errorf("expected add_paused true, got %v", optsMap["add_paused"])
	}
}
