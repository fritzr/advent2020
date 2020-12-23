# advent2020
Advent of Code 2020 (https://adventofcode.com/2020)

## Running

Run a puzzle like this (from the root directory):

```sh
$ go run . --help
usage: go run advent2020 [-v] [day [args...]]

Run the puzzle the given day's puzzle. Additional puzzle-specific arguments
may be accepted for some puzzles. Add -h or --help after the day to find out.
```

The puzzle day defaults to 1.
With -v or --verbose, tell the puzzle to log debugging information.
Some puzzles accept additional arguments:

```sh
$ go run . 1 --help
```

## Dependencies

Everything except day 20 uses only core Go libraries. Day 20 uses `gonum/mat`
as an excuse for me to learn how Go deals with heavy numeric calculations.
