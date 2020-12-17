package p11

import (
  "fmt"
  "github.com/fritzr/advent2020/util"
)

type SeatMap struct {
  seats [][]byte
  width int
  nEmpty int
  nOccupied int
  nFloor int
}

const SEAT_ERROR = byte(' ')
const SEAT_FLOOR = byte('.')
const SEAT_EMPTY = byte('L')
const SEAT_OCCUPIED = byte('#')

func BCount(arr []byte, c byte) (n int) {
  for _, a := range arr {
    if a == c {
      n++
    }
  }
  return n
}

func NewSeatMap(lines []string) *SeatMap {
  m := new(SeatMap)
  m.seats = make([][]byte, len(lines))
  for row, line := range lines {
    m.seats[row] = []byte(line)
  }
  if len(lines) > 0 {
    m.width = len(lines[0])
  } else {
    m.width = 0
  }
  // Count how many of stuff we have.
  for _, row := range m.seats {
    m.nEmpty += BCount(row, SEAT_EMPTY)
    m.nOccupied += BCount(row, SEAT_OCCUPIED)
  }
  m.nFloor = len(m.seats) * m.width - (m.nEmpty + m.nOccupied)
  return m
}

type Coordinate struct {
  row int
  col int
}

func (c *Coordinate) Add(row int, col int) {
  c.row += row
  c.col += col
}

func (c *Coordinate) Eq(o *Coordinate) bool {
  return c.row == o.row && c.col == o.col
}

func (m *SeatMap) At(c *Coordinate) byte {
  if c.row < 0 || c.row >= len(m.seats) || c.col < 0 || c.col >= m.width {
    return SEAT_ERROR
  }
  return m.seats[c.row][c.col]
}

// Simulate one step of people occupying seats.
//
// Returns true if anything changed.
// The stand() and sit() functions return whether an occupied seat becomes
// empty, or an empty seat becomes occupied, respectively.
func (m *SeatMap) fill(occupied func(*SeatMap, int, int) int,
                       stand func(*SeatMap, int) bool,
                       sit func(*SeatMap, int) bool) bool {
  // Apply rules simultaneously by storing changes before comitting them.
  newOccupied := make([]Coordinate, 0, m.Empty())
  newEmpty := make([]Coordinate, 0, m.Occupied())
  for rowIndex, row := range m.seats {
    for colIndex, status := range row {
      if status != SEAT_FLOOR {
        n := occupied(m, rowIndex, colIndex)
        // Some empty seats may become occupied.
        if status == SEAT_EMPTY && sit(m, n) {
          newOccupied = append(newOccupied, Coordinate{rowIndex, colIndex})
        // Some .
        } else if status == SEAT_OCCUPIED && stand(m, n) {
          newEmpty = append(newEmpty, Coordinate{rowIndex, colIndex})
        }
      }
    }
  }
  // Commit seat changes all at once.
  for _, oPos := range newOccupied {
    m.seats[oPos.row][oPos.col] = SEAT_OCCUPIED
  }
  for _, ePos := range newEmpty {
    m.seats[ePos.row][ePos.col] = SEAT_EMPTY
  }
  // Update counts.
  m.nOccupied += len(newOccupied) - len(newEmpty)
  m.nEmpty += len(newEmpty) - len(newOccupied)
  // Return true as long as something changed.
  return len(newOccupied) > 0 || len(newEmpty) > 0
}

func (m *SeatMap) countVisibleFrom(row int, col int, rd int, cd int) (n int) {
  pos := Coordinate{row + rd, col + cd}
  for ; m.At(&pos) == SEAT_FLOOR;
        pos.Add(rd, cd) {
  }
  if m.At(&pos) == SEAT_OCCUPIED {
    return 1
  }
  return 0
}

// Number of visible (line-of-sight) occupied seats.
func (m *SeatMap) nVisible(row int, col int) (n int) {
  // Cardinal directions.
  n += m.countVisibleFrom(row, col,  1,  0) // down
  n += m.countVisibleFrom(row, col, -1,  0) // up
  n += m.countVisibleFrom(row, col,  0, -1) // left
  n += m.countVisibleFrom(row, col,  0,  1) // right
  // Diagonals in all four directions.
  n += m.countVisibleFrom(row, col, -1, -1) // upper left
  n += m.countVisibleFrom(row, col, -1,  1) // upper right
  n += m.countVisibleFrom(row, col,  1, -1) // lower left
  n += m.countVisibleFrom(row, col,  1,  1) // lower right
  return n
}

// Number of adjacent occupied seats.
func (m *SeatMap) nAdjacent(seatRow int, seatCol int) (n int) {
  for rowIndex := seatRow - 1;
      rowIndex < seatRow + 2 && rowIndex < len(m.seats);
      rowIndex++ {
    for colIndex := seatCol - 1;
        colIndex < seatCol + 2 && colIndex < m.width;
        colIndex++ {
      if (rowIndex >= 0 && colIndex >= 0 && (
          m.seats[rowIndex][colIndex] == SEAT_OCCUPIED && (
            rowIndex != seatRow || colIndex != seatCol))) {
        n++
      }
    }
  }
  return n
}

// Whether to stand based on how many occupied seats are visible.
func (m *SeatMap) standVisible(visible int) bool {
  return visible >= 5
}

// Whether to stand based on how many occupied seats are adjacent.
func (m *SeatMap) standAdjacent(adjacent int) bool {
  return adjacent >= 4
}

// Whether to sit based on how many occupied seats we know about.
func (m *SeatMap) sit(adjacent int) bool {
  return adjacent == 0
}

// Simulate sitting/standing up by counting adjacent occupied seats.
func (m *SeatMap) FillAdjacent() bool {
  return m.fill((*SeatMap).nAdjacent,
                (*SeatMap).standAdjacent,
                (*SeatMap).sit)
}

// Simulate sitting/standing up by counting occupied seats in line of sight.
func (m *SeatMap) FillVisible() bool {
  return m.fill((*SeatMap).nVisible,
                (*SeatMap).standVisible,
                (*SeatMap).sit)
}

// Number of floor spaces.
func (m *SeatMap) Floor() int { return m.nFloor }
// Number of empty seats.
func (m *SeatMap) Empty() int { return m.nEmpty }
// Number of occupied seats.
func (m *SeatMap) Occupied() int { return m.nOccupied }

func Part1(lines []string) (int, int) {
  // Fill until nothing changes.
  seatMap := NewSeatMap(lines)
  nSteps := 0
  for seatMap.FillAdjacent() {
    nSteps++
  }
  return nSteps, seatMap.Occupied()
}

func Part2(lines []string) (int, int) {
  // Fill until nothing changes.
  seatMap := NewSeatMap(lines)
  nSteps := 0
  for seatMap.FillVisible() {
    nSteps++
  }
  return nSteps, seatMap.Occupied()
}

func Main(input_path string, verbose bool, args []string) error {
  lines, err := util.ReadLinesFromFile(input_path)
  if err != nil {
    return err
  }

  aSteps, aOccupied := Part1(lines)
  fmt.Printf("After %d adjacency steps, there are %d occupied seats.\n",
    aSteps, aOccupied)

  vSteps, vOccupied := Part2(lines)
  fmt.Printf("After %d visibility steps, there are %d occupied seats.\n",
    vSteps, vOccupied)

  return nil
}
