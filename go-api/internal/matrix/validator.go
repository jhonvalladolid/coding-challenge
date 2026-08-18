package matrix

import (
	"math"

	apperr "go-api/internal/errors"
)

func Validate(matrix [][]float64, maxDim int) error {
	if matrix == nil {
		return apperr.MatrixRequired()
	}

	if len(matrix) == 0 {
		return apperr.EmptyMatrix()
	}

	columnCount := -1
	for _, row := range matrix {
		if len(row) == 0 {
			return apperr.EmptyMatrixRow()
		}

		if columnCount == -1 {
			columnCount = len(row)
		} else if len(row) != columnCount {
			return apperr.IrregularMatrix()
		}

		for _, value := range row {
			if math.IsNaN(value) || math.IsInf(value, 0) {
				return apperr.InvalidMatrixValue()
			}
		}
	}

	rowCount := len(matrix)
	if rowCount > maxDim || columnCount > maxDim {
		return apperr.UnsupportedMatrixDimensions("The matrix dimensions exceed the maximum allowed.")
	}

	// gonum.org/v1/gonum/mat.QR.Factorize panics unless m >= n.
	if rowCount < columnCount {
		return apperr.UnsupportedMatrixDimensions(
			"QR factorization requires a matrix with at least as many rows as columns (m >= n).",
		)
	}

	return nil
}
