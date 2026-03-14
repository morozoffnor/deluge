package deluge

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type delugeMockServer struct {
	*httptest.Server
	responses map[string]interface{}
	requests  []delugeReq
}

func newMockServer() *delugeMockServer {
	mock := &delugeMockServer{
		responses: make(map[string]interface{}),
		requests:  []delugeReq{},
	}

	mock.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req delugeReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		mock.requests = append(mock.requests, req)

		response := mock.responses[req.Method]

		resp := delugeResp{
			ID:     req.ID,
			Result: mustMarshal(response),
			Error:  nil,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))

	return mock
}

func mustMarshal(v interface{}) json.RawMessage {
	ok, _ := json.Marshal(v)
	return ok
}

func (m *delugeMockServer) Close() {
	m.Server.Close()
}

func (m *delugeMockServer) setResponse(method string, result interface{}) {
	m.responses[method] = result
}

func (m *delugeMockServer) getLastRequest() *delugeReq {
	if len(m.requests) == 0 {
		return nil
	}
	return &m.requests[len(m.requests)-1]
}

func TestDelugeClient_Login(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		result      bool
		expectError bool
	}{
		{
			name:        "success",
			password:    "truepass",
			result:      true,
			expectError: false,
		},
		{
			name:        "fail",
			password:    "falsepass",
			result:      false,
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockServer()
			defer mock.Close()

			mock.setResponse("auth.login", tt.result)

			client, _ := New(mock.URL, "truepass", "", "")
			ok, err := client.Login()
			if err != nil {
				t.Errorf("DelugeClient.Login() error = %v", err)
			}

			if tt.expectError == true && ok != false {
				t.Errorf("DelugeClient.Login() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if tt.expectError == false && ok != true {
				t.Errorf("DelugeClient.Login() ok = %v, expectOK %v", ok, tt.expectError)
			}

			req := mock.getLastRequest()
			if req.Method != "auth.login" {
				t.Errorf("DelugeClient.Login() method = %v, expectMethod auth.login", req.Method)
			}
		})
	}
}
