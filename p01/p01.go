package p01

import (
  "io"
  "errors"
  "log"
  "os"
  "bufio"
  "fmt"
  "strconv"
)


// Part 1
func NumbersSummingTo(input []int, sum int) (int, int, error) {
  table := make([]int, sum, sum)

  // Once we see one number, we know the other number which sums with it.
  // Mark both once; when a number is marked twice, it and its pair sum to 2020.
  for _, value := range input {
    if value < sum {
      other_value := sum - value
      if table[other_value] > 0 {
        return other_value, value, nil
      } else {
        table[value] += 1
        table[other_value] += 1
      }
    }
  }
  return 0, 0, errors.New(fmt.Sprintf("no values sum to %d", sum))
}

// Part 2: N numbers (depth) from I which sum to S
func NNumbersSummingTo(depth int, input []int, sum int) ([]int, error) {
  if (depth < 2) {
    return []int{}, errors.New("don't bother summing less than 2 numbers")
  }

  if (depth == 2) {
    n1, n2, err := NumbersSummingTo(input, sum)
    if err != nil {
      return []int{}, err
    }
    return []int{n1, n2}, nil
  }

  for _, value := range input {
    if value < sum {
      others, err := NNumbersSummingTo(depth - 1, input, sum - value)
      if err == nil {
        return append(others, value), nil
      }
    }
  }

  return []int{}, errors.New(fmt.Sprintf("no %d numbers sum to %d", depth, sum))
}

func ReadNumbers(r io.Reader, verbose bool) ([]int, error) {
  result := make([]int, 0, 200)

  bufreader := bufio.NewReader(r)
  line_bytes, shortRead, err := bufreader.ReadLine()
  var i int
  for err == nil && !shortRead {
    line := string(line_bytes)
    if verbose {
      log.Printf("p01: read line: '%s'\n", line)
    }
    i, err = strconv.Atoi(line)
    if err == nil {
      result = append(result, i)
    } else {
      return result, err
    }
    line_bytes, shortRead, err = bufreader.ReadLine()
  }
  if shortRead {
    err = errors.New("short read!")
  }
  if err == io.EOF {
    err = nil
  }
  return result, err
}

func ReadNumbersFromFile(path string, verbose bool) ([]int, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()
  return ReadNumbers(file, verbose)
}

func LogNumbers(input []int) {
  log.Print("p01: got ", len(input), " values")
  i := 0
  for ; i < (len(input) - 8); i++ {
    log.Printf("%5d %5d %5d %5d %5d %5d %5d %5d\n",
      input[i], input[i+1], input[i+2], input[i+3],
      input[i+4], input[i+5], input[i+6], input[i+7])
  }
  for ; i < len(input); i++ {
    log.Printf("%5d ", input[i])
  }
  if len(input) > 0 {
    log.Printf("\n")
  }
}

func Main(input_path string, verbose bool) (error) {
  var input []int
  var err error
  input, err = ReadNumbersFromFile(input_path, verbose)
  if verbose {
    LogNumbers(input)
  }
  if err != nil {
    return err
  }

  // Part 1
  var n1 int
  var n2 int
  n1, n2, err = NumbersSummingTo(input, 2020)
  if err != nil {
    return err
  }
  fmt.Printf("Numbers which sum to 2020: %d, %d\n%d * %d = %d\n",
    n1, n2, n1, n2, n1 * n2)

  // Part 2
  var nums []int
  nums, err = NNumbersSummingTo(3, input, 2020)
  if err != nil {
    return err
  }
  if len(nums) != 3 {
    return errors.New(fmt.Sprintf("got %d numbers, expected 3", len(nums)))
  }
  fmt.Printf("Three numbers which sum to 2020: %d, %d, %d\n%d * %d * %d = %d\n",
    nums[0], nums[1], nums[2], nums[0], nums[1], nums[2],
    nums[0] * nums[1] * nums[2])

  return nil
}
