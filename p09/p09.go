package p09

import (
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

// Find a contiguous series of numbers with the desired sum.
func FindSumWindow(data []int, sum int) (lo int, hi int, windowSum int) {
  if len(data) == 0 {
    return 0, 0, 0
  }

  // The numbers are roughly increasing. Use a sliding window.
  // When the sum is too large, remove the lowest values.
  // When the sum is too small, add a higher value.
  // We may run into problems because the numbers do not increase monotonically.
  next := hi
  sgn := 1
  fmt.Printf("  pushing [%d] %d\n", hi, data[hi])
  for hi < len(data) && windowSum != sum {
    // Adjust the sum according to the last operation.
    windowSum += sgn * data[next]
    if gVerbose {
      fmt.Printf("Sum [%d:%d](%d,...,%d) = %d\n",
        lo, hi, data[lo], data[hi], windowSum)
    }
    if windowSum < sum {
      // Expand the window to increase the sum.
      sgn = 1

      /* XXX do we need to consider expanding the left side of the window?
      if lo > 0 && (hi >= len(data) - 1 || data[lo] < data[hi + 1]) {
        lo--
        next = lo
      } else { ... }
      */

      // Move up the right bound and add the newest value.
      hi++
      next = hi
      sgn = 1
      if gVerbose {
        if hi < len(data) {
          fmt.Printf("  pushing [%d] %d\n", hi, data[hi])
        } else {
          fmt.Printf("  at end [%d]\n", hi)
        }
      }
    } else if windowSum > sum {
      sgn = -1
      if lo < hi {
        // Subtract the oldest value and move up the left bound.
        if gVerbose {
          fmt.Printf("  popping [%d] %d\n", lo, data[lo])
        }
        next = lo
        lo++
      } else {
        // Bum window... lo == hi and the value here is too large.
        // Try to move past it.
        if gVerbose {
          fmt.Printf("  skipping [%d] %d\n", lo, data[lo])
        }
        next = lo
        lo++
        hi++
      }
    }
  }
  // We are supposed to return inclusive bounds, so don't return the size.
  if hi >= len(data) {
    hi = len(data) - 1
  }
  if lo >= len(data) {
    lo = len(data) - 1
  }
  return lo, hi, windowSum
}

func Usage() {
  fmt.Println("usage: advent2020 9 [main opts...] [-w window_size=25]\n")
  fmt.Println("")
  fmt.Println("The -w option allows you to customize the XMAS window size.\n")
}

func ParseArgs(args []string) (windowSize int, err error) {
  if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
    Usage()
    return windowSize, nil
  }

  if len(args) > 0 && (args[0] == "-w" || args[0] == "--window-size") {
    if len(args) == 1 {
      Usage()
      return windowSize, errors.New("-w expects an argument")
    }
    windowSize, err = strconv.Atoi(args[1])
    if err != nil {
      return windowSize, err
    }
    if windowSize <= 0 {
      return windowSize, errors.New("window size must be positive")
    }
  } else {
    windowSize = 25
  }

  return windowSize, nil
}

func VerifySum(data []int, sum int) error {
  if len(data) == 0 {
    return errors.New("sum window is empty")
  }
  value := data[0]
  for idx := 1; idx < len(data); idx++ {
    value += data[idx]
  }
  if value != sum {
    return errors.New(fmt.Sprintf("Sum of window is %d, expected %d",
      value, sum))
  }
  return nil
}

func Main(input_path string, verbose bool, args []string) error {
  gVerbose = verbose

  windowSize, err := ParseArgs(args)
  if err != nil {
    return err
  }

  var data []int
  data, err = util.ReadNumbersFromFile(input_path)
  if err != nil {
    return err
  }

  // Validate the input, first of all.
  validator := NewXMASValidator(windowSize)
  var (idx int; value int)
  for idx, value = range data {
    err = validator.Read(value)
    if err != nil {
      break
    }
  }

  // Part 1: expecting a bad value.
  if err == BadXMASValue {
    fmt.Printf("First bad value was [%d] %d\n", idx, value)
  } else if err == nil {
    return errors.New("all numbers were unexpectedly valid!")
  } else {
    return err
  }

  // Part 2: find a contiguous sequence of numbers which sum to the bad value.
  loIndex, hiIndex, sum := FindSumWindow(data[:idx], value)
  if sum == value {
    lo := data[loIndex]
    hi := data[hiIndex]
    window := data[loIndex:(hiIndex+1)]
    fmt.Printf("Sum window: [%d:%d] (%d,...,%d)\n", loIndex, hiIndex, lo, hi)
    if err = VerifySum(window, sum); err != nil {
      return err
    }
    _, min := util.IMin(window)
    _, max := util.IMax(window)
    fmt.Printf("  min=%d, max=%d, min + max=%d\n", min, max, min + max)
  } else {
    return errors.New(fmt.Sprintf(
      "failed to find sum window! stopped at [%d:%d] (%d,...,%d)\n",
        loIndex, hiIndex, data[loIndex], data[hiIndex]))
  }

  return nil
}
