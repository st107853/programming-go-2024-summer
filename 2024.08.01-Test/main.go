package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"test.com/shard"
)

func main() {

	//	Create(1_000_000)

	// The variables to track time and memory spent
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	t0 := time.Now()

	// Res is the result of the program's work
	// It shows the number of unique addresses in this file
	res := FileScaner("./test.txt")
	fmt.Printf("we got result: %v\n", res)

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
	file, err := os.Create("test2.txt")

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

// FilterScaner scans the file, fills the HashSet with strings converted to uint32
// and counts unique values
//
// It has one parameter: a string with the value of the address of the file being scanned
func FileScaner(name string) int {
	var in = make(chan uint32, 1000)

	file, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)

	// Result represents the HashSet
	n := int(min(255, fi.Size()/14275270+1))

	go func() {
		for scanner.Scan() {

			var ipUint uint32
			v := scanner.Text()

			data := strings.Split(v, ".")

			// Strings from data array converts to uint32 and adds to HashSet
			for i := 0; i < 3; i++ {
				val, _ := strconv.Atoi(data[i])

				ipUint += uint32(val)
				ipUint = ipUint << 8
			}
			val, _ := strconv.Atoi(data[3])
			ipUint += uint32(val)

			in <- ipUint
		}
		close(in)
	}()

	cnt := Split(in, n)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	// Return the number of unique addresses in this file
	return cnt
}

func Split(source <-chan uint32, n int) int {
	fmt.Println(n)
	var wG sync.WaitGroup // Use WaitGroup to wait until
	wG.Add(n * n)
	var result = shard.NewShardedMap(n)
	// Create the dests slice
	// Create n destination channels

	res := make(chan int, 100)
	cnt := 0

	for i := 0; i < n*n; i++ {
		go func() {
			cnt := 0
			for val := range source {
				if ok := result.Get(val); !ok {
					result.Set(val, true)
					cnt++

				}
			}
			res <- cnt
			wG.Done()
		}()
	}
	wG.Wait()

	for i := 0; i < n*n; i++ {
		cnt += <-res
	}

	return cnt
}
