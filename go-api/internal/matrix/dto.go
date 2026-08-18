package matrix

type QRRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

type Factorization struct {
	Q [][]float64 `json:"q"`
	R [][]float64 `json:"r"`
}

type QRResponseData struct {
	OriginalMatrix [][]float64   `json:"originalMatrix"`
	Factorization  Factorization `json:"factorization"`
}

type HealthData struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}
