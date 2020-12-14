package p08

import (
  "fmt"
  "io"
  "os"
  "bufio"
  "errors"
  "strings"
  "strconv"
)

var gVerbose = false

type Simulator struct {
  insns []string
  accumulator int
  pc int
}

func (s *Simulator) acc(ops []string) (int, error) {
  value, err := strconv.Atoi(ops[0])
  if err != nil {
    return s.pc, err
  }
  if gVerbose {
    fmt.Printf("  ACC %d\n", value)
  }
  s.accumulator += value
  return s.pc + 1, nil
}

func (s *Simulator) nop(ops []string) (int, error) {
  if gVerbose {
    fmt.Printf("  NOP\n")
  }
  return s.pc + 1, nil
}

func (s *Simulator) jmp(ops []string) (int, error) {
  value, err := strconv.Atoi(ops[0])
  if gVerbose {
    fmt.Printf("  JMP %d\n", value)
  }
  if err != nil {
    return s.pc, err
  }
  return s.pc + value, err
}

// Op functions take a list of operands and return the new PC value,
// and an error if one occurs.
type opFunc func(*Simulator, []string) (int, error)

var opTable = map[string]opFunc {
  "acc": (*Simulator).acc,
  "nop": (*Simulator).nop,
  "jmp": (*Simulator).jmp }

func (s *Simulator) Jump(address int) error {
  if gVerbose {
    fmt.Printf("  PC %d => %d\n", s.pc, address)
  }
  if address >= 0 && address < len(s.insns) {
    s.pc = address
  } else {
    // natural termination
    if address == len(s.insns) {
      return io.EOF
    }
    return errors.New(fmt.Sprintf(
      "instruction target out of bounds at [%03d] %s", s.pc, s.insns[s.pc]))
  }
  return nil
}

func (s *Simulator) Step() error {
  if s.pc >= len(s.insns) {
    return errors.New(fmt.Sprintf("overflow at offset %d", s.pc))
  }
  if gVerbose {
    fmt.Printf("  STEP %d: insn=%s\n", s.pc, s.insns[s.pc])
  }
  ops := strings.Fields(s.insns[s.pc])
  if len(ops) == 0 {
    return errors.New(fmt.Sprintf("empty instruction at %d", s.pc))
  }
  opFunc := opTable[ops[0]]
  if opFunc == nil {
    return errors.New(fmt.Sprintf("unrecognized instruction '%s'", ops[0]))
  }
  newPc, err := opFunc(s, ops[1:])
  if err == nil {
    return s.Jump(newPc)
  }
  return err
}

func NewSimulator(input io.Reader) (*Simulator, error) {
  s := new(Simulator)
  s.insns = make([]string, 0, 650)
  scanner := bufio.NewScanner(input)
  scanner.Split(bufio.ScanLines)
  for scanner.Scan() {
    s.insns = append(s.insns, scanner.Text())
  }
  return s, scanner.Err()
}

func ReadSimulatorFromFile(path string) (*Simulator, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()
  return NewSimulator(file)
}

func Main(input_path string, verbose bool, args []string) error {
  gVerbose = verbose

  sim, err := ReadSimulatorFromFile(input_path)
  if err != nil {
    return err
  }

  // execution count
  xcount := make([]int, len(sim.insns))
  fmt.Printf("Loaded %d instructions\n", len(sim.insns))

  // Part 1: find the first loop.
  lastPc := -1
  if verbose {
    fmt.Printf("First: [%02d] %s (acc=%d)\n",
      sim.pc, sim.insns[sim.pc], sim.accumulator)
  }
  xcount[sim.pc] += 1
  err = sim.Step()
  for err == nil {
    if verbose {
      fmt.Printf("Next: [%02d] %s (acc=%d)\n",
        sim.pc, sim.insns[sim.pc], sim.accumulator)
    }
    // Loop detection.
    if xcount[sim.pc] > 0 {
      break
    }
    lastPc = sim.pc
    xcount[sim.pc] += 1
    err = sim.Step()
  }
  if err != nil {
    return err
  }

  fmt.Printf("Loop: [%02d] %s => [%02d] %s\n",
    lastPc, sim.insns[lastPc], sim.pc, sim.insns[sim.pc])

  fmt.Printf("Accumulator: %d\n", sim.accumulator)
  return nil
}
