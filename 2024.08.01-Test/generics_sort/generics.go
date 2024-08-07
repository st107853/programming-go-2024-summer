// Package generics provides utilites to scan, convert and sort files
package generics

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/exp/constraints"
)

// InterAddr represents the type of converted string
type InetAddr uint32

// Partition portions the array
func Partition[T constraints.Ordered](arr []T, low, high int) ([]T, int) {
	pivot := arr[high]
	i := low
	for j := low; j < high; j++ {
		if arr[j] < pivot {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}
	arr[i], arr[high] = arr[high], arr[i]
	return arr, i
}

// QuickSort implements quicksort of generics
func QuickSort[T constraints.Ordered](arr []T, low, high int) []T {
	var wg sync.WaitGroup
	wg.Add(2)

	if low < high {
		var p int
		arr, p = Partition(arr, low, high)
		go func() {
			defer wg.Done()
			arr = QuickSort(arr, low, p-1)
		}()
		go func() {
			defer wg.Done()
			arr = QuickSort(arr, p+1, high)
		}()

		wg.Wait()
	}

	return arr
}

// FielScaner scans the fiel with IPv4 and converts them to uint by byte shifting
//
// It has one parameter: a string with the value of the address of the file being scanned
func FileScaner(name string) ([]InetAddr, int) {

	var result []InetAddr
	var cnt int = 0

	file, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var ipUint InetAddr = 0
		cnt++
		data := scanner.Text()
		dataSplit := strings.Split(data, ".")

		for i := 0; i < 3; i++ {
			val, _ := strconv.Atoi(dataSplit[i])

			ipUint += InetAddr(val)

			ipUint = ipUint << 8
		}

		val, _ := strconv.Atoi(dataSplit[3])

		ipUint += InetAddr(val)

		result = append(result, ipUint)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return result, cnt
}
