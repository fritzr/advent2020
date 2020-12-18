
package p13

import (
  "fmt"
  "os"
  "io"
  "bufio"
  "strconv"
  "strings"
)

var gVerbose = false

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

func (b *BusSchedule) NextAvailable(time int64) map[int]int64 {
  nextAvailable := make(map[int]int64, len(b.buses))
  for _, busId := range b.buses {
    nextAvailable[busId] = int64(busId) * (1 + time / int64(busId))
  }
  return nextAvailable
}

// Compute the parameter T for integer n and m such that:
//   A * n = T and = B * m = T + O.
//
// In other words,
//   A * n = B * m - O.
//
// This function returns n. To obtain T itself, multiply the result by A.
func offsetMatchN(A int64, B int64, O int64) int64 {
  Q := O / A
  O = O % A

  // In the degenerate case where the desired offset is zero,
  // we can simply use n = B and m = A such that A * B = T = B * A.
  // Of course, with a non-zero offset divisible by A, we still need to
  // subtract by the quotient O / A.
  if O == 0 {
    return B - Q
  }

  // First find the quotient q and remainder r w.r.t B.
  // We will further decompose A * n into k * (B + r).
  q := B / A
  r := B % A

  // Find a = k * r = O mod A, the first multiple of r congruent to O modulo A.
  // XXX can we do this mathematically or must we iterate?
  k := int64(1)
  a := r
  for (a % A) != O {
    a += r
    k += 1
  }

  // Multiply the result by the original quotient of B / A,
  // then add (a / A) to obtain the first such number divisible by A.

  // Note that we tracked a = r * k already.
  // The formula for n is n = (B / A) * k + ((r * k) / A).
  // In other words, A * n = k * (B + r) = B * m + O such that k * r = O mod A.
  // We also need to subtract Q = O / A in case the requested offset is
  // larer than the root A.
  n := q * k + (a / A) - Q

  // Now one could compute m from the definition simply as:
  // m := (A * n + O) / B
  return n
}

func (b *BusSchedule) ConstrainedTime() int64 {
  lcm := int64(b.buses[0])
  t := int64(0)

  // Find the sub-result t for A and B which defines the sequence:
  //
  //   { t, t + lcm(A, B), ..., t + lcm(A, B) * n }
  //
  for _, busId := range b.buses[1:] {
    if gVerbose {
      fmt.Printf("  %d + %d * n = T\n", t, lcm )
      fmt.Printf("  %d + %d * n = T\n", t + int64(b.busOffsets[busId]), busId)
      fmt.Printf("===================\n")
    }
    t += lcm * offsetMatchN(lcm, int64(busId), t + int64(b.busOffsets[busId]))
    lcm *= int64(busId)
  }

  return t
}

func (b *BusSchedule) ConstrainedTimeBruteForce(maxIterations int64) int64 {
  // Find the time t for which bus x departs at t + busOffsets[x] (for all x).
  bus0 := int64(b.buses[0])
  t := bus0
  n := int64(1)
  for n < maxIterations + 1 {
    valid := true
    /*
    if n == int64(152683) {
      fmt.Printf("%d: time=%d\n", n, t)
    }
    */
    for id, time := range b.NextAvailable(t) {
      /*
      if n == int64(152683) {
        fmt.Printf("    %d requires offset %d, arrives at offset %d (%d)\n",
          id, b.busOffsets[id], time - t, time)
      }
      */
      if int64(id) != bus0 && (time - t) != int64(b.busOffsets[id]) {
        valid = false
        break
      }
    }
    if valid {
      return t
    }
    n++
    t += bus0
  }
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
  gVerbose = verbose
  timestamp, schedule, err := ReadScheduleFromFile(input_path)
  if err != nil {
    return err
  }

  // Part 1: find the earliest bus after the given timestamp.
  time := int64(timestamp)
  earliestBus := -1
  earliestWaitTime := int64(-1)
  nextAvailable := schedule.NextAvailable(time)
  for nextBus, nextTime := range nextAvailable {
    if verbose {
      fmt.Printf("  %d arrives next at %d\n", nextBus, nextTime)
    }
    wait := nextTime - time
    if earliestWaitTime < 0 || wait < earliestWaitTime {
      earliestBus = nextBus
      earliestWaitTime = wait
    }
  }

  earliestBusTime := nextAvailable[earliestBus]
  fmt.Printf("Next bus after %d is %d, arriving at %d (wait time %d).\n",
    time, earliestBus, earliestBusTime, earliestWaitTime)
  fmt.Printf("    %d x %d = %d\n",
    earliestBus, earliestWaitTime, int64(earliestBus)*earliestWaitTime)

  // Part 2: find timestamp which matches the scheduled wait times.
  constrainedTime := schedule.ConstrainedTime()
  fmt.Printf("Schedule matches constraints starting at %d.\n", constrainedTime)
  return nil
}
