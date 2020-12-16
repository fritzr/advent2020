package p12

import (
  "fmt"
  "io"
  "os"
  "bufio"
  "strconv"
  "errors"
  "github.com/fritzr/advent2020/util"
)

var gVerbose = false

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

  // Relative position of the waypoint. Used for waypoint travel.
  wayLat int
  wayLong int
}



// You Must Build A Boat.
func NewBoat(lat int, long int, head int, wayLat int, wayLong int) *Boat {
  return &Boat{lat, long, head, wayLat, wayLong}
}

// Reset the position of the boat.
func (b *Boat) Set(lat int, long int) {
  b.lat = lat
  b.long = long
}

func head2Coord(head int) (lat int, long int) {
  if head % 2 == 1 { // N=1, S=3 => 1, -1
    // Positive latitude is in the north heading.
    lat = 2 - head
  } else { // {E=0, W=2} => 1, -1
    // Positive longitude is in the east heading.
    long = 1 - head
  }
  return lat, long
}

// Move a point along a vector given in magnitude and heading.
func move(lat int, long int, mag int, head int) (int, int) {
  latDelta, longDelta := head2Coord(head)
  lat += mag * latDelta
  long += mag * longDelta
  return lat, long
}

// Move a point component-wise.
func moveTo(lat int, long int, latDelta int, longDelta int) (int, int) {
  lat += latDelta
  long += longDelta
  return lat, long
}

// Rotate a heading by a number of units (CCW 90 degree rotations).
func rotate(head int, headDelta int) int {
  head += headDelta
  head = head % len(Headings)
  if head < 0 {
    head += len(Headings)
  }
  return head
}

// Rotate a point a number of units (CCW 90 degree rotations) about the origin.
//
// headDelta must be one of {0, 1, 2, 3} for no rotation, 90CCW, 180, 90CW.
func rotatePoint(lat int, long int, headDelta int) (int, int) {
  // Swap for odd rotations (+/-90)
  if headDelta % 2 != 0 {
    tmpLat := lat
    lat = long
    long = tmpLat
  }

  // Negate latitude for -90, 180
  // lat *= -1 * ((headDelta / 2) * 2 - 1)
  if headDelta == 2 || headDelta == 3 {
    lat *= -1
  }

  // Negate longitude for 90, 180
  // long *= -1 * ((((headDelta + 1) / 2) % 2) * 2 - 1)
  if headDelta == 1 || headDelta == 2 {
    long *= -1
  }

  return lat, long

  /*
  switch headDelta {
  // Rotate CCW: swap and negate longitude
  case 1:
    lat := b.wayLat
    b.wayLat = b.wayLong
    b.wayLong = -lat
  // Flip: don't swap... negate both
  case 2:
    b.wayLat = -b.wayLat
    b.wayLong = -b.wayLong
  // Rotate CW: swap and negate latitude
  case 3:
    lat := b.wayLat
    b.wayLat = -b.wayLong
    b.wayLong = lat
  }
  */
}


// Follow a Direction.
func (b *Boat) Move(d Direction) {
  heading := Headings[d.action]
  if heading != 0 {
    // Move the waypoint with magnitude and heading.
    b.lat, b.long = move(b.lat, b.long, d.value, heading - 1)
    if gVerbose {
      fmt.Printf("Moving %s by %d to (%s, %s)\n",
        HeadingStrings[heading-1], d.value, b.LatStr(), b.LongStr())
    }
  } else {
    // Rotate the boat (adjust the heading).
    rotation := Rotations[d.action]
    if rotation != 0 {
      b.head = rotate(b.head, rotation * (d.value / 90))
      if gVerbose {
        fmt.Printf("Turning %c by %d => %s\n",
          d.action, d.value, b.HeadStr())
      }
    } else {
      // Move the boat along its current heading.
      b.lat, b.long = move(b.lat, b.long, d.value, b.head)
      if gVerbose {
        fmt.Printf("Moving %s by %d to (%s, %s)\n",
          b.HeadStr(), d.value, b.LatStr(), b.LongStr())
      }
    }
  }
}

// Follow several Directions.
func (b *Boat) Follow(directions []Direction) {
  for _, d := range directions {
    b.Move(d)
  }
}

// Move the waypoint (or move the boat to the waypoint).
func (b *Boat) MoveWaypoint(w Direction) {
  heading := Headings[w.action]
  if heading != 0 {
    // Move the waypoint with magnitude and heading.
    b.wayLat, b.wayLong = move(b.wayLat, b.wayLong, w.value, heading - 1)
  } else {
    rotation := Rotations[w.action]
    if rotation != 0 {
      // Rotate the waypoint about the boat N degrees.
      headDelta := rotate(0, (rotation * (w.value / 90)))
      b.wayLat, b.wayLong = rotatePoint(b.wayLat, b.wayLong, headDelta)
    } else {
      // Move to the waypoint N times.
      b.lat, b.long = moveTo(b.lat, b.long,
        b.wayLat * w.value, b.wayLong * w.value)
    }
  }
}

// Follow directions which modify the waypoint.
func (b *Boat) FollowWaypoint(directions []Direction) {
  for _, d := range directions {
    b.MoveWaypoint(d)
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

// Distance in the L1 norm (Manhattan distance) from a point.
func (b *Boat) L1Distance(fromLat int, fromLong int) int {
  return util.IAbs(b.lat - fromLat) + util.IAbs(b.long - fromLong)
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

func printBoat(boat *Boat, fromLat int, fromLong int) {
  fmt.Println("New position and heading:", boat.Str())
  fmt.Printf("Manhattan distance from start: %d\n",
    boat.L1Distance(fromLat, fromLong))
}

func Main(input_path string, verbose bool, args []string) error {
  gVerbose = verbose
  directions, err := ReadDirectionsFromFile(input_path)
  if err != nil {
    return err
  }

  // Part 1 -- follow directions using turtle mechanics.
  boat := NewBoat(/*pos:*/ 0, 0, /*head:*/0/*E*/, /*waypoint:*/ 1/*N*/, 10/*E*/)
  boat.Follow(directions)
  printBoat(boat, 0, 0)

  // Part 2 -- follow directions using waypoint mechanics.
  boat.Set(0, 0)
  boat.FollowWaypoint(directions)
  printBoat(boat, 0, 0)

  return nil
}
