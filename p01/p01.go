package p01

import (
  "strings"
  "errors"
  "log"
  "fmt"
  "strconv"
  "github.com/fritzr/advent2020/util"
)

func NumbersSummingToIter(sum int, Next func() bool, Get func() int) (int, int, error) {
  table := make(map[int]int)

  // Once we see one number, we know the other number which sums with it.
  // Mark both once; when a number is marked twice, it and its pair sum to 2020.
  for Next() {
    value := Get()
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

// Part 1
func NumbersSummingTo(input []int, sum int) (int, int, error) {
  idx := -1
  return NumbersSummingToIter(sum,
    func() bool { idx++; return idx < len(input) },
    func() int { return input[idx] })
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

func Usage() {
  fmt.Println("usage: go run advent2020 1 [N [SUM=2020]]")
  fmt.Println()
  fmt.Println("Find N numbers which sum to SUM in the input.")
  fmt.Println("If no args are given, print the results required by the puzzle.")
  fmt.Println("If N is given, the default SUM is 2020 (from the puzzle).")
}

func do_sum(input []int, N int, sum int) error {
  var err error
  var result []int
  result, err = NNumbersSummingTo(N, input, sum)
  if err != nil {
    return err
  }
  if len(result) != N {
    return errors.New(fmt.Sprintf("expected %d numbers, reported %d",
    N, len(result)))
  }

  // String representation of the numbers (1, 2, 3, ...)
  var nrep strings.Builder
  for idx, value := range result {
    nrep.WriteString(strconv.Itoa(value))
    if idx != len(result) - 1 {
      nrep.WriteString(", ")
    }
  }

  fmt.Printf("%d numbers which sum to %d: %s\n  Product: %d\n",
    N, sum, nrep.String(), util.Product(result))

  return nil
}


func Day1(input []int) error {
  var err error

  // Part 1
  if err = do_sum(input, 2, 2020); err != nil {
    return err
  }

  // Part 2
  if err = do_sum(input, 3, 2020); err != nil {
    return err
  }

  return nil
}

func Main(input_path string, verbose bool, args []string) (error) {
  // Read the input.
  var input []int
  var err error
  input, err = util.ReadNumbersFromFile(input_path)
  if verbose {
    LogNumbers(input)
  }
  if err != nil {
    return err
  }

  // Check args. On empty, just do what the puzzle asked for.
  if len(args) == 0 {
    return Day1(input)
  }

  // Otherwise, grab N and then look for the SUM.

  if args[0] == "-h" || args[0] == "--help" {
    Usage()
    return nil
  }

  // N
  var N int
  N, err = strconv.Atoi(args[0])
  if err != nil {
    return err
  }
  args = args[1:]

  // SUM
  sum := 2020
  if len(args) > 0 {
    sum, err = strconv.Atoi(args[0])
    if err != nil {
      return err
    }
    args = args[1:]
  }

  return do_sum(input, N, sum)
}
