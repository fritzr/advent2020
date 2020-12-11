package p06

import (
  "io"
  "bufio"
  "os"
  "fmt"
  "unicode"
  "github.com/fritzr/advent2020/util"
)

type ResponseGroup struct {
  any map[byte]int
  all map[byte]int
  members int
}

func NewResponseGroup(data string) ResponseGroup {
  // If we got any responses, we got one more than the number of newlines,
  // since ScanLineGroups never returns trailing newlines, so initialize
  // members to 1.
  group := ResponseGroup{make(map[byte]int, 26), make(map[byte]int, 26), 1}
  for _, question := range data {
    // Record the questions which anyone answered.
    if !unicode.IsSpace(question) {
      group.any[byte(question)]++
    } else if question == '\n' {
      group.members++
    }
  }
  // If we didn't actually get any responses, reset members to zero.
  if len(group.any) == 0 {
    group.members = 0
  } else {
    // Record the questions which everyone answered.
    for question, count := range group.any {
      if count == group.members {
        group.all[question] = 1
      }
    }
  }
  return group
}

func ReadResponseGroups(input io.Reader) ([]ResponseGroup, error) {
  scanner := bufio.NewScanner(input)
  scanner.Split(util.ScanLineGroups)

  groups := make([]ResponseGroup, 0, 1024)
  for scanner.Scan() {
    groups = append(groups, NewResponseGroup(scanner.Text()))
  }
  return groups, scanner.Err()
}

func ReadResponseGroupsFromFile(path string) ([]ResponseGroup, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()
  return ReadResponseGroups(file)
}

func Main(input_path string, verbose bool, args []string) error {

  responses, err := ReadResponseGroupsFromFile(input_path)
  if err != nil {
    return err
  }

  fmt.Printf("Read responses from %d groups.\n", len(responses))

  any_sum := 0
  all_sum := 0
  for _, group := range responses {
    any_sum += len(group.any)
    all_sum += len(group.all)
  }

  // Part 1
  fmt.Printf("Sum of \"any\" response counts is: %d\n", any_sum)
  // Part 2
  fmt.Printf("Sum of \"all\" response counts is: %d\n", all_sum)

  return nil
}
