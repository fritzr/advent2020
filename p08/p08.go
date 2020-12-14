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
  return ExecAt(insn, s.pc)
}

func (s *Simulator) Insn(address int) string {
  if address >= 0 && address < len(s.insns) {
    return s.insns[address]
  }
  return string()
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

type pathEntry struct {
  // PC value of the instruction.
  address int
  // Whether this instruction is a branch point.
  branch bool
  // If this instruction is a branch point, whether the instruction's alternate
  // form has been attempted.
  alternate bool // stateNOP or stateJMP
}

type ProgramFixer struct {
  s *Simulator

  // Search path for a program which terminates.
  //
  // Each entry in the path represents an instruction which we've executed to
  // get here. If we detect an infinite loop, we walk back through the path
  // stack to the most recent branch point, then try the alternate branch path
  // (JMP to NOP or vice-versa). If both lead to an infinite loop, we cut that
  // path too, add it to the set of dead nodes, and walk back again.
  path []pathEntry
  pathIndex int

  // Set of dead nodes.
  //
  // If we ever reach a dead node, the program is guaranteed not to terminate.
  dead map[int]bool

  // Instructions visited.
  //
  // For each instruction I in the path, visited[I] is true.
  // This is essentially a fast way to determine if we're about to cause a
  // loop. When we're walking back in the path stack, we also need to clear
  // the visited flag for each instruction we pop to keep them in sync.
  visited []bool
}

func NewProgramFixer(s *Simulator) *ProgramFixer {
  f := new(ProgramFixer)
  f.s = s
  f.path = make([]pathEntry, len(f.insns))
  f.dead = make(map[int]bool, len(f.insns))
  f.visited = make([]bool, len(f.insns))
  return f
}

func isBranch(insn string) bool {
  return (
    strings.Fields(insn)[0] == "jmp" || strings.Fields(insn)[0] == "nop")
}

func (f *ProgramFixer) newPathEntry(addr int) pathEntry {
  return pathEntry{addr, isBranch(f.s.insns[addr]), false}
}

func (f *ProgramFixer) stepLoop(addr int) (bool, error) {
  if f.pathIndex >= len(f.path) || addr >= len(f.s.insns) {
    return false, io.EOF
  }

  // See if this instruction causes a loop...
  f.s.pc = addr
  f.path[f.pathIndex] = newPathEntry(addr)
  f.pathIndex++
  err = f.s.Step()
  if err != nil {
    return false, err
  }
  if f.visited[f.s.pc] {
    return markLoop(f.s.pc)
  }

  // If not, see if the next instruction causes a loop.
  var loops bool
  loops, err = stepLoop(f.s.pc)
  if err != nil {
    return false, err
  }

  // If yes, try the alternate instruction.
  if loops {
    f.pathIndex--
    if !f.path[f.pathIndex].alternate {
      f.s.ExecAt(addr)
    }
  }
}

func (f *ProgramFixer) Fix() error {
  _, err := stepLoop()
  return err
}

func Part1(sim *Simulator) error {
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

func Part2(sim *Simulator) {
  fixer := NewProgramFixer(sim)
  return fixer.Fix()
}

func Main(input_path string, verbose bool, args []string) error {
  gVerbose = verbose

  sim, err := ReadSimulatorFromFile(input_path)
  if err != nil {
    return err
  }

  if err = Part1(sim); err != nil {
    return err
  }

  if err = Part2(sim); err != nil {
    returne rr
  }

  return nil
}
