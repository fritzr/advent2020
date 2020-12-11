package p03

import (
  "fmt"
  "io"
  "os"
  "github.com/fritzr/advent2020/p01"
)

func get_line_width(r *os.File) (int, error) {
  nread := 1
  pos := 0
  var err error
  char := make([]byte, 1)
  for char[0] != '\n' && nread > 0 && err == nil {
    nread, err = r.Read(char)
    pos += nread
  }
  return pos - 1, err
}

func TobogganSled(rows *os.File, slope int, speed int) (int, int, error) {
  const tree_char = '#'
  const open_char = '.'


  // Skip the first line, since we always start in the top-left (an open space)
  line_width, err := get_line_width(rows)
  if err != nil {
    return 0, 0, err
  }

  ntrees := 0
  nopen := 0
  line_pos := 0
  nread := 1
  char := make([]byte, 1)

  _, err = rows.Seek(1, io.SeekStart)
  if err != nil {
    return 0, 0, err
  }
  for nread > 0 && err == nil {
    // need to seek to end of line, then to the next line position
    line_remaining := line_width - line_pos
    line_pos = (line_pos + slope) % line_width

    // right N, down 1 (line width)
    _, err = rows.Seek(
      int64(line_remaining + line_pos + (speed - 1) * (line_width + 1)),
      io.SeekCurrent)
    if err != nil {
      return 0, 0, err
    }

    nread, err = rows.Read(char)
    if nread == len(char) && err == nil {
      if char[0] == tree_char {
        ntrees++
      } else {
        nopen++
      }
    }
  }

  if err == io.EOF {
    err = nil
  }

  return ntrees, nopen, err
}

func do_sled(file *os.File, slope int, speed int) (int, int, error) {
  var (ntrees int; nopen int; err error)
  _, err = file.Seek(0, io.SeekStart)
  if err != nil {
    return 0, 0, err
  }

  ntrees, nopen, err = TobogganSled(file, slope, speed)
  if err != nil {
    return ntrees, nopen, err
  }

  fmt.Printf("Slope %d x %d: I dodged %d trees and hit %d.\n",
    slope, speed, nopen, ntrees)
  return ntrees, nopen, err
}

func Main(input_path string, verbose bool, args []string) error {
  file, err := os.Open(input_path)
  if err != nil {
    return err
  }
  defer file.Close()

  var ntrees int
  trees := make([]int, 5)
  treenum := 0
  slow_slopes := []int{1, 3, 5, 7}
  for _, slope := range slow_slopes {
    ntrees, _, err = do_sled(file, slope, 1)
    if err != nil {
      return err
    }
    trees[treenum] = ntrees
    treenum++
  }

  ntrees, _, err = do_sled(file, 1, 2)
  if err != nil {
    return err
  }
  trees[treenum] = ntrees

  fmt.Printf("Product: %d\n", p01.Product(trees))

  return nil
}
