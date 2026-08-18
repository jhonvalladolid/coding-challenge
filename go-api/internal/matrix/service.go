package matrix

import (
	"context"

	"go-api/internal/statistics"
)

type StatisticsProvider interface {
	Calculate(ctx context.Context, requestID string, q, r [][]float64) (*statistics.Result, error)
}

type Service struct {
	maxMatrixDim int
	stats        StatisticsProvider
}

func NewService(maxMatrixDim int, stats StatisticsProvider) *Service {
	return &Service{
		maxMatrixDim: maxMatrixDim,
		stats:        stats,
	}
}

func (s *Service) ProcessQR(ctx context.Context, requestID string, matrix [][]float64) (*QRResponseData, error) {
	if err := Validate(matrix, s.maxMatrixDim); err != nil {
		return nil, err
	}

	q, r, err := FactorizeQR(matrix)
	if err != nil {
		return nil, err
	}

	stats, err := s.stats.Calculate(ctx, requestID, q, r)
	if err != nil {
		return nil, err
	}

	return &QRResponseData{
		OriginalMatrix: matrix,
		Factorization: Factorization{
			Q: q,
			R: r,
		},
		Statistics: stats.Statistics,
		Diagonal:   stats.Diagonal,
	}, nil
}
