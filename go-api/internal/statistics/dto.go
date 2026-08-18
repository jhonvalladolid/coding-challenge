package statistics

type Matrices struct {
	Q [][]float64 `json:"q"`
	R [][]float64 `json:"r"`
}

type Request struct {
	Matrices Matrices `json:"matrices"`
}

type Statistics struct {
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
	Average float64 `json:"average"`
	Sum     float64 `json:"sum"`
}

type Diagonal struct {
	Q           bool `json:"q"`
	R           bool `json:"r"`
	AnyDiagonal bool `json:"anyDiagonal"`
}

type Result struct {
	Statistics Statistics `json:"statistics"`
	Diagonal   Diagonal   `json:"diagonal"`
}

type nodeData struct {
	Statistics *Statistics `json:"statistics"`
	Diagonal   *Diagonal   `json:"diagonal"`
}

type nodeResponse struct {
	Success bool      `json:"success"`
	Data    *nodeData `json:"data"`
}
