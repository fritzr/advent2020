
package main

import (
  "log"
  "os"
  "fmt"
  "path"
  "strconv"
  // "github.com/fritzr/advent2020/common"
  "github.com/fritzr/advent2020/p01"
  "github.com/fritzr/advent2020/p02"
  "github.com/fritzr/advent2020/p03"
  "github.com/fritzr/advent2020/p04"
  "github.com/fritzr/advent2020/p05"
  "github.com/fritzr/advent2020/p06"
  "github.com/fritzr/advent2020/p07"
)

type AdventMain func(path string, verbose bool, args []string) (error)

func main() {
  args := os.Args[1:]

  if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
    fmt.Printf("usage: go run %s [-v] [day [args...]]\n", path.Base(os.Args[0]))
    fmt.Println("")
    fmt.Println(
"Run the puzzle the given day's puzzle. Additional puzzle-specific arguments")
    fmt.Println(
"may be accepted for some puzzles. Add -h or --help after the day to find out.")
    os.Exit(1)
  }

  verbose := false
  if len(args) > 0 && (args[0] == "-v" || args[0] == "--verbose") {
    verbose = true
    args = args[1:]
  }

  day := 1
  var puzzle AdventMain = p01.Main
  var err error
  if len(args) > 0 {
    day, err = strconv.Atoi(args[0])
    if err != nil {
      log.Fatal(err)
    }

    switch day {
    default: log.Fatal(fmt.Printf("unimplemented day '%d'\n", day))
    case 1: puzzle = p01.Main
    case 2: puzzle = p02.Main
    case 3: puzzle = p03.Main
    case 4: puzzle = p04.Main
    case 5: puzzle = p05.Main
    case 6: puzzle = p06.Main
    case 7: puzzle = p07.Main
    }

    args = args[1:]
  }

  RunPuzzle(day, puzzle, verbose, args)
}

func RunPuzzle(day int, puzzle AdventMain, verbose bool, args []string) {
  path := path.Join(".", fmt.Sprintf("p%02d", day), "input")
  err := puzzle(path, verbose, args)
  if err != nil {
    log.Fatal(err)
  }
}
