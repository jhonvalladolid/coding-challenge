package matrix

type Service struct {
	maxMatrixDim int
}

func NewService(maxMatrixDim int) *Service {
	return &Service{maxMatrixDim: maxMatrixDim}
}

func (s *Service) Factorize(matrix [][]float64) (*QRResponseData, error) {
	if err := Validate(matrix, s.maxMatrixDim); err != nil {
		return nil, err
	}

	q, r, err := FactorizeQR(matrix)
	if err != nil {
		return nil, err
	}

	return &QRResponseData{
		OriginalMatrix: matrix,
		Factorization: Factorization{
			Q: q,
			R: r,
		},
	}, nil
}
