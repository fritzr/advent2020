package p17

import (
  "fmt"
  "strconv"
  "strings"
  "github.com/fritzr/advent2020/util"
)

type pocketDimensionKey uint64

type PocketDimension struct {
  // DOK (dictionary of keys) Sparse Matrix for N dimensions.
  // Map coordinates like (x,y,z,...) to cell state.
  // We use a string key as a substitute for using a variable integer array as
  // a map key, which ought to be faster but is not supported by Go maps.
  // We could potentially use bitwise ops to embed a custom hash
  cells map[pocketDimensionKey]bool
  ndim int
  nactive int

  // The extents {min, max} of coordinates ever activated in each dimension.
  // This will grow monotonically away from zero.
  // extents [][2]int
}

func sgn(val int) uint64 {
  if val < 0 {
    return 1
  }
  return 0
}

// We partition the bits of the integer hash by coordinate.
// For example, in a 2-dimensional coordinate system, the hash consists of
// a 32-bit x coordinate and a 32-bit y coordinate; in 3 dimensions, the
// hash consists of |_63/3_|=21 bits for each dimension, etc...
func pdHash(coords []int, ndim int) pocketDimensionKey {
  nbits := 64 / ndim
  // We store the number's magnitude and sign.
  mask := uint64((1 << (nbits - 1) - 1))
  ulimit := mask
  nshift := 0
  var hash uint64
  for _, coord := range coords[:ndim] {
    magnitude := uint64(util.IAbs(coord))
    if magnitude > ulimit {
      panic("coordinate too large for current hash")
    }
    hash |= (magnitude << nshift) & mask
    mask <<= nbits
    nshift += nbits
    // Store the sign bit as the uppermost bit in the nbits mask.
    hash |= sgn(coord) << (nshift - 1)
  }
  return pocketDimensionKey(hash)
}

func signDecode(magnitude uint64, signBitPos int) int {
  if magnitude & (1 << signBitPos) != 0 {
    return -int(magnitude)
  }
  return int(magnitude)
}

// Reverse hash.
// Decode the numbers in the hash.
func pdRHash(hash pocketDimensionKey, ndim int) (coords []int) {
  coords = make([]int, ndim)
  nbits := 64 / ndim
  mask := uint64((1 << nbits) - 1)
  nshift := 0
  uHash := uint64(hash)
  for index := 0; index < ndim; index++ {
    coords[index] = signDecode((uHash & mask) >> nshift, nshift + nbits - 1)
    mask <<= nbits
    nshift += nbits
  }
  return coords
}

func NewPocketDimension(ndim int) *PocketDimension {
  d := new(PocketDimension)
  d.cells = make(map[pocketDimensionKey]bool)
  d.ndim = ndim
  // d.extents = make([][2]int, ndim)
  return d
}

func (d *PocketDimension) IsActive(coords []int) bool {
  return d.cells[pdHash(coords, d.ndim)]
}

func (d *PocketDimension) Activate(coords []int) {
  hash := pdHash(coords, d.ndim)
  if !d.cells[hash] {
    d.cells[hash] = true
    d.nactive++
  }
}

func (d *PocketDimension) Deactivate(coords []int) {
  // Is deleting too slow?
  // It makes it far easier to iterate over all active cells.
  hash := pdHash(coords, d.ndim)
  if d.cells[hash] {
    delete(d.cells, hash)
    d.nactive--
  }
}

func SliceEqual(s1 []int, s2 []int) bool {
  if len(s1) != len(s2) { return false }
  for idx := 0; idx < len(s1); idx++ {
    if s1[idx] != s2[idx] {
      return false
    }
  }
  return true
}

// Count active cells adjacent to the chosen hyperplane.
//
// The hyperplane is formed by fixing the coordinates in the latter argument.
// Varying the next non-fixed coordinate from the origin in {X-1, X, X+1}
// forms the adjacent hyperplanes.
//
// In the base case, all coordinates are fixed and we can check if it is active.
func (d *PocketDimension) countActive(origin []int, fixedCoords []int) int {
  var count int
  if len(origin) == len(fixedCoords) {
    if !SliceEqual(origin, fixedCoords) && d.IsActive(fixedCoords) {
      return 1
    }
    return 0
  } else {
    varyCoord := origin[len(fixedCoords)]
    fixedCoords = append(fixedCoords, varyCoord - 1)
    lastIndex := len(fixedCoords)-1

    // fixedCoords[lastIndex] = varyCoord - 1
    count += d.countActive(origin, fixedCoords)
    fixedCoords[lastIndex] = varyCoord
    count += d.countActive(origin, fixedCoords)
    fixedCoords[lastIndex] = varyCoord + 1
    count += d.countActive(origin, fixedCoords)
  }
  return count
}

// Count active neighbors.
// We need to visit every adjacent cell.
func (d *PocketDimension) ActiveNeighbors(coords []int) int {
  return d.countActive(coords, []int{})
}

// Visit active cells in no particular order.
func (d *PocketDimension) VisitActive(f func(coords []int)) {
  for hash, active := range d.cells {
    // assert(active) -- enforced by Deactivate()
    if !active { panic("we are expecting deactivated cells to be deleted") }
    // Reverse (decode) the hash key to obtain the coordinates.
    f(pdRHash(hash, d.ndim))
  }
}

// Simulate one cycle.
func (d *PocketDimension) Simulate() {
  // TODO
}

func (d *PocketDimension) ActiveCount() int {
  return d.nactive
}

func (d *PocketDimension) GetExtents() [][2]int {
  extents := make([][2]int, d.ndim)
  d.VisitActive(func(coords []int) {
    for dim, val := range coords {
      if val < extents[dim][0] {
        extents[dim][0] = val
      }
      if val > extents[dim][1] {
        extents[dim][1] = val
      }
    }
  })
  return extents
}

// String representation for debugging.
func (d *PocketDimension) ActiveStr() string {
  var s strings.Builder
  const maxLine = 8
  nLine := 0
  d.VisitActive(func(coords []int) {
    if nLine == maxLine {
      s.WriteString("\n")
      nLine = 0
    }
    s.WriteString("  (")
    for _, coord := range coords[:len(coords)-1] {
      s.WriteString(strconv.Itoa(coord))
      s.WriteByte(',')
    }
    s.WriteString(strconv.Itoa(coords[len(coords)-1]))
    s.WriteString(")")
    nLine += 1
  })
  return s.String()
}

func Main(input_path string, verbose bool, args []string) error {
  lines, err := util.ReadLinesFromFile(input_path)
  if err != nil {
    return err
  }

  // Part 1: Activate cells from the plane specified in the input.
  dim := NewPocketDimension(3)
  for rowCoord, line := range lines {
    for colCoord, c := range []byte(line) {
      if c != '.' {
        dim.Activate([]int{rowCoord, colCoord, 0})
      }
    }
  }
  fmt.Printf("There are initially %d active cells.\n", dim.ActiveCount())
  if verbose {
    fmt.Println(dim.ActiveStr())
  }

  // Simulate 6 times and count active cells.
  const nSteps = 6
  for n := 0; n < nSteps; n++ {
    dim.Simulate()
  }
  fmt.Printf("After %d steps, there are %d active cells.\n",
    nSteps, dim.ActiveCount())

  return nil
}
