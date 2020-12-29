package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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
	"github.com/fritzr/advent2020/p14"
	"github.com/fritzr/advent2020/p15"
	"github.com/fritzr/advent2020/p16"
	"github.com/fritzr/advent2020/p17"
	"github.com/fritzr/advent2020/p18"
	"github.com/fritzr/advent2020/p19"
	"github.com/fritzr/advent2020/p20"
)

type AdventMain func(path string, verbose bool, args []string) error

var puzzles = [...]AdventMain{
	p01.Main,
	p02.Main,
	p03.Main,
	p04.Main,
	p05.Main,
	p06.Main,
	p07.Main,
	p08.Main,
	p09.Main,
	p10.Main,
	p11.Main,
	p12.Main,
	p13.Main,
	p14.Main,
	p15.Main,
	p16.Main,
	p17.Main,
	p18.Main,
	p19.Main,
	p20.Main,
}

var verbose bool
var input string

const (
	verboseUsage = "enable debug messages"
	inputUsage   = "puzzle input path"
)

func init() {
	flag.BoolVar(&verbose, "verbose", false, verboseUsage)
	flag.BoolVar(&verbose, "v", false, verboseUsage)
	flag.StringVar(&input, "input", "", inputUsage)
	flag.StringVar(&input, "i", "", inputUsage)
}

func Usage() {
	fmt.Fprintf(flag.CommandLine.Output(),
		`usage: %s [OPTIONS...] [--] [[day] [ARGS...]]

Run the given day's puzzle (defaults to the latest implemented puzzle).
Additional puzzle-specific arguments may be accepted for some puzzles.
Add -h or --help after the day to find out.
All puzzles accept '-v' to run verbose and '-i PATH' to override the input.

OPTIONS are:
`, path.Base(os.Args[0]))
	flag.PrintDefaults()
}

func main() {
	flag.Usage = Usage
	flag.Parse()

	// day
	day := len(puzzles)
	puzzle := puzzles[len(puzzles)-1]
	args := flag.Args()
	var err error
	if flag.NArg() > 0 {
		day, err = strconv.Atoi(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}

		if day <= len(puzzles) {
			puzzle = puzzles[day-1]
		} else {
			log.Fatal(fmt.Sprintf("unimplemented day '%d'\n", day))
		}
		args = args[1:]
	}

	// input override (-i path)
	if input == "" {
		input = path.Join(".", fmt.Sprintf("p%02d", day), "input")
	}

	// Run the selected puzzle. Pass additional arguments.
	err = puzzle(input, verbose, args)
	if err != nil {
		log.Fatal(err)
	}
}
