
package p13

import (
  "fmt"
  "os"
  "io"
  "bufio"
  "strconv"
  "strings"
)

type BusSchedule struct {
  active []int
  inactive map[int]bool
}

func NewBusSchedule(sched string) (*BusSchedule, error) {
  b := new(BusSchedule)
  b.active = make([]int, 0, 64)
  b.inactive = make(map[int]bool, 128)
  for index, strValue := range strings.Split(sched, ",") {
    if strValue == "x" {
      b.inactive[index] = true
    } else {
      value, err := strconv.Atoi(strValue)
      if err != nil {
        return nil, err
      }
      b.active = append(b.active, value)
    }
  }
  return b, nil
}

func (b *BusSchedule) NextAvailable(time int) map[int]int {
  nextAvailable := make(map[int]int, len(b.active))
  for _, active := range b.active {
    nextAvailable[active] = active * (1 + time / active)
  }
  return nextAvailable
}

func ReadSchedule(input io.Reader) (timestamp int, s *BusSchedule, err error) {
  scanner := bufio.NewScanner(input)
  scanner.Split(bufio.ScanLines)
  if scanner.Scan() {
    timestamp, err = strconv.Atoi(scanner.Text())
  }
  if err == nil && scanner.Scan() {
    s, err = NewBusSchedule(scanner.Text())
  }
  if err == nil {
    err = scanner.Err()
  }
  return timestamp, s, err
}

func ReadScheduleFromFile(path string) (int, *BusSchedule, error) {
  file, err := os.Open(path)
  if err != nil {
    return -1, nil, err
  }
  defer file.Close()
  return ReadSchedule(file)
}

func Main(input_path string, verbose bool, args []string) error {
  timestamp, schedule, err := ReadScheduleFromFile(input_path)
  if err != nil {
    return err
  }

  earliestBus := -1
  earliestWaitTime := -1
  nextAvailable := schedule.NextAvailable(timestamp)
  for nextBus, nextTime := range nextAvailable {
    if verbose {
      fmt.Printf("  %d arrives next at %d\n", nextBus, nextTime)
    }
    wait := nextTime - timestamp
    if earliestWaitTime < 0 || wait < earliestWaitTime {
      earliestBus = nextBus
      earliestWaitTime = wait
    }
  }

  earliestBusTime := nextAvailable[earliestBus]
  fmt.Printf("Next bus after %d is %d, arriving at %d (wait time %d).\n",
    timestamp, earliestBus, earliestBusTime, earliestWaitTime)
  fmt.Printf("    %d x %d = %d\n",
    earliestBus, earliestWaitTime, earliestBus*earliestWaitTime)

  return nil
}
