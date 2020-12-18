package p14

import (
  "fmt"
  "errors"
  "strings"
  "strconv"
  "github.com/fritzr/advent2020/util"
)

type BitSystem struct {
  setMask uint64
  clrMask uint64
  mem map[uint64]uint64
}

func (s *BitSystem) Mask(value uint64) uint64 {
  return (value | s.setMask) &^ s.clrMask
}

func (s *BitSystem) Write(address uint64, value uint64) {
  s.mem[address] = s.Mask(value)
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
  var sum uint64
  for _, value := range s.mem {
    sum += value
  }
  return sum
}

func NewBitSystem() *BitSystem {
  s := new(BitSystem)
  s.mem = make(map[uint64]uint64)
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

func parseInsns(lines []string) ([]BitInsn, error) {
  fieldList := make([]BitInsn, 0, len(lines))
  for _, line := range lines {
    fields := strings.Split(line, " = ")
    if len(fields) != 2 {
      return fieldList, errors.New(
        fmt.Sprintf("wrong number of fields for instruction '%s'", line))
    }
    // Memory write: mem[INDEX] = VALUE
    if strings.HasPrefix(fields[0], "mem[") {
      indexStr := fields[0][4:len(fields[0])-1]
      index, err := strconv.Atoi(indexStr)
      if err != nil {
        return fieldList, err
      }
      value, err := strconv.Atoi(fields[1])
      fieldList = append(fieldList,
        BitInsn{INSN_WRITE, uint64(index), uint64(value)})
    } else if fields[0] == "mask" {
      // Mask update: mask = XXXXX1XXX0XX...
      setMask, clrMask, err := parseMask(fields[1])
      if err != nil {
        return fieldList, err
      }
      fieldList = append(fieldList,
        BitInsn{INSN_MASK, setMask, clrMask})
    } else {
      return fieldList, errors.New(
        fmt.Sprintf("unrecognized mnemonic '%s'", fields[0]))
    }
  }
  return fieldList, nil
}

func Main(input_path string, verbose bool, args []string) error {
  lines, err := util.ReadLinesFromFile(input_path)
  if err != nil {
    return err
  }

  insns, err2 := parseInsns(lines)
  if err2 != nil {
    return err2
  }

  s := NewBitSystem()
  err = s.ExecAll(insns)
  if err != nil {
    return err
  }

  fmt.Printf("After all %d instructions, the sum in memory is %d.\n",
    len(insns), s.MemorySum())

  return nil
}
