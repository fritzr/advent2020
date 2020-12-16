package p09

import (
  "fmt"
  "errors"
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

func Main(input_path string, verbose bool, args []string) error {
  gVerbose = verbose
  data, err := util.ReadNumbersFromFile(input_path)
  if err != nil {
    return err
  }

  // Validate the input, first of all.
  validator := NewXMASValidator(25)
  var value int
  for _, value = range data {
    err = validator.Read(value)
    if err != nil {
      break
    }
  }

  // Part 1: expecting a bad value.
  if err == BadXMASValue {
    fmt.Printf("First bad value was %d\n", value)
  } else if err == nil {
    return errors.New("all numbers were unexpectedly valid!")
  } else {
    return err
  }

  // Part 2: find a contiguous sequence of numbers which sum to the bad value.

  return nil
}
