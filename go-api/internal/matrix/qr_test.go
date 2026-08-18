package matrix

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

const qrTolerance = 1e-10

func TestFactorizeQR_Properties(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		matrix [][]float64
	}{
		{
			name: "square with negatives",
			matrix: [][]float64{
				{12, -51, 4},
				{6, 167, -68},
				{-4, 24, -41},
			},
		},
		{
			name: "rectangular supported",
			matrix: [][]float64{
				{1, 2},
				{3, 4},
				{5, 6},
			},
		},
		{
			name: "decimals",
			matrix: [][]float64{
				{0.5, -1.25},
				{2.5, 0.75},
			},
		},
		{
			name:   "1x1",
			matrix: [][]float64{{3.5}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			qValues, rValues, err := FactorizeQR(tc.matrix)
			if err != nil {
				t.Fatalf("FactorizeQR returned error: %v", err)
			}

			a := toDense(tc.matrix)
			q := toDense(qValues)
			r := toDense(rValues)

			assertApproxEqual(t, "A ≈ Q×R", mul(q, r), a, qrTolerance)
			assertApproxEqual(t, "Qᵀ×Q ≈ I", mul(q.T(), q), identity(len(qValues)), qrTolerance)
		})
	}
}

func mul(x, y mat.Matrix) *mat.Dense {
	var product mat.Dense
	product.Mul(x, y)
	return &product
}

func identity(n int) *mat.Dense {
	dense := mat.NewDense(n, n, nil)
	for i := 0; i < n; i++ {
		dense.Set(i, i, 1)
	}
	return dense
}

func assertApproxEqual(t *testing.T, label string, got, want mat.Matrix, tol float64) {
	t.Helper()
	if !mat.EqualApprox(got, want, tol) {
		t.Fatalf("%s failed\ngot:\n%v\nwant:\n%v", label, mat.Formatted(got), mat.Formatted(want))
	}
}
