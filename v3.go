package main

import (
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
)

const (
	NUM_WORKERS = 32
	BUF_SIZE    = 100
)

var wg sync.WaitGroup

func v3(filename string) map[string]*WeatherStation {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	M := make(map[string]*WeatherStation)

	// Create concurrent workers
	inputChan := make(chan []byte, BUF_SIZE)
	outputChan := make(chan map[string]*WeatherStation, NUM_WORKERS)
	for i := 0; i < NUM_WORKERS; i++ {
		wg.Add(1)
		go worker(inputChan, outputChan)
	}

	// Read from file into 4K buffers for optimal read speeds
	buffer := make([]byte, 2048*2048)
	leftoverBuffer := make([]byte, 256)
	leftoverSize := 0
	for {
		n, err := f.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatal(err)
		}

		// Find last newline character in buffer
		newline_idx := n - 1
		for newline_idx >= 0 {
			if buffer[newline_idx] == '\n' {
				break
			}
			newline_idx--
		}

		// Only process to the last newline. The remaining data is a partial
		// record that needs to be processed in the next iteration. Save it in
		// a temporary buffer.
		// If there was any leftover data from the previous iteration, prepend
		// it to the current buffer.
		to_process := make([]byte, newline_idx+1+leftoverSize)

		copy(to_process, leftoverBuffer[:leftoverSize])
		copy(to_process[leftoverSize:], buffer[:newline_idx+1])
		copy(leftoverBuffer, buffer[newline_idx+1:n])

		leftoverSize = n - (newline_idx + 1)

		// Perform processing operation
		inputChan <- to_process
	}

	// Done input. Wait for all workers to finish
	close(inputChan)
	wg.Wait()
	close(outputChan)

	// Merge partial results into one map
	for result := range outputChan {
		for name, ws := range result {
			if ws2, ok := M[name]; ok {
				ws2.NumVal += ws.NumVal
				ws2.SumVal += ws.SumVal
				ws2.MinVal = min(ws2.MinVal, ws.MinVal)
				ws2.MaxVal = max(ws2.MaxVal, ws.MaxVal)
			} else {
				M[name] = ws
			}
		}
	}

	return M
}

func worker(input chan []byte, output chan map[string]*WeatherStation) {
	defer wg.Done()
	result := make(map[string]*WeatherStation)

	for {
		// Get buffer from the input channel
		buffer, more := <-input
		if !more {
			break
		}

		// Parse the buffer
		idx_start := 0
		for idx_start < len(buffer) {
			// Read string up to ';' character
			idx_end := idx_start + 1
			for buffer[idx_end] != ';' {
				idx_end++
			}
			name := string(buffer[idx_start:idx_end])
			idx_start = idx_end + 1

			// Read float up to '\n' character
			idx_end = idx_start + 1
			for buffer[idx_end] != '\n' {
				idx_end++
			}
			val, err := strconv.ParseFloat(string(buffer[idx_start:idx_end]), 64)
			if err != nil {
				log.Fatal(err)
			}
			idx_start = idx_end + 1

			// Add to map
			if w, ok := result[name]; ok {
				w.NumVal++
				w.SumVal += val
				w.MinVal = min(w.MinVal, val)
				w.MaxVal = max(w.MaxVal, val)
			} else {
				result[name] = &WeatherStation{
					NumVal: 1,
					SumVal: val,
					MinVal: val,
					MaxVal: val,
				}
			}
		}
	}

	output <- result
}
