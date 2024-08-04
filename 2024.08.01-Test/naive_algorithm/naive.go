package naive

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type InetAddr uint32

var result = make(map[InetAddr]bool)

func FileScaner(name string) int {

	file, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K

	for scanner.Scan() {

		var ipUint InetAddr
		v := scanner.Text()

		data := strings.Split(v, ".")

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

	return len(result)
}
