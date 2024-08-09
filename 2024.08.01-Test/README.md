Comparison of approaches to solving the problem.

The following solution is presented on the main branch.
At first the FileScanner function scans a file and returns
an array of uint32 numbers.Converting IPv4 strings to uint32
allows you to allocate less memory.

Conversion works as follows:
The string is split into a dataSplit array of four lines.
Array elements are converted into a number and added to the
ipUint variable.
Then ipUint is shifted 8 bits to the left, freeing up space
for writing the next element from the dataSplit array.

Next, the resulting array is quickly sorted. This will
simplify the calculation of unique IPv4s. Sorting with Go
generics is faster than existing ones in the sort package.
This saves time. If there are more than 10^6 elements in the
file, goroutines are launched. They help calculate the
number of unique values in the sorted array. The number of
goroutines launched is calculated based on Amdahl's law.

The table below shows the results of the time spent on the
operation of the described algorithm and the naive algorithm
from the simplification branch.


| amount   |	main	| simplification   |
|----------|--------|------------------|
| 10^4     | ~10ms  |    ~5ms 	       |
| 10^5     | ~80ms  |   ~35ms 	       |
| 10^6     | ~800ms | ~400ms 	         |
|  10^7    |  6-7s  |  5-6s 	         |

In all cases, from 101 to 104 kilobytes of memory were used.
