package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"runtime"
	"strconv"
	"time"

	"algorithm.com/naive"
)

func main() {

	//	Create(10_000_000)

	// The variables to track time and memory spent
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	t0 := time.Now()

	// Res is the result of the program's work
	// It shows the number of unique addresses in this file
	res := naive.FileScaner("./test.txt")

	fmt.Println(res)

	t1 := time.Now()

	fmt.Println(t1.Sub(t0))

	fmt.Printf("Total allocated memory (in bytes): %d\n", memStats.Alloc)
	fmt.Printf("Heap memory (in bytes): %d\n", memStats.HeapAlloc)
	fmt.Printf("Number of garbage collections: %d\n", memStats.NumGC)
}

// Create creates the "test.txt" file with n various IPv4 addresses
//
// It has one parameter: an int instance indicating the number of IPv4
func Create(n int) {
	file, err := os.Create("test.txt")

	if err != nil {
		fmt.Println("Unable to create")
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
