package matrix

import (
	"math"
	"testing"

	apperr "go-api/internal/errors"
)

func TestValidate_ValidMatrices(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		matrix [][]float64
	}{
		{name: "1x1", matrix: [][]float64{{7}}},
		{name: "square", matrix: [][]float64{{1, 2}, {3, 4}}},
		{name: "rectangular m>n", matrix: [][]float64{{1, 2}, {3, 4}, {5, 6}}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if err := Validate(tc.matrix, 200); err != nil {
				t.Fatalf("expected valid matrix, got %v", err)
			}
		})
	}
}

func TestValidate_InvalidMatrices(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		matrix [][]float64
		code   string
	}{
		{name: "nil", matrix: nil, code: "MATRIX_REQUIRED"},
		{name: "empty", matrix: [][]float64{}, code: "EMPTY_MATRIX"},
		{name: "empty row", matrix: [][]float64{{1, 2}, {}}, code: "EMPTY_MATRIX_ROW"},
		{name: "irregular", matrix: [][]float64{{1, 2}, {3}}, code: "IRREGULAR_MATRIX"},
		{name: "nan", matrix: [][]float64{{math.NaN()}}, code: "INVALID_MATRIX_VALUE"},
		{name: "plus inf", matrix: [][]float64{{math.Inf(1)}}, code: "INVALID_MATRIX_VALUE"},
		{name: "minus inf", matrix: [][]float64{{math.Inf(-1)}}, code: "INVALID_MATRIX_VALUE"},
		{name: "m < n", matrix: [][]float64{{1, 2, 3}, {4, 5, 6}}, code: "UNSUPPORTED_MATRIX_DIMENSIONS"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := Validate(tc.matrix, 200)
			appErr, ok := apperr.AsAppError(err)
			if !ok {
				t.Fatalf("expected AppError, got %v", err)
			}
			if appErr.Code != tc.code {
				t.Fatalf("expected code %s, got %s", tc.code, appErr.Code)
			}
		})
	}
}
