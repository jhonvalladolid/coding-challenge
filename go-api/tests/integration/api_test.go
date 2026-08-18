package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-api/internal/config"
	"go-api/internal/matrix"
)

func doRequest(t *testing.T, method, path, body, requestID string) (int, http.Header, map[string]any) {
	t.Helper()

	var reader io.Reader
	if body != "" {
		reader = bytes.NewBufferString(body)
	}

	req := httptest.NewRequest(method, path, reader)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if requestID != "" {
		req.Header.Set("X-Request-ID", requestID)
	}

	resp, err := matrix.NewApp(config.Config{AppEnv: "test", MaxMatrixDim: 200}).Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	var payload map[string]any
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &payload); err != nil {
			t.Fatalf("unmarshal response: %v\nbody: %s", err, raw)
		}
	}

	return resp.StatusCode, resp.Header, payload
}

func TestHealth(t *testing.T) {
	status, header, payload := doRequest(t, http.MethodGet, "/health", "", "")

	if status != http.StatusOK {
		t.Fatalf("status=%d", status)
	}
	if header.Get("X-Request-ID") == "" {
		t.Fatal("expected generated X-Request-ID")
	}
	if payload["success"] != true {
		t.Fatalf("success=%v", payload["success"])
	}

	data := payload["data"].(map[string]any)
	if data["service"] != "go-api" || data["status"] != "ok" {
		t.Fatalf("data=%v", data)
	}
}

func TestHealthReusesRequestID(t *testing.T) {
	status, header, payload := doRequest(t, http.MethodGet, "/health", "", "test-123")
	if status != http.StatusOK {
		t.Fatalf("status=%d", status)
	}
	if header.Get("X-Request-ID") != "test-123" {
		t.Fatalf("header=%s", header.Get("X-Request-ID"))
	}

	meta := payload["meta"].(map[string]any)
	if meta["requestId"] != "test-123" {
		t.Fatalf("meta=%v", meta)
	}
}

func TestFactorizeSuccess(t *testing.T) {
	body := `{
		"matrix": [
			[12, -51, 4],
			[6, 167, -68],
			[-4, 24, -41]
		]
	}`

	status, header, payload := doRequest(t, http.MethodPost, "/api/v1/matrices/qr", body, "")
	if status != http.StatusOK {
		t.Fatalf("status=%d payload=%v", status, payload)
	}
	if payload["success"] != true {
		t.Fatalf("payload=%v", payload)
	}
	if header.Get("X-Request-ID") == "" {
		t.Fatal("expected X-Request-ID")
	}

	data := payload["data"].(map[string]any)
	if _, ok := data["originalMatrix"]; !ok {
		t.Fatal("missing originalMatrix")
	}
	factorization := data["factorization"].(map[string]any)
	if factorization["q"] == nil || factorization["r"] == nil {
		t.Fatalf("factorization=%v", factorization)
	}
}

func TestFactorizeValidationErrors(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		body   string
		status int
		code   string
	}{
		{name: "empty body", body: "", status: 400, code: "INVALID_REQUEST_BODY"},
		{name: "malformed json", body: "{", status: 400, code: "INVALID_REQUEST_BODY"},
		{name: "matrix missing", body: `{}`, status: 400, code: "MATRIX_REQUIRED"},
		{name: "empty matrix", body: `{"matrix":[]}`, status: 400, code: "EMPTY_MATRIX"},
		{name: "empty row", body: `{"matrix":[[1,2],[]]}`, status: 400, code: "EMPTY_MATRIX_ROW"},
		{name: "irregular", body: `{"matrix":[[1,2],[3]]}`, status: 400, code: "IRREGULAR_MATRIX"},
		{name: "unsupported dimensions", body: `{"matrix":[[1,2,3],[4,5,6]]}`, status: 422, code: "UNSUPPORTED_MATRIX_DIMENSIONS"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, header, payload := doRequest(t, http.MethodPost, "/api/v1/matrices/qr", tc.body, "err-1")
			if status != tc.status {
				t.Fatalf("status=%d payload=%v", status, payload)
			}
			if payload["success"] != false {
				t.Fatalf("success=%v", payload["success"])
			}
			errBody := payload["error"].(map[string]any)
			if errBody["code"] != tc.code {
				t.Fatalf("code=%v payload=%v", errBody["code"], payload)
			}
			if header.Get("X-Request-ID") != "err-1" {
				t.Fatalf("header=%s", header.Get("X-Request-ID"))
			}
			meta := payload["meta"].(map[string]any)
			if meta["requestId"] != "err-1" {
				t.Fatalf("meta=%v", meta)
			}
		})
	}
}
