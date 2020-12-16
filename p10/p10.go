package p10

import (
  "fmt"
  "sort"
  "errors"
  "github.com/fritzr/advent2020/util"
)

var gVerbose = false

func Main(input_path string, verbose bool, args []string) error {
  adapters, err := util.ReadNumbersFromFile(input_path)
  if err != nil {
    return err
  }

  sort.Ints(adapters)
  adapters = append(adapters, 3 + adapters[len(adapters)-1])
  diffs := make([]int, 4)

  last := 0
  for idx, joltage := range adapters {
    diff := joltage - last
    if diff <= 3 {
      diffs[diff]++
    } else {
      return errors.New(fmt.Sprintf(
        "unusable adapter from [%d] %d => [%d] %d",
        idx - 1, last, idx, joltage))
    }
    last = joltage
  }

  fmt.Printf("Diffs (%d):\n", len(adapters))
  util.PrintArray(diffs)
  fmt.Printf("  [1] %d * [3] %d = %d\n",
    diffs[1], diffs[3], diffs[1] * diffs[3])

  return nil
}
