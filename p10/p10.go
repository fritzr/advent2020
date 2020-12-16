package p10

import (
  "fmt"
  "sort"
  "errors"
  "github.com/fritzr/advent2020/util"
)

var gVerbose = false

func combFrom(idx int, adapters []int, memo map[int]int) (nComb int) {
  // Special case: combFrom(-1) uses the implicit starting adapter (0).
  var adapter int
  if idx >= 0 {
    // Base case 1: memoized lookup
    if memo[idx] != 0 {
      return memo[idx]
    }
    adapter = adapters[idx]
  }

  // Number of combinations is the number of ways we can arrange our child
  // nodes.
  nextIndex := idx + 1
  for nextIndex < len(adapters) && adapters[nextIndex] - adapter <= 3 {
    nComb += combFrom(nextIndex, adapters, memo)
    nextIndex++
  }

  // Base case 2: leaf node.
  if nComb == 0 {
    nComb = 1
  }
  memo[idx] = nComb
  return nComb
}

func AdapterCombinations(adapters []int) int {
  // Memoized table of combinations possible starting from a given adapter.
  memo := make(map[int]int, len(adapters))
  return combFrom(-1, adapters, memo)
}

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

  fmt.Printf("# of ways to arrange adapters: %d\n",
    AdapterCombinations(adapters))

  return nil
}
