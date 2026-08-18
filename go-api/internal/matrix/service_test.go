package matrix

import (
	"context"
	"errors"
	"testing"

	apperr "go-api/internal/errors"
	"go-api/internal/statistics"
)

type stubStatistics struct {
	requestID string
	q         [][]float64
	r         [][]float64
	result    *statistics.Result
	err       error
}

func (s *stubStatistics) Calculate(_ context.Context, requestID string, q, r [][]float64) (*statistics.Result, error) {
	s.requestID = requestID
	s.q = q
	s.r = r
	return s.result, s.err
}

func TestProcessQR_Success(t *testing.T) {
	t.Parallel()

	stub := &stubStatistics{
		result: &statistics.Result{
			Statistics: statistics.Statistics{Max: 4, Min: 0, Average: 1.375, Sum: 11},
			Diagonal:   statistics.Diagonal{Q: true, R: false, AnyDiagonal: true},
		},
	}
	service := NewService(200, stub)

	out, err := service.ProcessQR(context.Background(), "rid-1", [][]float64{
		{1, 0},
		{0, 1},
	})
	if err != nil {
		t.Fatalf("ProcessQR: %v", err)
	}

	if stub.requestID != "rid-1" {
		t.Fatalf("requestID=%s", stub.requestID)
	}
	if len(stub.q) == 0 || len(stub.r) == 0 {
		t.Fatal("expected Q and R to be sent to statistics")
	}
	if out.Factorization.Q == nil || out.Factorization.R == nil {
		t.Fatal("expected Q and R in the result")
	}
	if out.Statistics != stub.result.Statistics {
		t.Fatalf("statistics=%v", out.Statistics)
	}
	if out.Diagonal != stub.result.Diagonal {
		t.Fatalf("diagonal=%v", out.Diagonal)
	}
}

func TestProcessQR_StatisticsError(t *testing.T) {
	t.Parallel()

	stub := &stubStatistics{err: apperr.StatisticsServiceUnavailable()}
	service := NewService(200, stub)

	_, err := service.ProcessQR(context.Background(), "rid-2", [][]float64{{1}})
	if !errors.Is(err, stub.err) {
		appErr, ok := apperr.AsAppError(err)
		if !ok || appErr.Code != "STATISTICS_SERVICE_UNAVAILABLE" {
			t.Fatalf("expected statistics unavailable, got %v", err)
		}
	}
}

func TestProcessQR_ValidationDoesNotCallStatistics(t *testing.T) {
	t.Parallel()

	stub := &stubStatistics{err: errors.New("should not be called")}
	service := NewService(200, stub)

	_, err := service.ProcessQR(context.Background(), "rid-3", nil)
	appErr, ok := apperr.AsAppError(err)
	if !ok || appErr.Code != "MATRIX_REQUIRED" {
		t.Fatalf("expected MATRIX_REQUIRED, got %v", err)
	}
	if stub.requestID != "" {
		t.Fatal("statistics client should not be called for invalid matrices")
	}
}
