package util

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func ReadLines(input io.Reader) ([]string, error) {
	lines := make([]string, 0, 1024)
	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func ReadLinesFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ReadLines(file)
}

func ReadNumbers(input io.Reader) ([]int, error) {
	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanWords)
	data := make([]int, 0, 1024)
	for scanner.Scan() {
		value, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return data, err
		}
		data = append(data, value)
	}
	return data, scanner.Err()
}

func ReadNumbersFromFile(path string) ([]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ReadNumbers(file)
}

func PrintArray(array []int) {
	for idx, val := range array {
		fmt.Printf("  [%2d] %d\n", idx, val)
	}
}

type Set map[int]bool

func (s Set) String() string {
	var ss strings.Builder
	ss.WriteRune('{')
	count := 0
	for value, _ := range s {
		ss.WriteString(strconv.Itoa(value))
		if count != len(s)-1 {
			ss.WriteString(", ")
		}
		count++
	}
	ss.WriteRune('}')
	return ss.String()
}

func first_not_of(haystack []byte, hay byte) int {
	for idx, c := range haystack {
		if c != hay {
			return idx // needle
		}
	}
	return -1
}

// SplitFunc for a bufio.Scanner which splits input into groups of lines
// separated by groups of blank lines.
func ScanLineGroups(data []byte, atEOF bool) (advance int, token []byte, err error) {
	var (
		consecutive int
		start       int
		end         int
		c           byte
	)

	// Capture until we see consecutive empty lines.
	for advance, c = range data {
		if c == '\n' {
			if consecutive == 0 {
				end = advance
			}
			consecutive++
		} else {
			if consecutive > 1 {
				break
			}
			consecutive = 0
			end = start
		}
	}

	// Didn't find a complete token, expand the buffer.
	if consecutive < 2 && !atEOF {
		return 0, nil, nil
	}

	// Found a token (maybe).
	if end > start {
		token = data[start:end]
	}

	// Eat trailing newlines.
	next := first_not_of(data[advance:], '\n')
	if next < 0 {
		advance = len(data)
	} else {
		advance += next
	}

	return advance, token, err
}

func ScanInput(input io.Reader, scan bufio.SplitFunc) ([]string, error) {
	scanner := bufio.NewScanner(input)
	scanner.Split(scan)
	result := make([]string, 0)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return result, scanner.Err()
}

func ReadLineGroups(input io.Reader) ([]string, error) {
	return ScanInput(input, ScanLineGroups)
}

func ReadLineGroupsFromFile(path string) ([]string, error) {
	result, err := ReadFile(path, func(input io.Reader) (interface{}, error) {
		return ReadLineGroups(input)
	})
	if groups, ok := result.([]string); ok {
		return groups, err
	}
	return nil, err
}

func ReadFile(path string, read func(input io.Reader) (interface{}, error)) (
	interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return read(file)
}

func FieldsToInts(strings []string) (ints []int, err error) {
	ints = make([]int, len(strings))
	for index, str := range strings {
		ints[index], err = strconv.Atoi(str)
		if err != nil {
			break
		}
	}
	return ints, err
}

func Product(numbers []int) int {
	result := numbers[0]
	for _, value := range numbers[1:] {
		result *= value
	}
	return result
}

func IPow(root int, exp int) int {
	root = 1
	for n := 0; n < exp; n++ {
		root *= exp
	}
	return root
}

func IAbs(i int) int {
	if i < 0 {
		i *= -1
	}
	return i
}

func IMin(numbers []int) (int, int) {
	if len(numbers) == 0 {
		return -1, 0
	}
	min := numbers[0]
	var min_idx int
	for idx := 1; idx < len(numbers); idx++ {
		if numbers[idx] < min {
			min_idx = idx
			min = numbers[idx]
		}
	}
	return min_idx, min
}

func IMax(numbers []int) (int, int) {
	if len(numbers) == 0 {
		return -1, 0
	}
	max := numbers[0]
	var max_idx int
	for idx := 1; idx < len(numbers); idx++ {
		if numbers[idx] > max {
			max_idx = idx
			max = numbers[idx]
		}
	}
	return max_idx, max
}

func StringIsSubset(subset string, superset string) bool {
	return 0 > strings.IndexFunc(subset, func(r rune) bool {
		return strings.IndexRune(superset, r) < 0
	})
}

// Rotate an index around in a ring by a positive or negative value.
//
// The result is always positive and less than length.
func Rotate(index int, by int, length int) int {
	if by < 0 {
		index += (by % length)
		if index < 0 {
			index += length
		}
	} else {
		index = (index + by) % length
	}
	return index
}

// Like a ring.Ring, but uses an array buffer rather than a linked list.
type RingBuffer struct {
	_data []interface{}

	// Points to the tail, i.e. the last (most recently inserted) element.
	//
	// The head, i.e. first (least recently inserted) is at (_tail + 1) % Len().
	_tail int
}

func NewRingBuffer(size int) *RingBuffer {
	r := new(RingBuffer)
	r._data = make([]interface{}, size)
	r._tail = size - 1
	return r
}

func (r *RingBuffer) rotated(n int) int {
	return Rotate(r._tail, n, r.Len())
}

func (r *RingBuffer) head() int {
	return r.rotated(1)
}

func (r *RingBuffer) tail() int {
	return r._tail
}

// Push a new value.
//
// Returns the least-recently inserted value which was overwritten.
func (r *RingBuffer) Push(value interface{}) interface{} {
	r.Move(1)
	tail := r.tail()
	lru := r._data[tail]
	r._data[tail] = value
	return lru
}

// Pop the most-recently inserted value.
func (r *RingBuffer) Pop() interface{} {
	result := r._data[r.tail()]
	r.Move(-1)
	return result
}

// Return the most recently inserted (last) element.
func (r *RingBuffer) Last() interface{} {
	return r._data[r.tail()]
}

// Return the least recently inserted (first) element.
func (r *RingBuffer) First() interface{} {
	return r._data[r.head()]
}

func (r *RingBuffer) Len() int {
	return len(r._data)
}

// Move n % r.Len() elements backward (n < 0) or forward (n >= 0).
func (r *RingBuffer) Move(n int) *RingBuffer {
	r._tail = r.rotated(n)
	return r
}

// Get the n-th least recently used element.
//
// For example, Get(0) === First(),
// and Get(-1) === Get(Len()-1) === Last().
func (r *RingBuffer) Get(n int) interface{} {
	return r._data[Rotate(r.head(), n, r.Len())]
}

// Get the n-th most recently used element.
//
// For example, GetLast(0) == Last(),
// and Getlast(-1) === GetLast(Len()-1) === First().
func (r *RingBuffer) GetLast(n int) interface{} {
	return r._data[Rotate(r.tail(), -n, r.Len())]
}

// Call f on each element of the ring in forward order.
func (r *RingBuffer) Do(f func(interface{})) {
	// Do the head, then wrap around until we reach head again.
	head := r.head()
	f(r._data[head])
	for idx := r.rotated(2); idx != head; idx++ {
		f(r._data[idx])
	}
}

// Call f on each element of the ring in forward order so long as f returns true
func (r *RingBuffer) DoWhile(f func(interface{}) bool) {
	// Do the head, then wrap around until we reach head again.
	head := r.head()
	if f(r._data[head]) {
		for idx := r.rotated(2); idx != head && f(r._data[idx]); idx++ {
		}
	}
}
