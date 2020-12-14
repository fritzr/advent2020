package p12

import (
  "fmt"
  "io"
  "os"
  "bufio"
  "strconv"
  "errors"
)

type Direction struct {
  action byte
  value int
}

// The true value we use is Headings[a] - 1, so we can check for != 0.
var Headings = map[byte]int { 'E': 1, 'N': 2, 'W': 3, 'S': 4 }

// The direction to change the heading for 'left' or 'right' rotation.
// Left is clockwise (E->N->W->S) while right is counter-clockwise.
var Rotations = map[byte]int { 'L': 1, 'R': -1 }

// Parse a Direction from its string representation.
func NewDirection(input string) (Direction, error) {
  action := input[0]
  if Headings[action] == 0 && Rotations[action] == 0 && action != 'F' {
    return Direction{}, errors.New(fmt.Sprintf("invalid action '%c'", action))
  }
  value, err := strconv.Atoi(input[1:])
  if err != nil {
    return Direction{}, err
  }
  if Rotations[action] != 0 && (value % 90 != 0) {
    return Direction{}, errors.New(fmt.Sprintf(
      "non-cardinal rotation at '%s'", input))
  }
  return Direction{action, value}, nil
}

type Boat struct {
  lat int
  long int
  head int // cardinal heading: 0 is East, 1 is North, etc...
}

// You Must Build A Boat.
func NewBoat() *Boat {
  return new(Boat)
}

// Move the boat along a vector (without rotating it).
func (b *Boat) Move(mag int, head int) {
  if head % 2 == 1 { // N=1, S=3 => 1, -1
    // Positive latitude is in the north heading.
    latDelta := 2 - head
    b.lat = b.lat + (mag * latDelta)
  } else { // {E=0, W=2} => 1, -1
    // Positive longitude is in the east heading.
    longDelta := 1 - head
    b.long = b.long + (mag * longDelta)
  }
}

// Rotate the boat by adjusting the heading.
func (b *Boat) Rotate(headDelta int) {
  b.head += headDelta
  b.head = b.head % len(Headings)
  if b.head < 0 {
    b.head += len(Headings)
  }
}

// Follow a Direction.
func (b *Boat) MoveDirection(d Direction) {
  heading := Headings[d.action]
  if heading != 0 {
    b.Move(d.value, heading - 1)
  } else {
    rotation := Rotations[d.action]
    if rotation != 0 {
      b.Rotate(rotation * (d.value / 90))
    } else {
      b.Move(d.value, b.head)
    }
  }
}

// Follow several Directions.
func (b *Boat) Follow(directions []Direction) {
  for _, d := range directions {
    b.MoveDirection(d)
  }
}

func (b *Boat) LatStr() string {
  if b.lat >= 0 {
    return strconv.Itoa(b.lat) + "N"
  }
  return strconv.Itoa(-b.lat) + "S"
}

func (b *Boat) LongStr() string {
  if b.long >= 0 {
    return strconv.Itoa(b.long) + "E"
  }
  return strconv.Itoa(-b.long) + "W"
}

var HeadingStrings = []string { "E", "N", "W", "S" }

func (b *Boat) HeadStr() string {
  return HeadingStrings[b.head]
}

func (b *Boat) Str() string {
  return fmt.Sprintf("(%s, %s) heading %s",
    b.LatStr(), b.LongStr(), b.HeadStr())
}

func IAbs(i int) int {
  if i < 0 {
    i *= -1
  }
  return i
}

// Distance in the L1 norm (Manhattan distance) from a point.
func (b *Boat) L1Distance(fromLat int, fromLong int) int {
  return IAbs(b.lat - fromLat) + IAbs(b.long - fromLong)
}

func ReadDirections(input io.Reader) ([]Direction, error) {
  scanner := bufio.NewScanner(input)
  scanner.Split(bufio.ScanLines)
  directions := make([]Direction, 0, 768)
  for scanner.Scan() {
    direction, err := NewDirection(scanner.Text())
    if err != nil {
      return directions, err
    }
    directions = append(directions, direction)
  }
  return directions, scanner.Err()
}

func ReadDirectionsFromFile(path string) ([]Direction, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()
  return ReadDirections(file)
}

func Main(input_path string, verbose bool, args []string) error {
  directions, err := ReadDirectionsFromFile(input_path)
  if err != nil {
    return err
  }

  boat := NewBoat()
  boat.Follow(directions)
  fmt.Println("New position and heading:", boat.Str())
  fmt.Printf("Manhattan distance from start: %d\n", boat.L1Distance(0, 0))

  return nil
}
