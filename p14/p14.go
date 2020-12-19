package p14

import (
  "fmt"
  "errors"
  "strings"
  "strconv"
  "math/bits"
  "github.com/fritzr/advent2020/util"
)

var gVerbose = false

type BitMemory interface {
  Write(address uint64, value uint64, setMask uint64, clrMask uint64)
  Sum() uint64
  Addresses() uint64
}

func sumStore(store map[uint64]uint64) uint64 {
  var sum uint64
  for _, value := range store {
    sum += value
  }
  return sum
}

// Flat sparse memory.
//
// Writes will mask the value being written.
type FlatMemory struct {
  store map[uint64]uint64
}

func NewFlatMemory() *FlatMemory {
  m := new(FlatMemory)
  m.store = make(map[uint64]uint64)
  return m
}

func (m *FlatMemory) Write(addr uint64, value uint64, set uint64, clr uint64) {
  m.store[addr] = (value | set) &^ clr
}

func (m *FlatMemory) Sum() uint64 {
  return sumStore(m.store)
}

func (m *FlatMemory) Addresses() uint64 {
  return uint64(len(m.store))
}

// Floating-address sparse memory.
//
// Writes will mask the address, potentially causing writes to many addresses.
// The don't-care bits cause the address to assume all possible combinations.
type FloatMemory struct {
  // To be general we can't do something cheeky like storing the mask as
  // a string and only expanding the addresses at Sum() time, because writes
  // to "different" addresses may cover the same address. For example,
  // writing "0b0*x1" and "0b0*0x" will both write to address 0b11.
  store map[uint64]uint64
  width int
}

func NewFloatMemory(maskWidth int) *FloatMemory {
  m := new(FloatMemory)
  m.store = make(map[uint64]uint64)
  m.width = maskWidth
  return m
}

func (m *FloatMemory) writeAll(addr uint64, value uint64, xBits []int) {
  if len(xBits) == 0 {
    if gVerbose {
      fmt.Printf("Writing %#x to %#x\n", value, addr)
    }
    m.store[addr] = value
  } else {
    nextBit := uint64(1 << xBits[0])
    nextBits := xBits[1:]
    if gVerbose {
      fmt.Printf(".. clearing %#x\n", nextBit)
    }
    m.writeAll(addr &^ nextBit, value, nextBits)
    if gVerbose {
      fmt.Printf(".. setting %#x\n", nextBit)
    }
    m.writeAll(addr | nextBit, value, nextBits)
  }
}

func (m *FloatMemory) Write(addr uint64, value uint64, set uint64, clr uint64) {
  // Don't care mask.
  fullMask := uint64((1 << m.width) - 1)
  set &= fullMask
  clr &= fullMask
  dontCare := fullMask &^ (set | clr)
  // Expand the write to all referenced addresses.
  xIndexes := make([]int, 0, m.width)
  var nshift int
  for dontCare != 0 && len(xIndexes) < m.width {
    bitIndex := bits.TrailingZeros64(dontCare)
    xIndexes = append(xIndexes, bitIndex + nshift)
    dontCare >>= bitIndex + 1
    nshift += bitIndex + 1
  }
  if gVerbose {
    fmt.Printf("write(%#x(%d), %#x(%d), %#x, %#x)",
      addr, addr, value, value, set, clr)
  }
  // We no longer clear according to the clear mask.
  addr |= set
  if gVerbose {
    fmt.Printf(" (masked=%#x(%d))\n", addr, addr)
    fmt.Printf("  Indexes = ")
    util.PrintArrayInline(xIndexes)
    fmt.Println()
  }
  m.writeAll(addr, value, xIndexes)
}

func (m *FloatMemory) Sum() uint64 {
  return sumStore(m.store)
}

func (m *FloatMemory) Addresses() uint64 {
  return uint64(len(m.store))
}

// Bit system.
//
// Can execute simple instructions such as INSN_WRITE and INSN_MASK.
// The backing Memory decides what to do with the current mask on writes.
type BitSystem struct {
  setMask uint64
  clrMask uint64
  mem interface{BitMemory}
}

func (s *BitSystem) Write(address uint64, value uint64) {
  s.mem.Write(address, value, s.setMask, s.clrMask)
}

func (s *BitSystem) SetMask(setMask uint64, clrMask uint64) {
  s.setMask = setMask
  s.clrMask = clrMask
}

func (s *BitSystem) Exec(insn *BitInsn) error {
  switch (insn.insn) {
  case INSN_WRITE: s.Write(insn.op1, insn.op2)
  case INSN_MASK: s.SetMask(insn.op1, insn.op2)
  default: return errors.New(fmt.Sprintf("invalid instruction %d", insn.insn))
  }
  return nil
}

func (s *BitSystem) ExecAll(insns []BitInsn) error {
  for _, insn := range insns {
    if err := s.Exec(&insn); err != nil {
      return err
    }
  }
  return nil
}

func (s *BitSystem) MemorySum() uint64 {
  return s.mem.Sum()
}

func NewBitSystem(memory interface{BitMemory}) *BitSystem {
  s := new(BitSystem)
  s.mem = memory
  return s
}

const INSN_MASK = uint(1)
const INSN_WRITE = uint(2)

type BitInsn struct {
  insn uint  // INSN_MASK or INSN_WRITE
  op1 uint64 // memory index or set mask
  op2 uint64 // write value or clear mask
}

func parseMask(maskStr string) (setMask uint64, clrMask uint64, err error) {
  // Read string characters.
  // This will build the masks in reverse order, which is perfect.
  if len(maskStr) > 64 {
    err = errors.New(fmt.Sprintf("mask expression too long (%d)", len(maskStr)))
  } else {
    for strIndex := 0; strIndex < len(maskStr); strIndex++ {
      switch maskStr[strIndex] {
      case '1': setMask |= 1
      case '0': clrMask |= 1
      default:  // don't care
      }
      // Don't shift out the last bit.
      if strIndex != len(maskStr) - 1 {
        setMask <<= 1
        clrMask <<= 1
      }
    }
  }
  return setMask, clrMask, err
}

func parseInsns(lines []string) ([]BitInsn, int, error) {
  fieldList := make([]BitInsn, 0, len(lines))
  var maskWidth int
  for _, line := range lines {
    fields := strings.Split(line, " = ")
    if len(fields) != 2 {
      return fieldList, 0, errors.New(
        fmt.Sprintf("wrong number of fields for instruction '%s'", line))
    }
    // Memory write: mem[INDEX] = VALUE
    if strings.HasPrefix(fields[0], "mem[") {
      indexStr := fields[0][4:len(fields[0])-1]
      index, err := strconv.Atoi(indexStr)
      if err != nil {
        return fieldList, 0, err
      }
      value, err := strconv.Atoi(fields[1])
      fieldList = append(fieldList,
        BitInsn{INSN_WRITE, uint64(index), uint64(value)})
    } else if fields[0] == "mask" {
      // Mask update: mask = XXXXX1XXX0XX...
      if len(fields[1]) > maskWidth {
        maskWidth = len(fields[1])
      }
      setMask, clrMask, err := parseMask(fields[1])
      if err != nil {
        return fieldList, 0, err
      }
      fieldList = append(fieldList,
        BitInsn{INSN_MASK, setMask, clrMask})
    } else {
      return fieldList, 0, errors.New(
        fmt.Sprintf("unrecognized mnemonic '%s'", fields[0]))
    }
  }
  return fieldList, maskWidth, nil
}

func doExec(sys *BitSystem, insns []BitInsn) error {
  err := sys.ExecAll(insns)
  if err != nil {
    return err
  }

  fmt.Printf("The sum of %d memory words is %d.\n",
    sys.mem.Addresses(), sys.MemorySum())
  return nil
}

func Main(input_path string, verbose bool, args []string) error {
  gVerbose = verbose
  lines, err := util.ReadLinesFromFile(input_path)
  if err != nil {
    return err
  }

  insns, maskWidth, err2 := parseInsns(lines)
  if err2 != nil {
    return err2
  }

  // Part 1: standard flat memory.
  flat := NewFlatMemory()
  s1 := NewBitSystem(flat)
  fmt.Printf("Flat memory: ")
  err = doExec(s1, insns)
  if err != nil {
    return err
  }

  // Part 2: special floating-address memory.
  floating := NewFloatMemory(maskWidth)
  s2 := NewBitSystem(floating)
  fmt.Printf("Floating memory: ")
  err = doExec(s2, insns)
  if err != nil {
    return err
  }

  return nil
}
