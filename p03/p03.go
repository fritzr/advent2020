package p03

import (
  "fmt"
  "io"
  "bufio"
  "os"
)

func get_line_width(r io.Reader) (int, error) {
  lineReader := bufio.NewReader(r)
  line, _, err := lineReader.ReadLine()
  return len(line), err
}

func TobogganSled(rows *os.File, slope_width int) (int, int, error) {
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

  for nread > 0 && err == nil {
    // need to seek to end of line, then to the next line position
    line_remaining := line_width - line_pos
    line_pos = (line_pos + slope_width) % line_width

    // right N, down 1 (line width)
    rows.Seek(int64(line_remaining + line_pos), os.SEEK_CUR)

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

func Main(input_path string, verbose bool, args []string) error {
  file, err := os.Open(input_path)
  if err != nil {
    return err
  }
  defer file.Close()

  var (ntrees int; nopen int)
  ntrees, nopen, err = TobogganSled(file, 3)
  if err != nil {
    return err
  }

  fmt.Printf("I dodged %d trees and hit %d.\n", nopen, ntrees)
  return nil
}
