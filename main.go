package main

const (
	DATA1M = "/usr/share/1brc-data/measurements-1M.txt"
	DATA1B = "/usr/share/1brc-data/measurements-1B.txt"
)

var FILENAME = DATA1B

type WeatherStation struct {
	NumVal int
	SumVal float64
	MinVal float64
	MaxVal float64
}

func main() {
	// v1(FILENAME)
	v2(FILENAME)
}
