package p08

import (
  "fmt"
  "io"
  "errors"
  "strings"
  "strconv"
  "github.com/fritzr/advent2020/util"
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
  return s.Exec(s.insns[s.pc])
}

// Execute an instruction instead of the current instruction.
//
// First override the PC to be the given value.
func (s *Simulator) ExecAt(insn string, pc int) error {
  s.pc = pc
  ops := strings.Fields(insn)
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

// Execuate an instruction instead of the current instruction.
//
// The PC will updated according to the exec'd instruction.
// This essentially overrides the current instruction.
func (s *Simulator) Exec(insn string) error {
  return s.ExecAt(insn, s.pc)
}

func (s *Simulator) Insn(address int) string {
  if address >= 0 && address < len(s.insns) {
    return s.insns[address]
  }
  return ""
}

func (s *Simulator) Reset() {
  s.pc = 0
  s.accumulator = 0
}

func (s *Simulator) FindLoop() (int, int, error) {
  return s.FindLoopFrom(s.pc, s.insns[s.pc])
}

func (s *Simulator) FindLoopFrom(pc int, insn string) (int, int, error) {
  visited := make(map[int]bool)

  visited[pc] = true
  lastPc := pc
  if err := s.ExecAt(insn, pc); err != nil {
    return -1, -1, err
  }

  if gVerbose {
    fmt.Printf("Exec First: [%02d] %s (acc=%d)\n",
      s.pc, s.insns[s.pc], s.accumulator)
  }

  err := s.Step()
  for err == nil {
    if gVerbose {
      fmt.Printf("Exec Next: [%02d] %s (acc=%d)\n",
        s.pc, s.insns[s.pc], s.accumulator)
    }
    // Loop!
    if visited[s.pc] {
      return lastPc, s.pc, nil
    }
    lastPc = s.pc
    visited[s.pc] = true
    err = s.Step()
  }
  if err != nil {
    return -1, -1, err
  }

  return -1, -1, nil
}

func NewSimulator(insns []string) *Simulator {
  s := new(Simulator)
  s.insns = insns
  return s
}

func Part1(insns []string) (int, error) {
  sim := NewSimulator(insns)

  loopFrom, loopTo, err := sim.FindLoop()
  if err != nil {
    return sim.accumulator, err
  }
  if loopFrom < 0 || loopTo < 0 {
    return sim.accumulator, errors.New("No loop found!")
  }
  fmt.Printf("Loop: [%02d] %s => [%02d] %s\n",
    loopFrom, sim.insns[loopFrom], loopTo, sim.insns[loopTo])

  fmt.Printf("Accumulator: %d\n", sim.accumulator)
  return sim.accumulator, nil
}

func FixLoop(insns []string) (int, int, error) {
  sim := NewSimulator(insns)
  for pc, insn := range insns {
    if strings.HasPrefix(insn, "acc ") {
      continue
    }

    // See if either the regular instruction or its replacement cause a loop.
    var altInsn string
    if strings.HasPrefix(insn, "jmp ") {
      altInsn = "nop " + insn[4:]
    } else {
      altInsn = "jmp " + insn[4:]
    }

    for _, insn := range []string{insn, altInsn} {
      sim.insns[pc] = insn
      if gVerbose {
        fmt.Printf("Trying [%d] %s (acc=%d)...\n", pc, insn)
      }
      // TODO... we could probably do this smarter than running the
      // whole program each time.
      sim.Reset()
      from, to, err := sim.FindLoop()
      if err != nil {
        // Natural EOF, good!
        if err == io.EOF {
          err = nil // good!
        } else {
          pc = sim.pc // location of the error
        }
        return pc, sim.accumulator, err
      }
      // No loop... but this doesn't result in io.EOF.
      if from < 0 || to < 0 {
        return pc, sim.accumulator, nil
      }
      if gVerbose {
        fmt.Printf("... [%d] %s caused loop [%d] %s -> [%d] %s\n",
          pc, insn, from, insns[from], to, insns[to])
      }
    }

    // Restore the original instruction for the next attempt.
    sim.insns[pc] = insn
  }

  return -1, 0, errors.New("No fixable loops found")
}

func Part2(insns []string) (int, int, error) {
  fixedPc, acc, err := FixLoop(insns)
  if err != nil {
    return fixedPc, acc, err
  }
  fmt.Printf("Fixed loop at [%d] %s! Accumulator: %d\n",
    fixedPc, insns[fixedPc], acc)
  return fixedPc, acc, nil
}

func Main(input_path string, verbose bool, args []string) error {
  gVerbose = verbose

  insns, err := util.ReadLinesFromFile(input_path)
  if err != nil {
    return err
  }

  fmt.Printf("Loaded %d instructions\n", len(insns))

  if _, err = Part1(insns); err != nil {
    return err
  }

  if _, _, err = Part2(insns); err != nil {
    return err
  }

  return nil
}
