package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"runtime"
	"strconv"

	"sync"
	"time"

	"sorting.com/generics"
)

func main() {

	//	Create(10)

	// The variables to track time and memory spent
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	t0 := time.Now()

	// Data from the file and its length
	data, l := generics.FileScaner("./test.txt")

	t2 := time.Now()
	fmt.Println(t2.Sub(t0))

	generics.QuickSort(data, 0, l-1)

	// Sum is the result of the program's work
	// It shows the number of unique addresses in this file
	sum := 1

	// The gr shows the optimal amount of goroutines based on Amdahl’s law
	gr := min(l/(10^5), 1000)

	switch {
	case gr < 10:
		for j := 0; j < l-1; j++ {
			if data[j] != data[j+1] {
				sum++
			}
		}
	default:
		var wg sync.WaitGroup

		// Chan for tracking the result of goroutine work
		res := make(chan int, gr)
		wg.Add(gr - 1)

		// The loop distributes the array into gr-1 goroutines
		for i := 0; i < gr-1; i++ {
			go func(p int) {
				local := 0
				j := (p * l / gr)
				for ; j < ((p + 1) * l / gr); j++ {
					if data[j] != data[j+1] {
						local++
					}
				}
				res <- local
				wg.Done()
			}(i)
		}

		// The last part of the array
		for j := (gr - 1) * l / gr; j < l-1; j++ {
			if data[j] != data[j+1] {
				sum++
			}
		}
		wg.Wait()

		// Summarizing results from goroutines
		for j := 0; j < gr-1; j++ {
			v := <-res

			sum += v
		}
	}
	t1 := time.Now()

	fmt.Println(t1.Sub(t0))

	fmt.Println(sum)

	fmt.Printf("Total allocated memory (in bytes): %d\n", memStats.Alloc)
	fmt.Printf("Heap memory (in bytes): %d\n", memStats.HeapAlloc)
	fmt.Printf("Number of garbage collections: %d\n", memStats.NumGC)
}

// Create creates the "test.txt" file with n various IPv4 addresses
//
// It has one parameter: an int instance indicating the number of IPv4
func Create(n int) {
	if n < 1 {
		log.Fatal("it is not enough")
	}

	file, err := os.Create("test.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	for i := 0; i < n; i++ {
		text := ""

		for j := 0; j < 4; j++ {
			text += strconv.Itoa(rand.IntN(255)) + "."
		}

		file.WriteString(text[:len(text)-1] + "\n")
	}
}
