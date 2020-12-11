package util

import (
  "strings"
)

func first_not_of(haystack []byte, hay byte) int {
  for idx, c := range haystack {
    if c != hay {
      return idx // needle
    }
  }
  return -1
}

// SplitFunc for a bufio.Scanner which splits input into groups of lines
// separated by groups of blank lines.
func ScanLineGroups(data []byte, atEOF bool) (advance int, token []byte, err error) {
  var (consecutive int
       start int
       end int
       c byte)

  // Capture until we see consecutive empty lines.
  for advance, c = range data {
    if c == '\n' {
      if consecutive == 0 {
        end = advance
      }
      consecutive++
    } else {
      if (consecutive > 1) {
        break
      }
      consecutive = 0
      end = start
    }
  }

  // Didn't find a complete token, expand the buffer.
  if (consecutive < 2 && !atEOF) {
    return 0, nil, nil
  }

  // Found a token (maybe).
  if (end > start) {
    token = data[start:end]
  }

  // Eat trailing newlines.
  next := first_not_of(data[advance:], '\n')
  if next < 0 {
    advance = len(data)
  } else {
    advance += next
  }

  return advance, token, err
}

func Product(numbers []int) int {
  result := numbers[0]
  for _, value := range numbers[1:] {
    result *= value
  }
  return result
}

func StringIsSubset(subset string, superset string) bool {
  return 0 > strings.IndexFunc(subset, func(r rune) bool {
    return strings.IndexRune(superset, r) < 0 })
}
