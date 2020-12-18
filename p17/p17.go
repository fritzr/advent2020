package p17

import (
  "fmt"
  "strconv"
  "strings"
  "errors"
  "github.com/fritzr/advent2020/util"
)

var gVerbose = false

type pocketDimensionKey uint64

type PocketDimension struct {
  // DOK (dictionary of keys) Sparse Matrix for N dimensions.
  // Map coordinates like (x,y,z,...) to cell state.
  // We use a string key as a substitute for using a variable integer array as
  // a map key, which ought to be faster but is not supported by Go maps.
  // We could potentially use bitwise ops to embed a custom hash
  cells map[pocketDimensionKey]bool
  ndim int
  //nactive int

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
      panic(fmt.Sprintf("coordinate %d too large for current hash", magnitude))
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
    return -int(magnitude &^ (1 << signBitPos))
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
    coords[index] = signDecode((uHash & mask) >> nshift, nbits - 1)
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

func (d *PocketDimension) isActiveHash(hash pocketDimensionKey) bool {
  return d.cells[hash]
}

func (d *PocketDimension) IsActive(coords []int) bool {
  return d.isActiveHash(pdHash(coords, d.ndim))
}

func (d *PocketDimension) activateHash(hash pocketDimensionKey) {
  /*
  if !d.cells[hash] {
  */
    d.cells[hash] = true
  /*
    d.nactive++
  }
  */
}

func (d *PocketDimension) Activate(coords []int) {
  d.activateHash(pdHash(coords, d.ndim))
}

func (d *PocketDimension) deactivateHash(hash pocketDimensionKey) {
  // Deleting inactive cells makes it far easier to count and iterate over
  // all active cells, and may help memory usage.
  /*
  if d.cells[hash] {
  */
    delete(d.cells, hash)
  /*
    d.nactive--
  }
  */
}

func (d *PocketDimension) Deactivate(coords []int) {
  d.deactivateHash(pdHash(coords, d.ndim))
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
func (d *PocketDimension) visitPartialNeighbors(origin []int,
                                                fixedCoords []int,
                                                f func(coords []int)) {
  if len(origin) == len(fixedCoords) {
    if !SliceEqual(origin, fixedCoords) {
      f(fixedCoords)
    }
  } else {
    varyCoord := origin[len(fixedCoords)]
    fixedCoords = append(fixedCoords, varyCoord - 1)
    lastIndex := len(fixedCoords)-1

    // fixedCoords[lastIndex] = varyCoord - 1
    d.visitPartialNeighbors(origin, fixedCoords, f)
    fixedCoords[lastIndex] = varyCoord
    d.visitPartialNeighbors(origin, fixedCoords, f)
    fixedCoords[lastIndex] = varyCoord + 1
    d.visitPartialNeighbors(origin, fixedCoords, f)
  }
}

// Visit neighbor (adjacent) cells.
//
// Does not visit the origin cell itself.
func (d *PocketDimension) VisitNeighbors(coords []int, f func(coords []int)) {
  d.visitPartialNeighbors(coords, []int{}, f)
}

// Visit all active cells in no particular order.
func (d *PocketDimension) VisitActive(f func(coords []int)) {
  for hash, active := range d.cells {
    // assert(active) -- enforced by Deactivate()
    if !active { panic("we are expecting deactivated cells to be deleted") }
    // Reverse (decode) the hash key to obtain the coordinates.
    f(pdRHash(hash, d.ndim))
  }
}

func (d *PocketDimension) ActiveNeighbors(coords []int) int {
  count := 0
  d.VisitNeighbors(coords, func(neighbor []int) {
    if d.IsActive(neighbor) {
      count += 1
    }
  })
  return count
}

func coordStr(coords []int) string {
  var s strings.Builder
  s.WriteByte('(')
  for _, coord := range coords[:len(coords)-1] {
    s.WriteString(strconv.Itoa(coord))
    s.WriteByte(',')
  }
  s.WriteString(strconv.Itoa(coords[len(coords)-1]))
  s.WriteByte(')')
  return s.String()
}

// Simulate one cycle.
func (d *PocketDimension) Simulate() {
  const STABLE = 1
  const ACTIVATE = 2
  const DEACTIVATE = 3

  // To avoid infinite recursion, we mark PENDING every node before we visit
  // its neighbors.
  const PENDING = 4

  // Store instructions for every cell we end up visiting.
  // These are all executed once after we visit the cells.
  exec := make(map[pocketDimensionKey]int)

  var stateStr func(int) string
  if gVerbose {
    stateStr = func(state int) string {
      return []string{"NONE", "STABLE", "ACTIVE", "INACTIVE", "PENDING"}[state]
    }
  }

  // Visit every active cell and their neighbors.
  // Though there are infinitely many cells, only active cells and the
  // inactive cells adjacent to them can ever change state.
  // We must be very careful not to infinitely recurse.
  d.VisitActive(func(active []int) {
    activeHash := pdHash(active, d.ndim)
    // Only process this cell if we haven't already processed it.
    if exec[activeHash] == 0 {
      exec[activeHash] = PENDING

      if gVerbose {
        fmt.Printf("  visiting   active %s (hash=%x)\n", coordStr(active), activeHash)
      }

      // Count active neighbors of the active cell.
      activeN := 0
      d.VisitNeighbors(active, func(n []int) {
        nHash := pdHash(n, d.ndim)
        if d.isActiveHash(nHash) {
          activeN++
        }

        // While we're at it, visit any inactive neighbor cell n (once).
        if exec[nHash] == 0 && !d.isActiveHash(nHash) {
          exec[nHash] = PENDING

          if gVerbose {
            fmt.Printf("    visiting inactive %s (hash=%x)\n",
              coordStr(n), nHash)
          }

          // To do that, count each active neighbor m of n.
          nN := 0
          d.VisitNeighbors(n, func(m []int) {
            if d.IsActive(m) {
              nN++
            }
          })
          // Inactive cells activate with exactly three neighbors.
          if nN == 3 {
            exec[nHash] = ACTIVATE
          } else {
            exec[nHash] = STABLE
          }
          if gVerbose {
            fmt.Printf("    ... %x state: inactive => %s (n=%d)\n",
              nHash, stateStr(exec[nHash]), nN)
          }
        }
      })

      // Active cells only remain active with 2 or 3 active neighbors.
      if activeN != 2 && activeN != 3 {
        exec[activeHash] = DEACTIVATE
      } else {
        exec[activeHash] = STABLE
      }
      if gVerbose {
        fmt.Printf("  ... %x state: active => %s (n=%d)\n",
          activeHash, stateStr(exec[activeHash]), activeN)
      }
    }
  })

  // Now we have a set of instructions to execute. Do it.
  //
  // We use panics as assertions to indicate such a state is not possible.
  for hash, state := range exec {
    switch state {
    case STABLE: // Do nothing.
    case ACTIVATE:
      if d.isActiveHash(hash) { panic("activating an active cell!") }
      d.activateHash(hash)
    case DEACTIVATE:
      if !d.isActiveHash(hash) { panic("deactivating an inactive cell!") }
      d.deactivateHash(hash)
    case PENDING: panic("cell state was not resolved!")
    }
  }
}

func (d *PocketDimension) ActiveCount() int {
  return len(d.cells)
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
    s.WriteString("  ")
    s.WriteString(coordStr(coords))
    nLine += 1
  })
  return s.String()
}

func Usage() {
  fmt.Println("usage: advent2020 17 [main opts...] [-n N]")
  fmt.Println("")
  fmt.Println("If N is given, run N iterations of the simulation (default 6).")
}

func Main(input_path string, verbose bool, args []string) error {
  gVerbose = verbose

  iterations := 6
  if len(args) > 0 {
    if args[0] == "-h" || args[0] == "--help" {
      Usage()
      return nil
    }
    if args[0] == "-n" {
      if len(args) < 2 {
        return errors.New("-n requires an argument")
      }
      var err error
      iterations, err = strconv.Atoi(args[1])
      if err != nil {
        return err
      }
    }
  }

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
    fmt.Println("Extents:")
    for dimNum, extents := range dim.GetExtents() {
      fmt.Printf("  (dim%d) [min=%d, max=%d]\n", dimNum, extents[0], extents[1])
    }
  }

  // Simulate 6 times and count active cells.
  for n := 0; n < iterations; n++ {
    dim.Simulate()
    if verbose {
      fmt.Printf("After step %d, actives are:\n%s\n", n + 1, dim.ActiveStr())
    }
  }
  fmt.Printf("After %d steps, there are %d active cells.\n",
    iterations, dim.ActiveCount())

  return nil
}
