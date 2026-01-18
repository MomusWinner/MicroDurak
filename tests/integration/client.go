package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

type HTTPClient struct {
	HTTPClient *http.Client
	AuthToken  string
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: true,
			},
		},
	}
}

func (c *HTTPClient) Request(
	t *testing.T,
	method, url string,
	body any,
) (*http.Response, []byte) {
	t.Helper()

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)

		if err != nil {
			t.Error(err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		t.Error(err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.AuthToken != "" {
		req.Header.Set("Authorization", c.AuthToken)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	return resp, respBody
}

func (c *HTTPClient) GET(t *testing.T, url string) []byte {
	_, body := c.Request(t, "GET", url, nil)
	return body
}

func (c *HTTPClient) POST(t *testing.T, url string, data interface{}) (*http.Response, []byte) {
	return c.Request(t, "POST", url, data)
}

type WSClient struct {
	AuthToken string
	Conn      *websocket.Conn
}

func NewWSClient(authToken string) *WSClient {
	return &WSClient{
		AuthToken: authToken,
	}
}

func (c *WSClient) Connect(t *testing.T, url string) {
	t.Helper()

	headers := http.Header{}
	headers.Set("Authorization", c.AuthToken)

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, resp, err := dialer.Dial(url, headers)
	if err != nil {
		t.Error("WebSocket connection failed", err)

	}

	if http.StatusSwitchingProtocols != resp.StatusCode {
		t.Error("Status code should be 101 (SwitchingProtocols)")
	}

	c.Conn = conn
}

func (c *WSClient) Send(t *testing.T, message any) {
	t.Helper()

	if c.Conn == nil {
		t.Error("WebSocket connection is not established")
	}

	var msgBytes []byte
	switch v := message.(type) {
	case string:
		msgBytes = []byte(v)
	case []byte:
		msgBytes = v
	default:
		jsonMsg, err := json.Marshal(message)
		if err != nil {
			t.Error(err)
		}
		msgBytes = jsonMsg
	}

	err := c.Conn.WriteMessage(websocket.TextMessage, msgBytes)
	if err != nil {
		t.Error(err)
	}
}

func (c *WSClient) Receive(t *testing.T, timeout time.Duration) ([]byte, error) {
	t.Helper()

	if c.Conn == nil {
		t.Error("WebSocket connection is not established")
	}

	if timeout > 0 {
		c.Conn.SetReadDeadline(time.Now().Add(timeout))
		defer c.Conn.SetReadDeadline(time.Time{})
	}

	_, message, err := c.Conn.ReadMessage()
	return message, err
}

func (c *WSClient) Close(t *testing.T) {
	t.Helper()
	if c.Conn != nil {
		err := c.Conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			t.Error(err)
		}
		c.Conn.Close()
	}
}
