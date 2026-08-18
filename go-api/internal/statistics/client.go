package statistics

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"go-api/internal/config"
	apperr "go-api/internal/errors"
	"go-api/internal/response"
)

const statisticsPath = "/api/v1/statistics"
const maxResponseBytes = 2 << 20

type Client struct {
	baseURL    string
	httpClient *http.Client
	logger     *slog.Logger
}

func NewClient(cfg config.Config) *Client {
	logger := slog.Default()
	if cfg.AppEnv == "test" {
		logger = slog.New(slog.DiscardHandler)
	}

	return &Client{
		baseURL: strings.TrimRight(cfg.StatisticsAPIURL, "/"),
		httpClient: &http.Client{
			Timeout: cfg.StatisticsAPITimeout,
		},
		logger: logger,
	}
}

func (c *Client) Calculate(ctx context.Context, requestID string, q, r [][]float64) (*Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	payload, err := json.Marshal(Request{
		Matrices: Matrices{Q: q, R: r},
	})
	if err != nil {
		return nil, apperr.InternalServerError()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+statisticsPath, bytes.NewReader(payload))
	if err != nil {
		return nil, apperr.InternalServerError()
	}

	req.Header.Set("Content-Type", "application/json")
	if requestID != "" {
		req.Header.Set(response.RequestIDHeader, requestID)
	}

	started := time.Now()
	resp, err := c.httpClient.Do(req)
	if err != nil {
		classified := classifyTransportError(err)
		c.logger.Error("statistics request failed",
			"requestId", requestID,
			"service", "node-statistics",
			"code", classified.Code,
			"durationMs", time.Since(started).Milliseconds(),
		)
		return nil, classified
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	if err != nil {
		return nil, apperr.StatisticsServiceBadResponse()
	}

	c.logger.Info("statistics request completed",
		"requestId", requestID,
		"service", "node-statistics",
		"status", resp.StatusCode,
		"durationMs", time.Since(started).Milliseconds(),
	)

	if resp.StatusCode >= 400 {
		return nil, apperr.StatisticsServiceError()
	}

	return parseNodeResponse(raw)
}

func parseNodeResponse(raw []byte) (*Result, error) {
	var envelope nodeResponse
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, apperr.StatisticsServiceBadResponse()
	}

	if !envelope.Success || envelope.Data == nil || envelope.Data.Statistics == nil || envelope.Data.Diagonal == nil {
		return nil, apperr.StatisticsServiceBadResponse()
	}

	return &Result{
		Statistics: *envelope.Data.Statistics,
		Diagonal:   *envelope.Data.Diagonal,
	}, nil
}

func classifyTransportError(err error) *apperr.AppError {
	if errors.Is(err, context.DeadlineExceeded) {
		return apperr.StatisticsServiceTimeout()
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return apperr.StatisticsServiceTimeout()
	}

	return apperr.StatisticsServiceUnavailable()
}
