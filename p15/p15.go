package p15

import (
  "io/ioutil"
  "strings"
  "strconv"
  "errors"
  "fmt"
)

func stringsToInts(strings []string) (ints []int, err error) {
  ints = make([]int, len(strings))
  for index, str := range strings {
    ints[index], err = strconv.Atoi(str)
    if err != nil {
      break
    }
  }
  return ints, err
}

func Usage() {
  fmt.Println("usage: advent2020 [main opts...] [-n turns]")
  fmt.Println()
  fmt.Println("Run the game for the specified number of turns (default 2020).")
}

func Main(input_path string, verbose bool, args []string) error {
  data, err := ioutil.ReadFile(input_path)
  if err != nil {
    return err
  }

  numbers, err := stringsToInts(strings.Split(strings.Trim(string(data), "\n"), ","))
  if err != nil {
    return err
  }

  nTurns := 2020

  if len(args) > 0 {
    if args[0] == "-h" || args[0] == "--help" {
      Usage()
      return nil
    }
    if args[0] == "-n" {
      if len(args) < 2 {
        return errors.New("option -n requires an argument")
      }
      nTurns, err = strconv.Atoi(args[1])
      if err != nil {
        return err
      }
    }
  }

  // Initialize the age map.
  lastSpoken := make(map[int]int)
  for turn, startNumber := range numbers {
    lastSpoken[startNumber] = turn + 1 // never assign turn 0
    if verbose {
      fmt.Printf("  Turn %4d: %d\n", turn + 1, startNumber)
    }
  }

  last := numbers[len(numbers) - 1]
  next := 0 // assuming the input numbers are all unique
  for turn := len(lastSpoken) + 1; turn <= nTurns; turn++ {
    if verbose {
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

  fmt.Printf("The %d-th number spoken is %d.\n", nTurns, last)

  return nil
}
