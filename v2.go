package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const V2_NUM_WORKERS = 32

func v2(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Start up workers
	line_chan := make(chan string)
	results_chan := make(chan map[string]*WeatherStation)
	for i := 0; i < V2_NUM_WORKERS; i++ {
		go work(line_chan, results_chan)
	}

	// Fan out lines
	scanner := bufio.NewScanner(f)
	go func() {
		defer close(line_chan)
		for scanner.Scan() {
			line_chan <- scanner.Text()
		}
	}()

	merged := make(map[string]*WeatherStation)
	for i := 0; i < V2_NUM_WORKERS; i++ {
		result := <-results_chan
		for name, ws := range result {
			w := merged[name]
			if w != nil {
				w.NumVal += ws.NumVal
				w.SumVal += ws.SumVal
				w.MinVal = min(w.MinVal, ws.MinVal)
				w.MaxVal = max(w.MaxVal, ws.MaxVal)
			} else {
				merged[name] = ws
			}
		}
	}

	fmt.Printf("{")
	out := ""
	for name, w := range merged {
		out += fmt.Sprintf("%s=%.1f/%.1f/%.1f, ",
			name,
			w.MinVal,
			w.SumVal/float64(w.NumVal),
			w.MaxVal,
		)
	}
	fmt.Printf("%s}\n", out[:len(out)-2])
}

func work(line_chan <-chan string, results_chan chan<- map[string]*WeatherStation) {
	results := make(map[string]*WeatherStation)
	for line := range line_chan {
		s := strings.Split(line, ";")
		name := s[0]
		val, err := strconv.ParseFloat(s[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		w := results[name]
		if w != nil {
			w.NumVal++
			w.SumVal += val
			w.MinVal = min(w.MinVal, val)
			w.MaxVal = max(w.MaxVal, val)
		} else {
			results[name] = &WeatherStation{
				NumVal: 1,
				SumVal: val,
				MinVal: val,
				MaxVal: val,
			}
		}
	}

	results_chan <- results
}
