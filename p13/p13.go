
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
  // List of bus IDs in the same order as the input.
  buses []int
  // Map bus ID to constraing position in the schedule.
  busOffsets map[int]int
}

func NewBusSchedule(sched string) (*BusSchedule, error) {
  b := new(BusSchedule)
  b.buses = make([]int, 0, 16)
  b.busOffsets = make(map[int]int, 16)
  for index, strValue := range strings.Split(sched, ",") {
    if strValue != "x" {
      busId, err := strconv.Atoi(strValue)
      if err != nil {
        return nil, err
      }
      b.buses = append(b.buses, busId)
      b.busOffsets[busId] = index
    }
  }
  return b, nil
}

func (b *BusSchedule) NextAvailable(time int) map[int]int {
  nextAvailable := make(map[int]int, len(b.buses))
  for _, busId := range b.buses {
    nextAvailable[busId] = busId * (1 + time / busId)
  }
  return nextAvailable
}

func (b *BusSchedule) ConstrainedTime() int {
  // Find the time t for which bus x departs at t + busOffsets[x].
  // TODO
  return -1
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

  // Part 1: find the earliest bus after the given timestamp.
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

  // Part 2: find timestamp which matches the schedule.
  constrainedTime := schedule.ConstrainedTime()
  fmt.Printf("Schedule matches constraints starting at %d.\n", constrainedTime)

  return nil
}
