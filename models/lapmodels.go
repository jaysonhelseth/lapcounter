package models

type LapModel struct {
	LapCount int64
	Distance float64
}

var current *LapModel
var UpdateCh = make(chan int64)

func DefaultLapModel() *LapModel {
	if current != nil {
		return current
	}

	current = &LapModel{
		LapCount: 0,
		Distance: 0.0,
	}
	return current
}


