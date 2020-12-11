package p01

import (
  "io"
  "strings"
  "errors"
  "log"
  "os"
  "bufio"
  "fmt"
  "strconv"
  "github.com/fritzr/advent2020/util"
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
  input, err = ReadNumbersFromFile(input_path, verbose)
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
