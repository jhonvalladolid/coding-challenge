package statistics

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-api/internal/config"
	apperr "go-api/internal/errors"
	"go-api/internal/response"
)

const validNodeBody = `{
  "success": true,
  "data": {
    "statistics": {"max": 4, "min": 0, "average": 1.375, "sum": 11},
    "diagonal": {"q": true, "r": false, "anyDiagonal": true}
  }
}`

func TestClientCalculate_Success(t *testing.T) {
	t.Parallel()

	var (
		gotMethod      string
		gotPath        string
		gotContentType string
		gotRequestID   string
		gotPayload     Request
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotContentType = r.Header.Get("Content-Type")
		gotRequestID = r.Header.Get(response.RequestIDHeader)
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotPayload)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(validNodeBody))
	}))
	defer server.Close()

	client := NewClient(config.Config{
		AppEnv:               "test",
		StatisticsAPIURL:     server.URL,
		StatisticsAPITimeout: time.Second,
	})

	q := [][]float64{{1, 0}, {0, 1}}
	r := [][]float64{{2, 3}, {0, 4}}
	result, err := client.Calculate(context.Background(), "abc-123", q, r)
	if err != nil {
		t.Fatalf("Calculate: %v", err)
	}

	if gotMethod != http.MethodPost || gotPath != "/api/v1/statistics" {
		t.Fatalf("method=%s path=%s", gotMethod, gotPath)
	}
	if gotContentType != "application/json" {
		t.Fatalf("content-type=%s", gotContentType)
	}
	if gotRequestID != "abc-123" {
		t.Fatalf("request id=%s", gotRequestID)
	}
	if len(gotPayload.Matrices.Q) != 2 || len(gotPayload.Matrices.R) != 2 {
		t.Fatalf("payload=%v", gotPayload)
	}
	if result.Statistics.Max != 4 || result.Diagonal.AnyDiagonal != true {
		t.Fatalf("result=%v", result)
	}
}

func TestClientCalculate_HTTPErrors(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		status     int
		body       string
		wantCode   string
		wantStatus int
	}{
		{name: "400", status: 400, body: `{"success":false}`, wantCode: "STATISTICS_SERVICE_ERROR", wantStatus: 502},
		{name: "500", status: 500, body: `{"success":false}`, wantCode: "STATISTICS_SERVICE_ERROR", wantStatus: 502},
		{name: "invalid json", status: 200, body: `{`, wantCode: "STATISTICS_SERVICE_BAD_RESPONSE", wantStatus: 502},
		{name: "missing data", status: 200, body: `{"success":true}`, wantCode: "STATISTICS_SERVICE_BAD_RESPONSE", wantStatus: 502},
		{name: "success false", status: 200, body: `{"success":false,"data":{"statistics":{"max":1,"min":0,"average":0,"sum":1},"diagonal":{"q":false,"r":false,"anyDiagonal":false}}}`, wantCode: "STATISTICS_SERVICE_BAD_RESPONSE", wantStatus: 502},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			t.Cleanup(server.Close)

			client := NewClient(config.Config{
				AppEnv:               "test",
				StatisticsAPIURL:     server.URL,
				StatisticsAPITimeout: time.Second,
			})

			_, err := client.Calculate(context.Background(), "rid", [][]float64{{1}}, [][]float64{{1}})
			appErr, ok := apperr.AsAppError(err)
			if !ok {
				t.Fatalf("expected AppError, got %v", err)
			}
			if appErr.Code != tc.wantCode || appErr.StatusCode != tc.wantStatus {
				t.Fatalf("code=%s status=%d", appErr.Code, appErr.StatusCode)
			}
		})
	}
}

func TestClientCalculate_Timeout(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(validNodeBody))
	}))
	defer server.Close()

	client := NewClient(config.Config{
		AppEnv:               "test",
		StatisticsAPIURL:     server.URL,
		StatisticsAPITimeout: 50 * time.Millisecond,
	})

	_, err := client.Calculate(context.Background(), "rid", [][]float64{{1}}, [][]float64{{1}})
	appErr, ok := apperr.AsAppError(err)
	if !ok || appErr.Code != "STATISTICS_SERVICE_TIMEOUT" || appErr.StatusCode != 504 {
		t.Fatalf("expected timeout, got %v", err)
	}
}

func TestClientCalculate_Unavailable(t *testing.T) {
	t.Parallel()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := listener.Addr().String()
	_ = listener.Close()

	client := NewClient(config.Config{
		AppEnv:               "test",
		StatisticsAPIURL:     "http://" + addr,
		StatisticsAPITimeout: 200 * time.Millisecond,
	})

	_, err = client.Calculate(context.Background(), "rid", [][]float64{{1}}, [][]float64{{1}})
	appErr, ok := apperr.AsAppError(err)
	if !ok || appErr.Code != "STATISTICS_SERVICE_UNAVAILABLE" || appErr.StatusCode != 503 {
		t.Fatalf("expected unavailable, got %v", err)
	}
}
