package matrix

import (
	apperr "go-api/internal/errors"
	"gonum.org/v1/gonum/mat"
)

func FactorizeQR(matrix [][]float64) (q [][]float64, r [][]float64, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = apperr.QRFactorizationError()
		}
	}()

	dense := toDense(matrix)

	var qr mat.QR
	qr.Factorize(dense)

	var qMat, rMat mat.Dense
	qr.QTo(&qMat)
	qr.RTo(&rMat)

	return fromDense(&qMat), fromDense(&rMat), nil
}

func toDense(matrix [][]float64) *mat.Dense {
	rows := len(matrix)
	cols := len(matrix[0])
	data := make([]float64, 0, rows*cols)
	for _, row := range matrix {
		data = append(data, row...)
	}
	return mat.NewDense(rows, cols, data)
}

func fromDense(dense *mat.Dense) [][]float64 {
	rows, cols := dense.Dims()
	out := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		out[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			out[i][j] = dense.At(i, j)
		}
	}
	return out
}
