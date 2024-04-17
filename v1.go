package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func v1(filename string) map[string]*WeatherStation {
	f, err := os.Open(filename)
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

	return M
}
