package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

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
	f, err := os.Open(FILENAME)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	M := make(map[string]*WeatherStation)

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ";")
		name := s[0]
		val, err := strconv.ParseFloat(s[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		if w, ok := M[name]; ok {
			w.NumVal++
			w.SumVal += val
			w.MinVal = min(w.MinVal, val)
			w.MaxVal = max(w.MaxVal, val)
		} else {
			M[name] = &WeatherStation{
				NumVal: 1,
				SumVal: val,
				MinVal: val,
				MaxVal: val,
			}
		}
	}

	fmt.Printf("{")
	out := ""
	for name, w := range M {
		out += fmt.Sprintf("%s=%.1f/%.1f/%.1f, ",
			name,
			w.MinVal,
			w.SumVal/float64(w.NumVal),
			w.MaxVal,
		)
	}
	fmt.Printf("%s}\n", out[:len(out)-2])
}
