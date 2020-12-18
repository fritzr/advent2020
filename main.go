
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
  "github.com/fritzr/advent2020/p08"
  "github.com/fritzr/advent2020/p09"
  "github.com/fritzr/advent2020/p10"
  "github.com/fritzr/advent2020/p11"
  "github.com/fritzr/advent2020/p12"
  "github.com/fritzr/advent2020/p13"
  "github.com/fritzr/advent2020/p17"
  "github.com/fritzr/advent2020/p18"
)

type AdventMain func(path string, verbose bool, args []string) (error)

func main() {
  args := os.Args[1:]

  if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
    fmt.Printf("usage: go run %s [day [-v] [-i path] [args...]]\n",
    path.Base(os.Args[0]))
    fmt.Println("")
    fmt.Println(
"Run the puzzle the given day's puzzle. Additional puzzle-specific arguments")
    fmt.Println(
"may be accepted for some puzzles. Add -h or --help after the day to find out.")
    fmt.Println(
"All puzzles accept '-v' to run verbose and '-i PATH' to override the input.")
    os.Exit(1)
  }

  // day
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
    case 8: puzzle = p08.Main
    case 9: puzzle = p09.Main
    case 10: puzzle = p10.Main
    case 11: puzzle = p11.Main
    case 12: puzzle = p12.Main
    case 13: puzzle = p13.Main
    case 17: puzzle = p17.Main
    case 18: puzzle = p18.Main
    }

    args = args[1:]
  }

  // -v
  verbose := false
  if len(args) > 0 && (args[0] == "-v" || args[0] == "--verbose") {
    verbose = true
    args = args[1:]
  }

  // input override (-i path)
  input := path.Join(".", fmt.Sprintf("p%02d", day), "input")
  if len(args) > 0 && (args[0] == "-i" || args[0] == "--input") {
    if len(args) < 2 {
      log.Fatal("missing argument to -i")
    }
    input = args[1]
    args = args[2:]
  }

  // Run the selected puzzle. Pass additional arguments.
  err = puzzle(input, verbose, args)
  if err != nil {
    log.Fatal(err)
  }
}
