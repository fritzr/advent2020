package p05

import (
  "io"
  "bufio"
  "fmt"
  "log"
  "os"
  "errors"
  "github.com/fritzr/advent2020/p04"
)

type Seat struct {
  row int
  column int
}

func (s *Seat) ID() int {
  return s.row * 8 + s.column
}

type BoardingPass struct {
  steps []byte
}

func Bisect(steps []byte, lchar byte, uchar byte, min int, max int) int {
  for _, step := range steps[:len(steps)-1] {
    span := (max + 1 - min) / 2
    switch(step) {
    case lchar: max -= span
    case uchar: min += span
    default: panic("Bisect: invalid steps")
    }
  }
  switch(steps[len(steps)-1]) {
  case lchar: return min
  case uchar: return max
  default: panic("Bisect: invalid steps")
  }
}

func (p *BoardingPass) Decode(nrows int, ncols int) Seat {
  row := Bisect(p.steps[:7], 'F', 'B', 0, nrows - 1)
  column := Bisect(p.steps[7:], 'L', 'R', 0, ncols - 1)
  return Seat{row, column}
}

func NewBoardingPass(line string) (BoardingPass, error) {
  if len(line) != 10 {
    return BoardingPass{},
      errors.New(fmt.Sprintf("BoardingPass: invalid length for '%s'", line))
  }
  if !p04.StringIsSubset(line[:7], "FB") || !p04.StringIsSubset(line[7:], "LR") {
    return BoardingPass{},
      errors.New(fmt.Sprintf("BoardingPass: invalid characters in '%s'", line))
  }
  return BoardingPass{[]byte(line)}, nil
}

func ReadBoardingPasses(r io.Reader) ([]BoardingPass, error) {
  scanner := bufio.NewScanner(r)
  scanner.Split(bufio.ScanLines)
  passes := make([]BoardingPass, 0, 867)
  for scanner.Scan() {
    pass, err := NewBoardingPass(scanner.Text())
    if err != nil {
      return passes, err
    }
    passes = append(passes, pass)
  }
  return passes, scanner.Err()
}

func ReadBoardingPassesFromFile(path string) ([]BoardingPass, error) {
  file, err := os.Open(path)
  if err != nil {
    return []BoardingPass{}, err
  }
  return ReadBoardingPasses(file)
}

func find_missing_seat(passes []BoardingPass, max_id int) int {
  // Find the missing seat... Start by sorting all seats by ID.
  seats := make([]Seat, max_id + 1)
  for _, pass := range passes {
    seat := pass.Decode(128, 8)
    seats[seat.ID()] = seat
  }

  // Look for a missing ID.
  for id := 1; id < len(seats) - 1; id++ {
    if seats[id].ID() == 0 && seats[id-1].ID() != 0 && seats[id+1].ID() != 0 {
      return id
    }
  }

  return -1
}

func Main(input_path string, verbose bool, args []string) error {
  passes, err := ReadBoardingPassesFromFile(input_path)
  if err != nil {
    return err
  }

  max_id := 0
  var max_seat *Seat
  var max_pass *BoardingPass
  for _, pass := range passes {
    seat := pass.Decode(128, 8)
    id := seat.ID()
    if verbose {
      log.Print(fmt.Sprintf("%s => (%d, %d) [ID=%d]\n",
        pass.steps, seat.row, seat.column, id))
    }
    if id > max_id {
      max_id = id
      max_seat = &seat
      max_pass = &pass
    }
  }

  fmt.Printf("Read %d boarding passes.\n", len(passes))

  // Part 1
  fmt.Printf("Highest seat from: %s => (%d, %d) [ID=%d]\n",
    max_pass.steps, max_seat.row, max_seat.column, max_id)

  // Part 2
  missing_id := find_missing_seat(passes, max_id)
  fmt.Printf("Missing seat ID: %d\n", missing_id)

  return nil
}
