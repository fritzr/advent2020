package p15

import (
  "io/ioutil"
  "strings"
  "strconv"
  "errors"
  "fmt"
  "github.com/fritzr/advent2020/util"
)

var gVerbose = false

func Usage() {
  fmt.Println("usage: advent2020 [main opts...] [-n turns]")
  fmt.Println()
  fmt.Println("Run the game for the specified number of turns.")
  fmt.Println("If not given, run once with 2020 turns and once with 30000000.")
}

func RambunctiousRecitation(init []int, rounds int) int {
  // Initialize the age map.
  lastSpoken := make(map[int]int)
  for turn, startNumber := range init {
    lastSpoken[startNumber] = turn + 1 // never assign turn 0
    if gVerbose {
      fmt.Printf("  Turn %4d: %d\n", turn + 1, startNumber)
    }
  }

  last := init[len(init) - 1]
  next := 0 // assuming the input numbers are all unique
  for turn := len(lastSpoken) + 1; turn <= rounds; turn++ {
    if gVerbose {
      fmt.Printf("  Turn %4d: %d\n", turn, next)
    }
    last = next
    if lastSpoken[last] == 0 {
      // If last turn was the first time the number was spoken, 0 is next.
      next = 0
    } else {
      // How many turns ago was the last time the number was spoken?
      next = turn - lastSpoken[last]
    }
    lastSpoken[last] = turn
  }
  return last
}

func Main(input_path string, verbose bool, args []string) error {
  gVerbose = verbose
  data, err := ioutil.ReadFile(input_path)
  if err != nil {
    return err
  }

  fields := strings.Split(strings.Trim(string(data), "\n"), ",")
  numbers, err := util.FieldsToInts(fields)
  if err != nil {
    return err
  }

  if len(args) > 0 {
    if args[0] == "-h" || args[0] == "--help" {
      Usage()
      return nil
    }
    if args[0] == "-n" {
      if len(args) < 2 {
        return errors.New("option -n requires an argument")
      }
      nTurns, err := strconv.Atoi(args[1])
      if err != nil {
        return err
      }

      // Just do the requested amount.
      fmt.Printf("The %d-th number spoken is %d.\n", nTurns,
        RambunctiousRecitation(numbers, nTurns))
      return nil
    }
  }

  // Part 1: 2020 turns
  nTurns := 2020
  fmt.Printf("The %d-th number spoken is %d.\n", nTurns,
    RambunctiousRecitation(numbers, nTurns))

  // Part 2:
  nTurns = 30000000
  fmt.Printf("The %d-th number spoken is %d.\n", nTurns,
    RambunctiousRecitation(numbers, nTurns))

  return nil
}
