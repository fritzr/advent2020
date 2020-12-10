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


func NumbersSummingTo(input []int, sum int) (int, int, error) {
  table := make([]int, sum, sum)

  for _, value := range input {
    other_value := sum - value
    if table[other_value] > 0 {
      return other_value, value, nil
    } else {
      table[value] += 1
      table[other_value] += 1
    }
  }
  return 0, 0, errors.New(fmt.Sprintf("no values sum to %d", sum))
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
  input, err := ReadNumbersFromFile(input_path, verbose)
  if verbose {
    LogNumbers(input)
  }
  if err != nil {
    return err
  }
  n1, n2, err := NumbersSummingTo(input, 2020)
  if err != nil {
    return err
  }
  fmt.Printf("Numbers which sum to 2020: %d, %d\n  %d * %d = %d\n",
    n1, n2, n1, n2, n1 * n2)
  return nil
}
