package p09

import (
  "io"
  "os"
  "bufio"
  "fmt"
  "errors"
  "strconv"
  "github.com/fritzr/advent2020/util"
  "github.com/fritzr/advent2020/p01"
)

var gVerbose = false

var BadXMASValue = errors.New("bad XMAS value")

type XMASValidator struct {
  size int
  count int
  buffer *util.RingBuffer
}

func NewXMASValidator(size int) *XMASValidator {
  v := new(XMASValidator)
  v.buffer = util.NewRingBuffer(size)
  v.size = size
  return v
}

func (v *XMASValidator) Read(value int) error {
  if v.count >= v.size {
    // Verify the value is a sum of two previous values.
    mruIndex := -1
    size := v.size
    // Called like: for Next() { x := Get(); ... }
    v1, v2, err := p01.NumbersSummingToIter(value,
      func() bool { // Next()
        mruIndex++
        return mruIndex < size
      },
      func() int { // Get()
        return v.buffer.GetLast(mruIndex).(int)
      })
    if err != nil {
      return BadXMASValue
    }
    if gVerbose {
      fmt.Printf("%d: sum of %d and %d\n", value, v1, v2)
    }
  }
  if gVerbose {
    fmt.Printf("inserted %d\n", value)
  }
  v.buffer.Push(value)
  v.count += 1 // number of elements added
  return nil
}

func ReadXMAS(input io.Reader, ringSize int) (int, error) {
  scanner := bufio.NewScanner(input)
  scanner.Split(bufio.ScanLines)
  validator := NewXMASValidator(ringSize)
  for scanner.Scan() {
    value, err := strconv.Atoi(scanner.Text())
    if err != nil {
      return value, err
    }
    if err = validator.Read(value); err != nil {
      return value, err
    }
  }
  return -1, scanner.Err()
}

func ReadXMASFromFile(path string, ringSize int) (int, error) {
  file, err := os.Open(path)
  if err != nil {
    return 0, err
  }
  defer file.Close()
  return ReadXMAS(file, ringSize)
}

func Main(input_path string, verbose bool, args []string) error {
  gVerbose = verbose
  firstInvalid, err := ReadXMASFromFile(input_path, 25)
  if err != nil && err != BadXMASValue {
    return err
  }

  // Part1: expecting a bad value
  if err == BadXMASValue {
    fmt.Printf("First bad value was %d\n", firstInvalid)
  }

  return nil
}
