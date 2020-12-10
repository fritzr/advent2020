
package main

import (
  "log"
  "os"
  "fmt"
  "path"
  "strconv"
  // "github.com/fritzr/advent2020/common"
  "github.com/fritzr/advent2020/p01"
)

type AdventMain func(path string, verbose bool) (error)

func main() {
  args := os.Args[1:]

  if len(args) > 0 && args[0] == "-h" {
    fmt.Printf("usage: go run %s [-v] [day]\n", path.Base(os.Args[0]))
    os.Exit(1)
  }

  verbose := false
  if len(args) > 0 && args[0] == "-v" {
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
    }
  }

  RunPuzzle(day, puzzle, verbose)
}

func RunPuzzle(day int, puzzle AdventMain, verbose bool) {
  path := path.Join(".", fmt.Sprintf("p%02d", day), "input")
  err := puzzle(path, verbose)
  if err != nil {
    log.Fatal(err)
  }
}
