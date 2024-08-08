// Package naive provides a "naive" algorithm for solving the test problem
package naive

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// InterAddr represents the type of converted string
type InetAddr uint32

// Result represents the HashSet
var result = make(map[InetAddr]bool)

// FilterScaner scans the file, fills the HashSet with strings converted to uint32
// and counts unique values
//
// It has one parameter: a string with the value of the address of the file being scanned
func FileScaner(name string) int {

	file, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		var ipUint InetAddr
		v := scanner.Text()

		data := strings.Split(v, ".")

		// Strings from data array converts to uint32 and adds to HashSet
		for i := 0; i < 3; i++ {
			val, _ := strconv.Atoi(data[i])

			ipUint += InetAddr(val)
			ipUint = ipUint << 8
		}
		val, _ := strconv.Atoi(data[3])
		ipUint += InetAddr(val)

		if _, ok := result[ipUint]; !ok {
			result[ipUint] = true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	// Return the number of unique addresses in this file
	return len(result)
}
