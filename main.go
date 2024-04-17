package main

import (
	"fmt"
	"log"
	"os"
	"sort"
)

type WeatherStation struct {
	NumVal int
	SumVal float64
	MinVal float64
	MaxVal float64
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("invalid number of arguments\nusage: %s <version>(v1,v2,etc) <filename>", os.Args[0])
	}
	version := os.Args[1]
	filename := os.Args[2]

	// Select 1brc version
	var vFunc func(string) map[string]*WeatherStation
	switch version {
	case "v1":
		vFunc = v1
	case "v2":
		vFunc = v2
	default:
		log.Fatalf("no version found for '%s'", version)
	}

	// Do work
    M := vFunc(filename)
	if len(M) == 0 {
		log.Fatal("map empty")
	}

	// Sort output
	order := make([]string, 0, len(M))
	for name, _ := range M {
		order = append(order, name)
	}
	sort.Strings(order)

	// Print results to stdout
	fmt.Printf("{")
	out := ""
	for _, name := range order {
		w := M[name]
		out += fmt.Sprintf("%s=%.1f/%.1f/%.1f, ",
			name,
			w.MinVal,
			w.SumVal/float64(w.NumVal),
			w.MaxVal,
		)
	}
	fmt.Printf("%s}\n", out[:len(out)-2])
}
